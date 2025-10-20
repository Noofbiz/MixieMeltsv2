import { useAuth, useCart } from "../context/Context";

const Header = ({ setPage }) => {
  const { user, logout } = useAuth();
  const { cart } = useCart();
  const cartItemCount = cart.reduce((sum, item) => sum + item.quantity, 0);

  return (
    <header className="bg-cream-200 shadow-md sticky top-0 z-10">
      <nav className="container mx-auto px-6 py-4 flex justify-between items-center">
        <div
          className="text-2xl font-bold font-serif text-brown-900 cursor-pointer"
          onClick={() => setPage("home")}
        >
          Mixie Melts
        </div>
        <div className="flex items-center space-x-6">
          <a
            href="#"
            className="text-brown-800 hover:text-terracotta-500 transition-colors"
            onClick={() => setPage("home")}
          >
            Home
          </a>
          <a
            href="#"
            className="text-brown-800 hover:text-terracotta-500 transition-colors"
            onClick={() => setPage("products")}
          >
            Products
          </a>
          <a
            href="#"
            className="relative text-brown-800 hover:text-terracotta-500 transition-colors"
            onClick={() => setPage("cart")}
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
                d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>
            {cartItemCount > 0 && (
              <span className="absolute -top-2 -right-3 bg-terracotta-500 text-cream-100 text-xs rounded-full h-5 w-5 flex items-center justify-center">
                {cartItemCount}
              </span>
            )}
          </a>
          {user ? (
            <div className="flex items-center space-x-4">
              <a
                href="#"
                className="text-brown-800 hover:text-terracotta-500 transition-colors"
                onClick={() => setPage("account")}
              >
                My Account
              </a>
              <button
                onClick={logout}
                className="px-4 py-2 text-sm text-terracotta-600 border border-terracotta-600 rounded-full hover:bg-terracotta-600 hover:text-cream-100 transition-colors"
              >
                Logout
              </button>
            </div>
          ) : (
            <button
              onClick={() => setPage("login")}
              className="px-5 py-2 text-sm bg-terracotta-500 text-cream-100 rounded-full hover:bg-terracotta-600 transition-colors"
            >
              Login
            </button>
          )}
        </div>
      </nav>
    </header>
  );
};

export default Header;
