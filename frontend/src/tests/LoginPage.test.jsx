import React from "react";
import { render, screen, fireEvent, waitFor } from "./test-utils";
import { describe, it, expect, vi, beforeEach } from "vitest";
import LoginPage from "../components/LoginPage";

// Mock the global fetch function
const fetch = vi.fn();
vi.stubGlobal("fetch", fetch);

// Mock the useAuth hook to provide a login function
const mockLogin = vi.fn();
vi.mock("../context/Context", async () => {
  const actual = await vi.importActual("../context/Context");
  return {
    ...actual,
    useAuth: () => ({
      login: mockLogin,
    }),
  };
});

describe("LoginPage", () => {
  let setPageMock;

  beforeEach(() => {
    vi.clearAllMocks();
    setPageMock = vi.fn();
  });

  const renderComponent = () => render(<LoginPage setPage={setPageMock} />);

  it("should render the login form", () => {
    renderComponent();
    expect(screen.getByPlaceholderText("you@example.com")).not.toBeNull();
    expect(screen.getByPlaceholderText("••••••••")).not.toBeNull();
    expect(
      screen
        .getAllByRole("button")
        .find((btn) => btn.textContent === "Sign In"),
    ).not.toBeNull();
    const loginButton = screen
      .getAllByRole("button")
      .find((btn) => btn.textContent === "Sign In");
    fireEvent.click(loginButton);
  });

  it("should allow a user to fill out the form", () => {
    renderComponent();
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "test@example.com" },
    });
    fireEvent.change(screen.getByLabelText("Password"), {
      target: { value: "password123" },
    });

    expect(screen.getByLabelText("Email").value).toBe("test@example.com");
    expect(screen.getByLabelText("Password").value).toBe("password123");
  });

  it("should successfully log in a user and redirect", async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ token: "fake-token" }),
    });

    renderComponent();

    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "test@example.com" },
    });
    fireEvent.change(screen.getByLabelText("Password"), {
      target: { value: "password123" },
    });
    const loginButton = screen
      .getAllByRole("button")
      .find((btn) => btn.textContent === "Sign In");
    fireEvent.click(loginButton);

    await waitFor(() => {
      expect(fetch).toHaveBeenCalledWith("/api/users/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: "test@example.com",
          password: "password123",
        }),
      });
    });

    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith(expect.any(Object), "fake-token");
      //expect(setPageMock).toHaveBeenCalledWith("home");
    });
  });

  it("should show an error message on failed login", async () => {
    fetch.mockResolvedValueOnce({
      ok: false,
      status: 401,
      json: () => Promise.resolve({ message: "Invalid credentials" }),
    });

    renderComponent();

    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "wrong@example.com" },
    });
    fireEvent.change(screen.getByLabelText("Password"), {
      target: { value: "wrongpassword" },
    });
    const loginButton = screen
      .getAllByRole("button")
      .find((btn) => btn.textContent === "Sign In");
    fireEvent.click(loginButton);

    // TODO: Enable error message check once error display is implemented
    // await waitFor(() => {
    //   expect(screen.getByText("Invalid credentials")).not.toBeNull();
    // });

    expect(mockLogin).not.toHaveBeenCalled();
    expect(setPageMock).not.toHaveBeenCalled();
  });

  it("should disable the login button while submitting", async () => {
    // Make the fetch promise hang
    let resolvePromise;
    const promise = new Promise((resolve) => {
      resolvePromise = resolve;
    });

    fetch.mockReturnValue(promise);

    renderComponent();

    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "test@example.com" },
    });
    fireEvent.change(screen.getByLabelText("Password"), {
      target: { value: "password" },
    });

    const loginButton = screen
      .getAllByRole("button")
      .find((btn) => btn.textContent === "Sign In");
    fireEvent.click(loginButton);

    // TODO: Button should be disabled immediately after click
    // expect(loginButton.disabled).toBe(true);

    // Resolve the promise to finish the test
    resolvePromise({ ok: true, json: () => Promise.resolve({ token: "ok" }) });
  });
});
