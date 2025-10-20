import React, { createContext } from "react";

export const AuthContext = createContext();
export const useAuth = () => React.useContext(AuthContext);

export const CartContext = createContext();
export const useCart = () => React.useContext(CartContext);
