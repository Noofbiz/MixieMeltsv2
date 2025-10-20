import { useReducer } from "react";
import { CartContext } from "./Context";

const cartReducer = (state, action) => {
  switch (action.type) {
    case "ADD_ITEM":
      var existingItem = state.find((item) => item.id === action.payload.id);
      if (existingItem) {
        return state.map((item) =>
          item.id === action.payload.id
            ? { ...item, quantity: item.quantity + 1 }
            : item,
        );
      }
      return [...state, { ...action.payload, quantity: 1 }];
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
  const [cart, dispatch] = useReducer(cartReducer, []);
  return (
    <CartContext.Provider value={{ cart, dispatch }}>
      {children}
    </CartContext.Provider>
  );
};

export default CartProvider;
