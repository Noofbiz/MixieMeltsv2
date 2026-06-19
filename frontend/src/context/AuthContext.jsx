import { useState } from "react";
import { AuthContext } from "./Context";

const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);

  const login = (user, token) => {
    localStorage.setItem("token", token);
    setUser(user);

    // Attempt to merge any guest (session) cart into the authenticated user's cart.
    // Fire-and-forget: we don't block UI login flow on merge success/failure.
    try {
      const headers = { "Content-Type": "application/json" };
      if (token) {
        headers["Authorization"] = `Bearer ${token}`;
      }
      // Include credentials so the browser sends the session cookie set for guest carts.
      fetch("/carts/merge", {
        method: "POST",
        headers,
        body: JSON.stringify({ user_id: user && user.id ? user.id : 0 }),
        credentials: "include",
      }).catch((err) => {
        // Non-fatal: log for debugging
        // eslint-disable-next-line no-console
        console.warn("Cart merge failed:", err);
      });
    } catch (err) {
      // eslint-disable-next-line no-console
      console.warn("Cart merge error:", err);
    }
  };

  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;
