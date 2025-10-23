import React from "react";
import { render, screen, waitFor } from "./test-utils";
import { describe, it, expect, vi, beforeEach } from "vitest";
import ProductListPage from "../components/ProductListPage";

// Mock the global fetch function
const fetch = vi.fn();
vi.stubGlobal("fetch", fetch);

const mockProducts = [
  {
    id: 1,
    name: "Vanilla Dream",
    price: 12.99,
    image: "vanilla.jpg",
    scent: "Sweet Vanilla",
  },
  {
    id: 2,
    name: "Lavender Fields",
    price: 12.99,
    image: "lavender.jpg",
    scent: "Calming Lavender",
  },
];

const mockSubscriptionBoxes = [
  {
    id: 1,
    name: "Monthly Surprise",
    description: "A curated box of seasonal melts.",
    price: 29.99,
    image: "box.jpg",
  },
];

describe("ProductListPage", () => {
  beforeEach(() => {
    // Reset mocks before each test
    vi.clearAllMocks();
  });

  it("should fetch products and subscription boxes and render them", async () => {
    // Setup mock fetch responses
    fetch.mockImplementation((url) => {
      if (url.endsWith("/products")) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve(mockProducts),
        });
      }
      if (url.endsWith("/products/subscription-boxes")) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve(mockSubscriptionBoxes),
        });
      }
      return Promise.reject(new Error(`Unknown endpoint: ${url}`));
    });

    render(<ProductListPage />);

    // Wait for the elements to appear
    await waitFor(() => {
      // Check for subscription box
      expect(screen.getByText("Monthly Surprise")).not.toBeNull();
      // Check for products
      expect(screen.getByText("Vanilla Dream")).not.toBeNull();
      expect(screen.getByText("Lavender Fields")).not.toBeNull();
    });

    // Verify fetch was called for both endpoints
    expect(fetch).toHaveBeenCalledWith("/products");
    expect(fetch).toHaveBeenCalledWith("/products/subscription-boxes");
  });

  it("should handle fetch errors gracefully", async () => {
    // Setup mock fetch to simulate an error
    fetch.mockRejectedValue(new Error("API is down"));

    render(<ProductListPage />);

    // Wait for the component to process the error
    await waitFor(() => {
      // The component should not crash and should render its static titles
      expect(screen.getAllByText("Subscription Boxes")).not.toBeNull();
      expect(screen.getAllByText("Our Collection")).not.toBeNull();
    });
  });
});
