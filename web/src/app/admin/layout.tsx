"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { AdminAuthProvider, useAdminAuth } from "@/lib/admin-auth-context";
import {
  ChartIcon,
  GlobeIcon,
  BookIcon,
  UsersIcon,
  HomeIcon,
  MenuIcon,
  XIcon,
  LogOutIcon,
  LoadingSpinner,
  ShieldIcon,
} from "@/components/icons";

const navItems = [
  { href: "/admin", icon: ChartIcon, label: "Dashboard", exact: true },
  { href: "/admin/languages", icon: GlobeIcon, label: "Languages" },
  { href: "/admin/content", icon: BookIcon, label: "Content" },
  { href: "/admin/users", icon: UsersIcon, label: "Users" },
  { href: "/admin/admins", icon: ShieldIcon, label: "Admins" },
];

function AdminContent({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const { admin, isLoading, isAuthenticated, logout } = useAdminAuth();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  // Skip auth check for login page
  const isLoginPage = pathname === "/admin/login";

  useEffect(() => {
    if (!isLoading && !isAuthenticated && !isLoginPage) {
      router.push("/admin/login");
    }
  }, [isLoading, isAuthenticated, isLoginPage, router]);

  // For login page, just render children
  if (isLoginPage) {
    return <>{children}</>;
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <LoadingSpinner className="w-8 h-8 text-purple" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  const handleLogout = () => {
    logout();
    router.push("/admin/login");
  };

  const isActive = (href: string, exact?: boolean) => {
    if (exact) return pathname === href;
    return pathname.startsWith(href) && href !== "/admin";
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Desktop Sidebar */}
      <aside className="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-64 lg:flex-col">
        <div className="flex flex-col flex-grow bg-gradient-to-b from-purple to-indigo">
          {/* Logo */}
          <div className="flex items-center h-16 px-6 border-b border-white/10">
            <Link href="/admin" className="flex items-center gap-2">
              <Image
                src="/logo.png"
                alt="Bawo"
                width={100}
                height={28}
                className="h-7 w-auto brightness-0 invert"
              />
              <span className="text-white/80 text-sm font-medium">Admin</span>
            </Link>
          </div>

          {/* Navigation */}
          <nav className="flex-1 px-4 py-6 space-y-2">
            {navItems.map((item) => {
              const active = item.exact ? pathname === item.href : isActive(item.href);
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-colors ${
                    active
                      ? "bg-white/20 text-white"
                      : "text-white/70 hover:bg-white/10 hover:text-white"
                  }`}
                >
                  <item.icon className="w-5 h-5" />
                  {item.label}
                </Link>
              );
            })}
          </nav>

          {/* Admin info & logout */}
          <div className="p-4 border-t border-white/10">
            <div className="px-4 py-2 mb-2">
              <p className="text-white font-medium truncate">{admin?.name}</p>
              <p className="text-white/60 text-sm truncate">{admin?.username}</p>
            </div>
            <button
              onClick={handleLogout}
              className="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium text-white/70 hover:bg-white/10 hover:text-white transition-colors"
            >
              <LogOutIcon className="w-5 h-5" />
              Sign Out
            </button>
            <Link
              href="/"
              className="mt-2 w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium text-white/70 hover:bg-white/10 hover:text-white transition-colors"
            >
              <HomeIcon className="w-5 h-5" />
              Back to Site
            </Link>
          </div>
        </div>
      </aside>

      {/* Mobile header */}
      <div className="lg:hidden fixed top-0 left-0 right-0 z-50 bg-purple border-b border-white/10">
        <div className="flex items-center justify-between h-16 px-4">
          <Link href="/admin" className="flex items-center gap-2">
            <Image
              src="/logo.png"
              alt="Bawo"
              width={80}
              height={24}
              className="h-6 w-auto brightness-0 invert"
            />
            <span className="text-white/80 text-xs font-medium">Admin</span>
          </Link>
          <button
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            className="p-2 text-white/80 hover:text-white"
          >
            {isMobileMenuOpen ? (
              <XIcon className="w-6 h-6" />
            ) : (
              <MenuIcon className="w-6 h-6" />
            )}
          </button>
        </div>
      </div>

      {/* Mobile menu */}
      {isMobileMenuOpen && (
        <div className="lg:hidden fixed inset-0 z-40 bg-purple pt-16">
          <nav className="p-4 space-y-2">
            {navItems.map((item) => {
              const active = item.exact ? pathname === item.href : isActive(item.href);
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={() => setIsMobileMenuOpen(false)}
                  className={`flex items-center gap-3 px-4 py-3 rounded-xl text-base font-medium transition-colors ${
                    active
                      ? "bg-white/20 text-white"
                      : "text-white/70 hover:bg-white/10 hover:text-white"
                  }`}
                >
                  <item.icon className="w-5 h-5" />
                  {item.label}
                </Link>
              );
            })}
            <button
              onClick={handleLogout}
              className="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-base font-medium text-white/70 hover:bg-white/10 hover:text-white transition-colors"
            >
              <LogOutIcon className="w-5 h-5" />
              Sign Out
            </button>
            <Link
              href="/"
              onClick={() => setIsMobileMenuOpen(false)}
              className="flex items-center gap-3 px-4 py-3 rounded-xl text-base font-medium text-white/70 hover:bg-white/10 hover:text-white transition-colors"
            >
              <HomeIcon className="w-5 h-5" />
              Back to Site
            </Link>
          </nav>
        </div>
      )}

      {/* Main content */}
      <main className="lg:pl-64 pt-16 lg:pt-0">
        <div className="min-h-screen p-4 lg:p-8">{children}</div>
      </main>
    </div>
  );
}

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AdminAuthProvider>
      <AdminContent>{children}</AdminContent>
    </AdminAuthProvider>
  );
}
