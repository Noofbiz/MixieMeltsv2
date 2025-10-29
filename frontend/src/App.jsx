import React, { useState, useEffect, useCallback } from "react";
import AuthProvider from "./context/AuthContext";
import CartProvider from "./context/CartContext";
import { useAuth } from "./context/Context";
import Header from "./components/Header";
import HomePage from "./components/HomePage";
import ProductListPage from "./components/ProductListPage";
import ProductDetailPage from "./components/ProductDetailPage";
import LoginPage from "./components/LoginPage";
import CartPage from "./components/CartPage";
import AccountPage from "./components/AccountPage";

/*
  App:
  - Supports a simple internal "routing" based on state and browser pathname.
  - If the URL matches /products/:id the ProductDetailPage is shown.
  - setPage accepts strings like "home", "products", "cart", "account" OR
    a product detail navigation via setPage("product", id).
  - When navigating we push state to history so users can bookmark /products/:id.
*/

function App() {
  const { login } = useAuth();

  // pageKind: "home" | "products" | "cart" | "login" | "account" | "product"
  const [pageKind, setPageKind] = useState("home");
  const [selectedProductId, setSelectedProductId] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      const fetchUser = async () => {
        try {
          const response = await fetch("/api/users/me", {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });
          if (response.ok) {
            const user = await response.json();
            login(user, token);
          }
        } catch (error) {
          console.error("Failed to fetch user", error);
        }
      };
      fetchUser();
    }
  }, [login]);

  // Initialize from current pathname so direct links work.
  useEffect(() => {
    const parsePath = () => {
      const p = window.location.pathname || "/";
      // match /products/:id
      const prodMatch = p.match(/^\/products\/(\d+)$/);
      if (prodMatch) {
        setPageKind("product");
        setSelectedProductId(Number(prodMatch[1]));
        return;
      }
      // other simple paths
      if (p === "/" || p === "") {
        setPageKind("home");
        return;
      }
      if (p === "/products") {
        setPageKind("products");
        return;
      }
      if (p === "/cart") {
        setPageKind("cart");
        return;
      }
      if (p === "/login") {
        setPageKind("login");
        return;
      }
      if (p === "/account") {
        setPageKind("account");
        return;
      }
      // default
      setPageKind("home");
    };

    parsePath();
    // Listen for popstate so back/forward work.
    const onPop = () => parsePath();
    window.addEventListener("popstate", onPop);
    return () => window.removeEventListener("popstate", onPop);
  }, []);

  // Unified function to change page and push history.
  const navigate = useCallback((kind, id) => {
    if (kind === "product" && id) {
      const newPath = `/products/${id}`;
      window.history.pushState({}, "", newPath);
      setSelectedProductId(Number(id));
      setPageKind("product");
      return;
    }

    // non-product pages
    let path = "/";
    switch (kind) {
      case "products":
        path = "/products";
        break;
      case "cart":
        path = "/cart";
        break;
      case "login":
        path = "/login";
        break;
      case "account":
        path = "/account";
        break;
      case "home":
      default:
        path = "/";
        break;
    }
    window.history.pushState({}, "", path);
    setSelectedProductId(null);
    setPageKind(kind);
  }, []);

  const renderPage = () => {
    switch (pageKind) {
      case "products":
        // allow ProductListPage to trigger navigation to a product detail by using
        // the browser URL (e.g. instruct users to open /products/:id) or via the global navigate function.
        // We pass a small helper via props to allow it to navigate.
        return <ProductListPage goTo={navigate} />;
      case "product":
        return (
          <ProductDetailPage
            productId={selectedProductId}
            onBack={() => navigate("products")}
          />
        );
      case "cart":
        return <CartPage />;
      case "login":
        return <LoginPage setPage={navigate} />;
      case "account":
        return <AccountPage />;
      case "home":
      default:
        return <HomePage setPage={navigate} />;
    }
  };

  return (
    <div className="bg-cream-100 min-h-screen font-sans text-brown-800">
      <Header setPage={navigate} />
      <main>{renderPage()}</main>
      <footer className="bg-cream-200 mt-12 py-6">
        <div className="container mx-auto text-center text-brown-700">
          &copy; 2025 Mixie Melts. All Rights Reserved.
        </div>
      </footer>
    </div>
  );
}

export default function AppWrapper() {
  return (
    <AuthProvider>
      <CartProvider>
        <App />
      </CartProvider>
    </AuthProvider>
  );
}
