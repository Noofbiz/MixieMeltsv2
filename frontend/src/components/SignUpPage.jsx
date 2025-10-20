import React, { useState } from "react";

const SignUpPage = ({ setIsLogin }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    try {
      const response = await fetch("/api/users/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Sign up failed");
      }

      setIsLogin(true);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="flex items-center justify-center py-20">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-2xl shadow-xl">
        <h1 className="text-3xl font-bold font-serif text-center text-brown-900">
          Create an Account
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
          <div>
            <label
              htmlFor="confirm-password"
              className="text-sm font-medium text-brown-800"
            >
              Confirm Password
            </label>
            <input
              id="confirm-password"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-terracotta-500 focus:border-terracotta-500"
              placeholder="••••••••"
              required
            />
          </div>
          <button
            type="submit"
            className="w-full py-3 px-4 text-lg font-semibold text-cream-100 bg-terracotta-500 rounded-lg hover:bg-terracotta-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-terracotta-500"
          >
            Sign Up
          </button>
        </form>

        <div className="text-center">
          <p className="text-sm text-brown-800">
            Already have an account?{" "}
            <button
              onClick={() => setIsLogin(true)}
              className="font-medium text-terracotta-600 hover:text-terracotta-500"
            >
              Sign in
            </button>
          </p>
        </div>
      </div>
    </div>
  );
};

export default SignUpPage;
