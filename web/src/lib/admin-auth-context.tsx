"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from "react";

interface Admin {
  id: string;
  username: string;
  email: string;
  name: string;
  is_active: boolean;
  is_superadmin: boolean;
}

interface AdminAuthContextType {
  admin: Admin | null;
  token: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

const AdminAuthContext = createContext<AdminAuthContextType | undefined>(undefined);

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:4000";

export function AdminAuthProvider({ children }: { children: ReactNode }) {
  const [admin, setAdmin] = useState<Admin | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const savedToken = localStorage.getItem("bawo_admin_token");
    if (savedToken) {
      validateToken(savedToken);
    } else {
      setIsLoading(false);
    }
  }, []);

  const validateToken = async (savedToken: string) => {
    try {
      const response = await fetch(`${API_URL}/admin/auth/validate`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token: savedToken }),
      });

      const data = await response.json();
      if (data.valid && data.admin) {
        setAdmin(data.admin);
        setToken(savedToken);
      } else {
        localStorage.removeItem("bawo_admin_token");
      }
    } catch (error) {
      console.error("Token validation failed:", error);
      localStorage.removeItem("bawo_admin_token");
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (username: string, password: string) => {
    const response = await fetch(`${API_URL}/admin/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || "Login failed");
    }

    const data = await response.json();
    setAdmin(data.admin);
    setToken(data.token);
    localStorage.setItem("bawo_admin_token", data.token);
  };

  const logout = async () => {
    if (token) {
      try {
        await fetch(`${API_URL}/admin/auth/logout`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ token }),
        });
      } catch (error) {
        console.error("Logout request failed:", error);
      }
    }
    setAdmin(null);
    setToken(null);
    localStorage.removeItem("bawo_admin_token");
  };

  return (
    <AdminAuthContext.Provider
      value={{
        admin,
        token,
        isLoading,
        isAuthenticated: !!admin,
        login,
        logout,
      }}
    >
      {children}
    </AdminAuthContext.Provider>
  );
}

export function useAdminAuth() {
  const context = useContext(AdminAuthContext);
  if (context === undefined) {
    throw new Error("useAdminAuth must be used within an AdminAuthProvider");
  }
  return context;
}
