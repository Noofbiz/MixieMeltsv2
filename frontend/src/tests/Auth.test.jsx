import React from "react";
import { render, screen, fireEvent, waitFor } from "./test-utils";
import { describe, it, expect, vi, beforeEach } from "vitest";
import AppWrapper from "../App";

// Mock the global fetch function
const fetch = vi.fn();
vi.stubGlobal("fetch", fetch);

// Mock localStorage
const localStorageMock = (() => {
  let store = {};
  return {
    getItem: vi.fn((key) => store[key] || null),
    setItem: vi.fn((key, value) => {
      store[key] = value.toString();
    }),
    clear: vi.fn(() => {
      store = {};
    }),
    removeItem: vi.fn((key) => {
      delete store[key];
    }),
  };
})();
Object.defineProperty(window, "localStorage", { value: localStorageMock });

const mockUser = { id: 1, email: "user@test.com", is_admin: false };
const mockAdmin = { id: 2, email: "admin@test.com", is_admin: true };

describe("Authentication and Authorization", () => {
  const renderComponent = () => render(<AppWrapper />);

  beforeEach(() => {
    fetch.mockClear();
    localStorage.clear();
  });

  it("should allow a user to log in and see account details", async () => {
    fetch.mockImplementation((url) => {
      if (url.startsWith("/products")) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([]),
        });
      }
      if (url === "/api/users/login") {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ token: "fake-token" }),
        });
      }
      if (url === "/api/users/me") {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve(mockUser),
        });
      }
      return Promise.reject(new Error(`unhandled fetch request for ${url}`));
    });

    renderComponent();

    // 1. Initial state: User is logged out
    expect(screen.getByText("Login")).not.toBeNull();
    expect(screen.queryByText("Logout")).toBeNull();

    // 2. Navigate to login page
    fireEvent.click(
      screen.getAllByRole("button").find((btn) => btn.textContent === "Login"),
    );
    expect(await screen.findByLabelText("Email")).not.toBeNull();

    // 4. Fill and submit login form
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: mockUser.email },
    });
    fireEvent.change(screen.getByLabelText("Password"), {
      target: { value: "password" },
    });
    fireEvent.click(
      screen
        .getAllByRole("button")
        .find((btn) => btn.textContent === "Sign In"),
    );

    // 5. Post-login state: User is logged in and sees account info
    expect(await screen.findByText("My Account")).not.toBeNull();
    fireEvent.click(
      screen.getAllByRole("button").find((btn) => btn.textContent === "Logout"),
    );

    // 6. Logout
    expect(await screen.findByText("Login")).not.toBeNull();
    expect(screen.queryByText("Logout")).toBeNull();
  });

  it("should show an error on failed login", async () => {
    fetch.mockImplementation((url) => {
      if (url.startsWith("/products")) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([]),
        });
      }
      if (url === "/api/users/login") {
        return Promise.resolve({
          ok: false,
          status: 401,
          json: () => Promise.resolve({ message: "Invalid credentials" }),
        });
      }
      return Promise.reject(new Error(`unhandled fetch request for ${url}`));
    });

    renderComponent();

    // Navigate to login page
    fireEvent.click(
      screen.getAllByRole("button").find((btn) => btn.textContent === "Login"),
    );
    expect(await screen.findByLabelText("Email")).not.toBeNull();

    // Fill and submit login form with wrong credentials
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "wrong@test.com" },
    });
    fireEvent.change(screen.getByLabelText("Password"), {
      target: { value: "wrongpassword" },
    });
    // Navigate to front page signed in
    fireEvent.click(
      screen
        .getAllByRole("button")
        .find((btn) => btn.textContent === "Sign In"),
    );

    // Verify error message is shown
    expect(await screen.findByText("Invalid credentials")).not.toBeNull();

    // Verify user is not logged in
    expect(screen.queryByText("Logout")).toBeNull();
  });

  it("should not show admin panel for a regular user (via UI login)", async () => {
    // Use UI login flow to authenticate as a regular user
    fetch.mockImplementation((url) => {
      if (url.startsWith("/products")) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([]),
        });
      }
      if (url === "/api/users/login") {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ token: "fake-token" }),
        });
      }
      if (url === "/api/users/me") {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve(mockUser),
        });
      }
      return Promise.reject(new Error(`unhandled fetch request for ${url}`));
    });

    renderComponent();

    // Navigate to login page// Navigate to login page
    fireEvent.click(
      screen.getAllByRole("button").find((btn) => btn.textContent === "Login"),
    );
    expect(await screen.findByLabelText("Email")).not.toBeNull();

    // Fill and submit login form
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: mockUser.email },
    });
    fireEvent.change(screen.getByLabelText("Password"), {
      target: { value: "password" },
    });
    fireEvent.click(
      screen
        .getAllByRole("button")
        .find((btn) => btn.textContent === "Sign In"),
    );

    // Ensure login completed and header updated
    expect(await screen.findByText("My Account")).not.toBeNull();
    expect(screen.getByText("Logout")).not.toBeNull();

    // Navigate to account page
    fireEvent.click(screen.getByText("My Account"));

    // Verify user can see their email but not the admin panel
    await waitFor(() => {
      expect(screen.getByText(`Hello, ${mockUser.email}!`)).not.toBeNull();
    });

    expect(screen.queryByText("Add New Product")).toBeNull();
    expect(screen.queryByText("Add New Subscription Box")).toBeNull();
  });

  it("should show admin panel for an admin user (via UI login)", async () => {
    // Use UI login flow to authenticate as an admin user
    fetch.mockImplementation((url) => {
      if (url.startsWith("/products")) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([]),
        });
      }
      if (url === "/api/users/login") {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ token: "admin-token" }),
        });
      }
      if (url === "/api/users/me") {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve(mockAdmin),
        });
      }
      return Promise.reject(new Error(`unhandled fetch request for ${url}`));
    });

    renderComponent();

    // Navigate to login page
    fireEvent.click(
      screen.getAllByRole("button").find((btn) => btn.textContent === "Login"),
    );
    expect(await screen.findByLabelText("Email")).not.toBeNull();

    // Fill and submit login form as admin
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: mockAdmin.email },
    });
    fireEvent.change(screen.getByLabelText("Password"), {
      target: { value: "adminpassword" },
    });
    fireEvent.click(
      screen
        .getAllByRole("button")
        .find((btn) => btn.textContent === "Sign In"),
    );

    // Ensure login completed and header updated
    expect(await screen.findByText("My Account")).not.toBeNull();
    expect(
      screen.getAllByRole("button").find((btn) => btn.textContent === "Logout"),
    ).not.toBeNull();

    // Navigate to account page
    fireEvent.click(screen.getAllByText("My Account")[1]);

    // Verify admin can see theiSr email AND the admin panel
    await waitFor(() => {
      expect(screen.getByText(`Hello, ${mockAdmin.email}!`)).not.toBeNull();
    });

    expect(screen.getByText("Add New Product")).not.toBeNull();
    expect(screen.getByText("Add New Subscription Box")).not.toBeNull();
  });
});
