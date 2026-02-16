"use client";

import { useState, useEffect } from "react";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { GoogleIcon, AppleIcon, LoadingSpinner } from "@/components/icons";
import { useAuth } from "@/lib/auth-context";

declare global {
  interface Window {
    google?: {
      accounts: {
        id: {
          initialize: (config: {
            client_id: string;
            callback: (response: { credential: string }) => void;
          }) => void;
          renderButton: (
            element: HTMLElement,
            config: { theme: string; size: string; width: number; text: string }
          ) => void;
        };
      };
    };
    AppleID?: {
      auth: {
        init: (config: {
          clientId: string;
          scope: string;
          redirectURI: string;
          usePopup: boolean;
        }) => void;
        signIn: () => Promise<{
          authorization: { id_token: string };
        }>;
      };
    };
  }
}

export default function LoginPage() {
  const router = useRouter();
  const { login, isAuthenticated, isLoading: authLoading } = useAuth();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isAuthenticated) {
      router.push("/learn");
    }
  }, [isAuthenticated, router]);

  useEffect(() => {
    // Load Google Sign-In script
    const googleScript = document.createElement("script");
    googleScript.src = "https://accounts.google.com/gsi/client";
    googleScript.async = true;
    googleScript.defer = true;
    document.body.appendChild(googleScript);

    googleScript.onload = () => {
      const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
      if (clientId && window.google) {
        window.google.accounts.id.initialize({
          client_id: clientId,
          callback: handleGoogleCallback,
        });
      }
    };

    return () => {
      document.body.removeChild(googleScript);
    };
  }, []);

  const handleGoogleCallback = async (response: { credential: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      await login(response.credential);
      router.push("/learn");
    } catch (err) {
      setError("Failed to sign in with Google. Please try again.");
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleGoogleLogin = () => {
    const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
    if (!clientId) {
      setError("Google Sign-In is not configured");
      return;
    }

    if (window.google) {
      const buttonWrapper = document.getElementById("google-btn-wrapper");
      if (buttonWrapper) {
        buttonWrapper.innerHTML = "";
        window.google.accounts.id.renderButton(buttonWrapper, {
          theme: "outline",
          size: "large",
          width: 300,
          text: "signin_with",
        });
        const button = buttonWrapper.querySelector("div[role='button']") as HTMLElement;
        if (button) {
          button.click();
        }
      }
    }
  };

  const handleAppleLogin = async () => {
    const clientId = process.env.NEXT_PUBLIC_APPLE_CLIENT_ID;
    if (!clientId) {
      setError("Apple Sign-In is not configured");
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      // Apple Sign-In requires additional setup
      // For now, show a placeholder message
      setError("Apple Sign-In coming soon");
    } catch (err) {
      setError("Failed to sign in with Apple. Please try again.");
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDevLogin = async () => {
    if (process.env.NODE_ENV !== "development") return;

    setIsLoading(true);
    setError(null);

    try {
      // Use a test token for dev - backend needs to handle this
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/dev-login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: "dev@bawo.test", name: "Dev User" }),
      });

      if (!response.ok) throw new Error("Dev login failed");

      const data = await response.json();
      await login(data.token);
      router.push("/learn");
    } catch (err) {
      setError("Dev login failed. Make sure backend is running.");
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-mint/20 to-purple/10">
        <LoadingSpinner className="w-8 h-8 text-teal" />
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-mint/20 to-purple/10 px-4">
      <div className="w-full max-w-md">
        {/* Logo */}
        <div className="text-center mb-8">
          <Link href="/">
            <Image
              src="/logo.png"
              alt="Bawo"
              width={150}
              height={40}
              className="h-12 w-auto mx-auto"
            />
          </Link>
          <h1 className="mt-6 text-3xl font-bold text-teal">Welcome Back</h1>
          <p className="mt-2 text-gray-600">
            Continue your language learning journey
          </p>
        </div>

        {/* Login Card */}
        <div className="bg-white rounded-2xl shadow-lg p-8">
          {error && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-600 text-sm">
              {error}
            </div>
          )}

          <div className="space-y-4">
            {/* Google Sign In */}
            <button
              onClick={handleGoogleLogin}
              disabled={isLoading}
              className="w-full flex items-center justify-center gap-3 bg-white border-2 border-gray-200 hover:border-gray-300 text-gray-700 font-medium px-6 py-3.5 rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? (
                <LoadingSpinner className="w-5 h-5" />
              ) : (
                <>
                  <GoogleIcon className="w-5 h-5" />
                  Continue with Google
                </>
              )}
            </button>

            {/* Apple Sign In */}
            <button
              onClick={handleAppleLogin}
              disabled={isLoading}
              className="w-full flex items-center justify-center gap-3 bg-black hover:bg-gray-900 text-white font-medium px-6 py-3.5 rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? (
                <LoadingSpinner className="w-5 h-5" />
              ) : (
                <>
                  <AppleIcon className="w-5 h-5" />
                  Continue with Apple
                </>
              )}
            </button>

            {/* Dev Login - only in development */}
            {process.env.NODE_ENV === "development" && (
              <>
                <div className="relative my-4">
                  <div className="absolute inset-0 flex items-center">
                    <div className="w-full border-t border-gray-200" />
                  </div>
                  <div className="relative flex justify-center text-sm">
                    <span className="px-2 bg-white text-gray-500">Dev Only</span>
                  </div>
                </div>
                <button
                  onClick={handleDevLogin}
                  disabled={isLoading}
                  className="w-full flex items-center justify-center gap-3 bg-purple/10 hover:bg-purple/20 text-purple font-medium px-6 py-3.5 rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isLoading ? (
                    <LoadingSpinner className="w-5 h-5" />
                  ) : (
                    "Dev Login (Skip OAuth)"
                  )}
                </button>
              </>
            )}
          </div>

          {/* Hidden wrapper for Google button */}
          <div id="google-btn-wrapper" className="hidden" />

          <div className="mt-8 text-center">
            <p className="text-sm text-gray-500">
              By continuing, you agree to our{" "}
              <a href="#" className="text-purple hover:underline">
                Terms of Service
              </a>{" "}
              and{" "}
              <a href="#" className="text-purple hover:underline">
                Privacy Policy
              </a>
            </p>
          </div>
        </div>

        {/* Back to home */}
        <div className="mt-8 text-center">
          <Link href="/" className="text-teal hover:text-mint transition-colors">
            &larr; Back to home
          </Link>
        </div>
      </div>
    </div>
  );
}
