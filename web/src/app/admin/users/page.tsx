"use client";

import { useEffect, useState } from "react";
import { useAdminAuth } from "@/lib/admin-auth-context";
import { createClient } from "@/lib/api";
import {
  TrashIcon,
  LoadingSpinner,
} from "@/components/icons";
import type { user } from "@/lib/client";

export default function UsersPage() {
  const { token } = useAdminAuth();
  const [users, setUsers] = useState<user.User[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const limit = 20;

  useEffect(() => {
    if (token) {
      loadUsers();
    }
  }, [token, page]);

  const loadUsers = async () => {
    if (!token) return;

    setIsLoading(true);
    const client = createClient(token);
    try {
      const res = await client.user.ListUsers({
        Limit: limit,
        Offset: page * limit,
      });
      setUsers(res.users);
      setTotal(res.total);
    } catch (error) {
      console.error("Failed to load users:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRoleChange = async (userId: string, newRole: string) => {
    if (!token) return;

    const client = createClient(token);
    try {
      await client.user.UpdateUserRole(userId, { role: newRole });
      await loadUsers();
    } catch (error) {
      console.error("Failed to update role:", error);
    }
  };

  const handleDelete = async (userId: string) => {
    if (!token) return;
    if (!confirm("Are you sure you want to delete this user? This action cannot be undone.")) {
      return;
    }

    const client = createClient(token);
    try {
      await client.user.DeleteUser(userId);
      await loadUsers();
    } catch (error) {
      console.error("Failed to delete user:", error);
    }
  };

  const totalPages = Math.ceil(total / limit);

  if (isLoading && users.length === 0) {
    return (
      <div className="min-h-[50vh] flex items-center justify-center">
        <LoadingSpinner className="w-8 h-8 text-purple" />
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl lg:text-3xl font-bold text-gray-900">Users</h1>
          <p className="text-gray-500">{total} total users</p>
        </div>
      </div>

      {/* Users Table */}
      <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">
                  User
                </th>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">
                  Email
                </th>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">
                  Provider
                </th>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">
                  Role
                </th>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">
                  Joined
                </th>
                <th className="text-right text-sm font-medium text-gray-500 px-6 py-4">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {users.map((u) => (
                <tr key={u.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      {u.avatar_url ? (
                        <img
                          src={u.avatar_url}
                          alt={u.name}
                          className="w-10 h-10 rounded-full"
                        />
                      ) : (
                        <div className="w-10 h-10 rounded-full bg-purple/20 flex items-center justify-center text-purple font-semibold">
                          {u.name?.charAt(0) || "U"}
                        </div>
                      )}
                      <span className="font-medium text-gray-900">{u.name}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-gray-600">{u.email}</td>
                  <td className="px-6 py-4">
                    <span className="text-sm text-gray-600 capitalize">{u.provider}</span>
                  </td>
                  <td className="px-6 py-4">
                    <select
                      value={u.role}
                      onChange={(e) => handleRoleChange(u.id, e.target.value)}
                      className={`text-sm px-3 py-1 rounded-lg border ${
                        u.role === "admin"
                          ? "border-purple/30 bg-purple/10 text-purple"
                          : "border-gray-200 bg-gray-50 text-gray-600"
                      }`}
                    >
                      <option value="user">User</option>
                      <option value="admin">Admin</option>
                    </select>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-500">
                    {new Date(u.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button
                      onClick={() => handleDelete(u.id)}
                      className="text-gray-400 hover:text-red-600 transition-colors"
                      title="Delete user"
                    >
                      <TrashIcon className="w-5 h-5" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-6 py-4 border-t border-gray-100">
            <p className="text-sm text-gray-500">
              Showing {page * limit + 1}-{Math.min((page + 1) * limit, total)} of {total}
            </p>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setPage(page - 1)}
                disabled={page === 0}
                className="px-3 py-1 text-sm rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Previous
              </button>
              <button
                onClick={() => setPage(page + 1)}
                disabled={page >= totalPages - 1}
                className="px-3 py-1 text-sm rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Next
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
