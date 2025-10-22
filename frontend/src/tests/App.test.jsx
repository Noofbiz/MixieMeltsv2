import { render, screen } from "./test-utils";
import { describe, it, expect, vi, beforeEach } from "vitest";
import AppWrapper from "../App";
import React from "react";

// Mocking fetch API as it's used in App.jsx's useEffect
const fetch = vi.fn();
vi.stubGlobal("fetch", fetch);

// Mocking localStorage as it's used to retrieve the auth token
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
    key: vi.fn((n) => Object.keys(store)[n]),
    get length() {
      return Object.keys(store).length;
    },
  };
})();
vi.stubGlobal("localStorage", localStorageMock);

describe("App", () => {
  beforeEach(() => {
    // Reset mocks and localStorage before each test
    vi.clearAllMocks();
    localStorage.clear();
  });

  it("should render the homepage by default", () => {
    fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve([]),
    });
    render(<AppWrapper />);
    // Check for a static element that is always present, like the footer
    const footerElement = screen.getByText(
      /© 2025 Mixie Melts. All Rights Reserved./i,
    );
    expect(footerElement).not.toBeNull();
  });

  it("should not fetch user if no token is in localStorage", () => {
    fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve([]),
    });
    render(<AppWrapper />);
    // The fetch in HomePage is called, but the user fetch in App is not.
    expect(fetch).toHaveBeenCalledWith("/products?limit=3");
    expect(fetch).not.toHaveBeenCalledWith("/api/users/me", expect.anything());
  });

  it("should attempt to fetch user if a token is in localStorage", () => {
    const token = "test-token";
    localStorage.setItem("token", token);

    fetch.mockImplementation((url) => {
      if (url.startsWith("/products")) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([]),
        });
      }
      if (url === "/api/users/me") {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ id: 1, name: "Test User" }),
        });
      }
      return Promise.reject(new Error(`unhandled fetch request for ${url}`));
    });

    render(<AppWrapper />);

    expect(fetch).toHaveBeenCalledWith("/api/users/me", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  });
});
