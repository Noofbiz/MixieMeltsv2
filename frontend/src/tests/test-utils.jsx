import React from "react";
import { render } from "@testing-library/react";
import AuthProvider from "../context/AuthContext";
import CartProvider from "../context/CartContext";

const AllTheProviders = ({ children }) => {
  return (
    <AuthProvider>
      <CartProvider>{children}</CartProvider>
    </AuthProvider>
  );
};

const customRender = (ui, options) =>
  render(ui, { wrapper: AllTheProviders, ...options });

// re-export everything
export * from "@testing-library/react";

// override render method
export { customRender as render };
