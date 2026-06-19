import React from "react";
import AuthProvider from "../context/AuthContext";
import CartProvider from "../context/CartContext";

/**
 * AllTheProviders
 *
 * A small wrapper component that composes all context providers used by tests.
 * Keeping this component in its own file ensures fast refresh and satisfies the
 * rule that files exported from the tests directory only export components.
 *
 * Tests can import this as the `wrapper` for `render` so components under test
 * receive the same context behavior as the application.
 */
const AllTheProviders = ({ children }) => (
  <AuthProvider>
    <CartProvider>{children}</CartProvider>
  </AuthProvider>
);

export default AllTheProviders;
