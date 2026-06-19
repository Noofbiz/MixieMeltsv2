import { useReducer, useEffect, useCallback } from "react";
import { CartContext } from "./Context";
import { useAuth } from "./Context";

/*
  CartContext now:
  - exposes helper functions to add/remove/update cart items that persist to the carts service (best-effort)
  - keeps a local copy of cart state for fast UI updates
  - on login (when `user` becomes available) it attempts to load the server-side cart and reconcile it
*/

const cartReducer = (state, action) => {
  switch (action.type) {
    case "SET_CART":
      // payload is an array of items that should replace the current cart
      return Array.isArray(action.payload) ? action.payload : state;

    case "ADD_ITEM": {
      const existingItem = state.find((item) => item.id === action.payload.id);
      if (existingItem) {
        return state.map((item) =>
          item.id === action.payload.id
            ? {
                ...item,
                quantity: item.quantity + (action.payload.quantity || 1),
              }
            : item,
        );
      }
      return [
        ...state,
        { ...action.payload, quantity: action.payload.quantity || 1 },
      ];
    }

    case "REMOVE_ITEM":
      return state.filter((item) => item.id !== action.payload.id);

    case "UPDATE_QUANTITY":
      return state
        .map((item) =>
          item.id === action.payload.id
            ? { ...item, quantity: Math.max(0, action.payload.quantity) }
            : item,
        )
        .filter((item) => item.quantity > 0);

    default:
      return state;
  }
};

const CartProvider = ({ children }) => {
  const { user } = useAuth();
  const [cart, dispatch] = useReducer(cartReducer, []);

  // Helper: convert server cart items to the frontend shape.
  // The server returns { id, cart_id, product_id, quantity, running_low, added_at }.
  // We map to an object that contains at least `id` (product id) and `quantity`
  // and preserve any existing fields (name, price, image) if present in current cart.
  const mapServerItemsToLocal = useCallback(
    (serverItems) => {
      return serverItems.map((si) => {
        // try to find product metadata from current local cart
        const existing = cart.find((c) => c.id === si.product_id);
        return {
          // keep frontend `id` as the product id (consistent with existing UI code)
          id: si.product_id,
          product_id: si.product_id,
          cart_item_id: si.id, // server-side cart-item id (useful later)
          quantity: si.quantity,
          running_low: si.running_low,
          added_at: si.added_at,
          // metadata fallback from existing state
          name: existing?.name,
          price: existing?.price,
          image: existing?.image,
        };
      });
    },
    [cart],
  );

  // Load server cart when user becomes available
  useEffect(() => {
    let mounted = true;
    const loadCart = async () => {
      if (!user || !user.id) {
        // no authenticated user -> keep local cart (or clear, based on your product requirements)
        return;
      }
      try {
        const resp = await fetch(`/carts/${user.id}`);
        if (!mounted) return;
        if (!resp.ok) return; // best-effort: ignore errors
        const data = await resp.json();
        if (data && Array.isArray(data.items)) {
          const local = mapServerItemsToLocal(data.items);
          dispatch({ type: "SET_CART", payload: local });
        }
      } catch (e) {
        // Log the error for visibility while keeping behavior best-effort.
        // eslint-disable-next-line no-console
        console.warn("loadCart error:", e);
      }
    };
    loadCart();
    return () => {
      mounted = false;
    };
  }, [user, mapServerItemsToLocal]);

  // Persist helpers (best-effort). We update local state immediately and then call the API.
  // After a successful server operation we refresh the server cart and reconcile.
  const refreshServerCart = useCallback(
    async (userId) => {
      if (!userId) return;
      try {
        const resp = await fetch(`/carts/${userId}`);
        if (!resp.ok) return;
        const data = await resp.json();
        if (data && Array.isArray(data.items)) {
          const local = mapServerItemsToLocal(data.items);
          dispatch({ type: "SET_CART", payload: local });
        }
      } catch (e) {
        // Log to aid debugging while otherwise ignoring transient errors.
        // eslint-disable-next-line no-console
        console.warn("refresh load error:", e);
      }
    },
    [mapServerItemsToLocal],
  );

  const addItemToCart = useCallback(
    async (product, quantity = 1) => {
      // optimistic local update
      dispatch({ type: "ADD_ITEM", payload: { ...product, quantity } });

      if (!user || !user.id) {
        // anonymous or not logged in: we only keep local cart
        return;
      }

      try {
        await fetch(`/carts/${user.id}/items`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ product_id: product.id, quantity }),
        });
        // refresh authoritative server representation
        await refreshServerCart(user.id);
      } catch (e) {
        // Log server persistence failures so we can debug issues while keeping UX optimistic.
        // eslint-disable-next-line no-console
        console.warn("addItemToCart persist error:", e);
      }
    },
    [user, refreshServerCart],
  );

  const removeItemFromCart = useCallback(
    async (productId) => {
      // optimistic local update
      dispatch({ type: "REMOVE_ITEM", payload: { id: productId } });

      if (!user || !user.id) {
        return;
      }

      // Best-effort server-side removal:
      // The carts service exposes DELETE /carts/{user_id}/items/{itemID} and needs the cart_item id.
      // Because the frontend stores product-id keyed items we request the server cart and attempt to delete the matching cart_item.
      try {
        const resp = await fetch(`/carts/${user.id}`);
        if (!resp.ok) return;
        const data = await resp.json();
        if (!data || !Array.isArray(data.items)) return;
        const match = data.items.find((it) => it.product_id === productId);
        if (!match) return;
        // call delete by cart_item id
        await fetch(`/carts/${user.id}/items/${match.id}`, {
          method: "DELETE",
        });
        await refreshServerCart(user.id);
      } catch (e) {
        // Log removal errors for diagnostics; operation is best-effort.
        // eslint-disable-next-line no-console
        console.warn("removeItemFromCart error:", e);
      }
    },
    [user, refreshServerCart],
  );

  const updateItemQuantity = useCallback(
    async (productId, quantity) => {
      // optimistic local update
      dispatch({
        type: "UPDATE_QUANTITY",
        payload: { id: productId, quantity },
      });

      if (!user || !user.id) {
        return;
      }

      try {
        // need to map productId -> cart_item id on the server
        const resp = await fetch(`/carts/${user.id}`);
        if (!resp.ok) return;
        const data = await resp.json();
        if (!data || !Array.isArray(data.items)) return;
        const match = data.items.find((it) => it.product_id === productId);
        if (!match) return;
        // PATCH expects body { quantity: <int> }
        await fetch(`/carts/${user.id}/items/${match.id}`, {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ quantity }),
        });
        await refreshServerCart(user.id);
      } catch (e) {
        // Log update errors for diagnostics; operation is best-effort.
        // eslint-disable-next-line no-console
        console.warn("updateItemQuantity error:", e);
      }
    },
    [user, refreshServerCart],
  );

  // Expose cart + helpers
  return (
    <CartContext.Provider
      value={{
        cart,
        // legacy dispatch for internal components/tests; prefer helpers
        dispatch,
        addItemToCart,
        removeItemFromCart,
        updateItemQuantity,
        refreshServerCart,
      }}
    >
      {children}
    </CartContext.Provider>
  );
};

export default CartProvider;
