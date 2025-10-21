import React, { useState } from "react";
import { useAuth } from "../context/Context";
import SignUpPage from "./SignUpPage";

const LoginPage = ({ setPage }) => {
  const [isLogin, setIsLogin] = useState(true);

  return isLogin ? (
    <LoginView setPage={setPage} setIsLogin={setIsLogin} />
  ) : (
    <SignUpPage setPage={setPage} setIsLogin={setIsLogin} />
  );
};

const LoginView = ({ setPage, setIsLogin }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    try {
      const response = await fetch("/api/users/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email: email, password: password }),
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Login failed");
      }

      const { token } = await response.json();
      // In a real app, you would fetch the user profile here
      login({ email }, token);
      setPage("home");
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="flex items-center justify-center py-20">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-2xl shadow-xl">
        <h1 className="text-3xl font-bold font-serif text-center text-brown-900">
          Welcome Back
        </h1>
        {error && <p className="text-red-500 text-center">{error}</p>}
        <form className="space-y-6" onSubmit={handleSubmit}>
          <div>
            <label
              htmlFor="email"
              className="text-sm font-medium text-brown-800"
            >
              Email
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-terracotta-500 focus:border-terracotta-500"
              placeholder="you@example.com"
              required
            />
          </div>
          <div>
            <label
              htmlFor="password"
              className="text-sm font-medium text-brown-800"
            >
              Password
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-terracotta-500 focus:border-terracotta-500"
              placeholder="••••••••"
              required
            />
          </div>
          <button
            type="submit"
            className="w-full py-3 px-4 text-lg font-semibold text-cream-100 bg-terracotta-500 rounded-lg hover:bg-terracotta-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-terracotta-500"
          >
            Sign In
          </button>
        </form>

        <div className="text-center">
          <p className="text-sm text-brown-800">
            Don't have an account?{" "}
            <button
              onClick={() => setIsLogin(false)}
              className="font-medium text-terracotta-600 hover:text-terracotta-500"
            >
              Sign up
            </button>
          </p>
        </div>

        <div className="relative flex py-5 items-center">
          <div className="flex-grow border-t border-gray-300"></div>
          <span className="flex-shrink mx-4 text-gray-500">
            Or continue with
          </span>
          <div className="flex-grow border-t border-gray-300"></div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <a
            href="/api/users/oauth/google/login"
            className="w-full flex items-center justify-center py-2.5 px-4 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-terracotta-500"
          >
            <svg className="w-5 h-5 mr-2" viewBox="0 0 48 48">
              <path
                fill="#FFC107"
                d="M43.611 20.083H42V20H24v8h11.303c-1.649 4.657-6.08 8-11.303 8c-6.627 0-12-5.373-12-12s5.373-12 12-12c3.059 0 5.842 1.154 7.961 3.039L38.802 8.94C34.343 4.932 28.36 2 24 2C11.82 2 2 11.82 2 24s9.82 22 22 22s22-9.82 22-22c0-1.341-.138-2.65-.389-3.917z"
              ></path>
              <path
                fill="#FF3D00"
                d="M6.306 14.691L14.63 21.13C12.553 23.336 12 26.233 12 28c0 4.418 3.582 8 8 8h6v-8h-6c-2.209 0-4-1.791-4-4s1.791-4 4-4h1.303c-1.65-4.657-6.08-8-11.303-8c-3.518 0-6.736 1.252-9.284 3.321z"
              ></path>
              <path
                fill="#4CAF50"
                d="M24 44c5.166 0 9.86-1.977 13.409-5.192L32.191 32.19C30.078 33.748 27.218 34 24 34c-4.418 0-8-3.582-8-8c0-1.767.553-4.664 2.63-6.87L6.306 14.691C3.866 17.657 2 21.734 2 24c0 10.318 8.343 18.673 18.673 19.962L24 44z"
              ></path>
              <path
                fill="#1976D2"
                d="M43.611 20.083H42V20H24v8h11.303c-1.649 4.657-6.08 8-11.303 8-3.218 0-6.078-.752-8.191-2.19L10.591 38.6C14.134 41.068 18.86 42 24 42c12.18 0 22-9.82 22-22c0-1.341-.138-2.65-.389-3.917z"
              ></path>
            </svg>
            Google
          </a>
          <a
            href="/api/users/oauth/facebook/login"
            className="w-full flex items-center justify-center py-2.5 px-4 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-terracotta-500"
          >
            <svg className="w-5 h-5 mr-2" fill="#1877F2" viewBox="0 0 24 24">
              <path d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.879V14.89H8.22v-2.89h2.218v-2.174c0-2.194 1.313-3.394 3.29-3.394.945 0 1.933.167 1.933.167v2.45h-1.25c-1.08 0-1.41.65-1.41 1.34v1.614h2.77l-.443 2.89h-2.327V21.88C18.343 21.128 22 16.991 22 12z"></path>
            </svg>
            Facebook
          </a>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
