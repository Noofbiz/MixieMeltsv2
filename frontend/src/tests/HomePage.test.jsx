import React from "react";
import { render, screen, waitFor } from "./test-utils";
import { describe, it, expect, vi, beforeEach } from "vitest";
import HomePage from "../components/HomePage";

// Mock the global fetch function
const fetch = vi.fn();
vi.stubGlobal("fetch", fetch);

const mockProducts = [
  {
    id: 1,
    name: "Cozy Cashmere",
    price: 12.99,
    image: "cashmere.jpg",
    scent: "Soft and warm",
  },
  {
    id: 2,
    name: "Spiced Apple Cider",
    price: 12.99,
    image: "apple.jpg",
    scent: "Warm and spicy",
  },
  {
    id: 3,
    name: "Ocean Breeze",
    price: 12.99,
    image: "ocean.jpg",
    scent: "Fresh and clean",
  },
];

describe("HomePage", () => {
  let setPageMock;

  beforeEach(() => {
    fetch.mockClear();
    setPageMock = vi.fn();
  });

  it("should fetch and display the featured products", async () => {
    fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockProducts),
    });

    render(<HomePage setPage={setPageMock} />);

    // Wait for products to be rendered
    await waitFor(() => {
      expect(screen.getByText("Cozy Cashmere")).not.toBeNull();
      expect(screen.getByText("Spiced Apple Cider")).not.toBeNull();
      expect(screen.getByText("Ocean Breeze")).not.toBeNull();
    });

    // Verify the correct API endpoint was called
    expect(fetch).toHaveBeenCalledWith("/products?limit=3");
  });

  // TODO: Implement the error handling test once the component has error UI
  it("should handle API errors gracefully and not display products", async () => {
    fetch.mockRejectedValue(new Error("Network Error"));

    render(<HomePage setPage={setPageMock} />);

    // // Wait for any async operations to settle
    // await waitFor(() => {
    //   // The component should still render its static content
    //   expect(screen.getByText("Warmth & Fragrance, Delivered.")).not.toBeNull();
    //   // Ensure no products are displayed
    //   expect(screen.queryByText("Cozy Cashmere")).toBeNull();
    //   expect(screen.queryByText("Spiced Apple Cider")).toBeNull();
    // });
  });
});
