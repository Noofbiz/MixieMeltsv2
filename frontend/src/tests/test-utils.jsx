/* eslint-disable react-refresh/only-export-components */
import { render } from "@testing-library/react";
import AllTheProviders from "./AllTheProviders";

/*
  test-utils.jsx

  Helper test utilities for the frontend tests.

  - Exports a `render` function that wraps components with the app providers used
    in tests (Auth + Cart).
  - Exports `setCartFetchMock` to allow tests to override the default fetch mock.
  - Re-exports everything from @testing-library/react for convenience.

  The default fetch mock handles cart-related endpoints used by CartContext so
  tests that render components which interact with the cart do not need to
  stub every endpoint. Tests that require more control should call
  `setCartFetchMock` or stub global fetch directly via their test framework.
*/

/**
 * Default fetch mock focused on cart-related endpoints:
 *  - POST /carts/merge
 *  - GET /carts/:userId
 *  - POST /carts/:userId/items
 *  - PATCH /carts/:userId/items/:itemId
 *  - DELETE /carts/:userId/items/:itemId
 *
 * Returns objects compatible with the Fetch API response usage (ok, status, json).
 */
const defaultFetch = (input, init = {}) => {
  const url = typeof input === "string" ? input : input?.url || "";
  const method = (init && init.method) || "GET";

  const makeResp = (status = 200, body = {}) => {
    return Promise.resolve({
      ok: status >= 200 && status < 300,
      status,
      json: async () => body,
      text: async () => JSON.stringify(body),
    });
  };

  // Normalize to pathname when absolute URLs are passed
  let path = url;
  try {
    if (url.startsWith("http://") || url.startsWith("https://")) {
      path = new URL(url).pathname;
    }
  } catch {
    // if URL parsing fails, fall back to raw string
  }

  // carts merge endpoint (accept variations: trailing slash, query string, absolute/relative URLs)
  // Match:
  //   /carts/merge
  //   /carts/merge/
  //   /carts/merge?foo=1
  //   /carts/merge/?foo=1
  if (/^\/carts\/merge(?:\/)?(?:\?.*)?$/.test(path) && method === "POST") {
    return makeResp(200, {});
  }

  // GET /carts/:userId
  if (/^\/carts\/\d+$/.test(path) && method === "GET") {
    return makeResp(200, { items: [] });
  }

  // POST /carts/:userId/items
  if (/^\/carts\/\d+\/items$/.test(path) && method === "POST") {
    return makeResp(201, {});
  }

  // PATCH or DELETE /carts/:userId/items/:itemId
  if (
    /^\/carts\/\d+\/items\/\d+$/.test(path) &&
    (method === "PATCH" || method === "DELETE")
  ) {
    return makeResp(200, {});
  }

  // Fallback: return a generic empty successful response so tests don't blow up.
  return makeResp(200, {});
};

/**
 * Replace the global fetch implementation used by tests.
 * - If `impl` is a function it will be used as global fetch.
 * - If `impl` is falsy, the default cart-aware mock will be installed.
 */
const setCartFetchMock = (impl) => {
  if (typeof impl === "function") {
    globalThis.fetch = impl;
  } else {
    globalThis.fetch = defaultFetch;
  }
};

// Ensure there is a default fetch mock available for tests that don't set one.
if (typeof globalThis.fetch === "undefined") {
  globalThis.fetch = defaultFetch;
}

/**
 * customRender
 *
 * Wrap renders with application providers used by the frontend so components
 * under test get access to AuthContext / CartContext behavior without each
 * test having to compose providers.
 */
const customRender = (ui, options) =>
  render(ui, { wrapper: AllTheProviders, ...options });

// re-export everything from testing-library for convenience in tests
export * from "@testing-library/react";

// override render with our wrapped render and export setCartFetchMock helper
export { customRender as render, setCartFetchMock };
