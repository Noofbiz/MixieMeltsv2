import { useCart } from "../context/Context";

const CartPage = () => {
  const { cart, dispatch } = useCart();
  const totalPrice = cart.reduce(
    (sum, item) => sum + item.price * item.quantity,
    0,
  );

  const handleQuantityChange = (id, quantity) => {
    dispatch({
      type: "UPDATE_QUANTITY",
      payload: { id, quantity: parseInt(quantity, 10) },
    });
  };
  const handleRemove = (itemToRemove) => {
    dispatch({ type: "REMOVE_ITEM", payload: itemToRemove });
  };

  const handleCheckout = () => {
    alert("Checkout initiated! This is where we would call the order service.");
  };

  return (
    <div className="container mx-auto px-6 py-12">
      <h1 className="text-4xl font-bold font-serif text-center text-brown-900 mb-12">
        Your Shopping Cart
      </h1>
      {cart.length === 0 ? (
        <p className="text-center text-brown-700 text-xl">
          Your cart is empty.
        </p>
      ) : (
        <div className="bg-white rounded-xl shadow-lg p-8">
          {cart.map((item) => (
            <div
              key={item.id}
              className="flex items-center justify-between border-b border-cream-200 py-4"
            >
              <div className="flex items-center space-x-4">
                <img
                  src={item.image}
                  alt={item.name}
                  className="w-20 h-20 object-cover rounded-lg"
                />
                <div>
                  <h2 className="text-lg font-semibold font-serif text-brown-900">
                    {item.name}
                  </h2>
                  <p className="text-brown-700">${item.price.toFixed(2)}</p>
                </div>
              </div>
              <div className="flex items-center space-x-4">
                <input
                  type="number"
                  value={item.quantity}
                  onChange={(e) =>
                    handleQuantityChange(item.id, e.target.value)
                  }
                  className="w-16 p-2 border rounded-md text-center focus:ring-terracotta-500 focus:border-terracotta-500"
                  min="1"
                />
                <p className="text-lg font-semibold w-24 text-right text-brown-800">
                  ${(item.price * item.quantity).toFixed(2)}
                </p>
                <button
                  onClick={() => handleRemove(item)}
                  className="text-red-500 hover:text-red-700"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    />
                  </svg>
                </button>
              </div>
            </div>
          ))}
          <div className="mt-8 flex justify-end items-center">
            <span className="text-2xl font-bold text-brown-900">
              Total: ${totalPrice.toFixed(2)}
            </span>
            <button
              onClick={handleCheckout}
              className="ml-6 px-8 py-3 bg-brown-800 text-cream-100 text-lg font-semibold rounded-full hover:bg-brown-900"
            >
              Proceed to Checkout
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default CartPage;
