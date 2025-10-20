import React, { useState, useEffect } from "react";
import AuthProvider from "./context/AuthContext";
import CartProvider from "./context/CartContext";
import { useAuth } from "./context/Context";
import Header from "./components/Header";
import HomePage from "./components/HomePage";
import ProductListPage from "./components/ProductListPage";
import LoginPage from "./components/LoginPage";
import CartPage from "./components/CartPage";
import AccountPage from "./components/AccountPage";

function App() {
  const [page, setPage] = useState("home");
  const { login } = useAuth();

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

  const renderPage = () => {
    switch (page) {
      case "products":
        return <ProductListPage />;
      case "cart":
        return <CartPage />;
      case "login":
        return <LoginPage setPage={setPage} />;
      case "account":
        return <AccountPage />;
      case "home":
      default:
        return <HomePage setPage={setPage} />;
    }
  };

  return (
    <div className="bg-cream-100 min-h-screen font-sans text-brown-800">
      <Header setPage={setPage} />
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
