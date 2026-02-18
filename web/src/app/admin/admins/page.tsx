"use client";

import { useEffect, useState } from "react";
import { useAdminAuth } from "@/lib/admin-auth-context";
import {
  PlusIcon,
  EditIcon,
  TrashIcon,
  XIcon,
  LoadingSpinner,
  ShieldIcon,
} from "@/components/icons";

interface Admin {
  id: string;
  username: string;
  email: string;
  name: string;
  is_active: boolean;
  is_superadmin: boolean;
  created_at: string;
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:4000";

export default function AdminsPage() {
  const { token, admin: currentAdmin } = useAdminAuth();
  const [admins, setAdmins] = useState<Admin[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [editingAdmin, setEditingAdmin] = useState<Admin | null>(null);
  const [selectedAdminId, setSelectedAdminId] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    username: "",
    email: "",
    password: "",
    name: "",
    is_active: true,
    is_superadmin: false,
  });
  const [newPassword, setNewPassword] = useState("");
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (token) {
      loadAdmins();
    }
  }, [token]);

  const loadAdmins = async () => {
    if (!token) return;

    try {
      const response = await fetch(`${API_URL}/admin/admins`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (!response.ok) {
        throw new Error("Failed to load admins");
      }

      const data = await response.json();
      setAdmins(data.admins || []);
    } catch (error) {
      console.error("Failed to load admins:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const openCreateModal = () => {
    setEditingAdmin(null);
    setFormData({
      username: "",
      email: "",
      password: "",
      name: "",
      is_active: true,
      is_superadmin: false,
    });
    setError(null);
    setShowModal(true);
  };

  const openEditModal = (admin: Admin) => {
    setEditingAdmin(admin);
    setFormData({
      username: admin.username,
      email: admin.email,
      password: "",
      name: admin.name,
      is_active: admin.is_active,
      is_superadmin: admin.is_superadmin,
    });
    setError(null);
    setShowModal(true);
  };

  const openPasswordModal = (adminId: string) => {
    setSelectedAdminId(adminId);
    setNewPassword("");
    setError(null);
    setShowPasswordModal(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;

    setIsSaving(true);
    setError(null);

    try {
      if (editingAdmin) {
        // Update existing admin
        const response = await fetch(`${API_URL}/admin/admins/${editingAdmin.id}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            name: formData.name,
            email: formData.email,
            is_active: formData.is_active,
            is_superadmin: formData.is_superadmin,
          }),
        });

        if (!response.ok) {
          const data = await response.json();
          throw new Error(data.message || "Failed to update admin");
        }
      } else {
        // Create new admin
        const response = await fetch(`${API_URL}/admin/admins`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            username: formData.username,
            email: formData.email,
            password: formData.password,
            name: formData.name,
          }),
        });

        if (!response.ok) {
          const data = await response.json();
          throw new Error(data.message || "Failed to create admin");
        }
      }

      setShowModal(false);
      await loadAdmins();
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsSaving(false);
    }
  };

  const handlePasswordChange = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || !selectedAdminId) return;

    setIsSaving(true);
    setError(null);

    try {
      const response = await fetch(`${API_URL}/admin/admins/${selectedAdminId}/password`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ new_password: newPassword }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || "Failed to change password");
      }

      setShowPasswordModal(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!token) return;
    if (id === currentAdmin?.id) {
      alert("You cannot delete your own account");
      return;
    }
    if (!confirm("Are you sure you want to delete this admin?")) return;

    try {
      const response = await fetch(`${API_URL}/admin/admins/${id}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || "Failed to delete admin");
      }

      await loadAdmins();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to delete admin");
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-[50vh] flex items-center justify-center">
        <LoadingSpinner className="w-8 h-8 text-purple" />
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl lg:text-3xl font-bold text-gray-900">Admin Users</h1>
        <button
          onClick={openCreateModal}
          className="flex items-center gap-2 bg-purple hover:bg-purple/90 text-white font-medium px-4 py-2 rounded-xl transition-colors"
        >
          <PlusIcon className="w-5 h-5" />
          Add Admin
        </button>
      </div>

      {/* Admins Table */}
      <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">Admin</th>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">Username</th>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">Email</th>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">Status</th>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">Role</th>
                <th className="text-left text-sm font-medium text-gray-500 px-6 py-4">Created</th>
                <th className="text-right text-sm font-medium text-gray-500 px-6 py-4">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {admins.map((admin) => (
                <tr key={admin.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-full bg-purple/20 flex items-center justify-center text-purple font-semibold">
                        {admin.name?.charAt(0) || "A"}
                      </div>
                      <span className="font-medium text-gray-900">{admin.name}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-gray-600">{admin.username}</td>
                  <td className="px-6 py-4 text-gray-600">{admin.email}</td>
                  <td className="px-6 py-4">
                    <span
                      className={`text-xs px-2 py-1 rounded-full ${
                        admin.is_active
                          ? "bg-green-100 text-green-700"
                          : "bg-red-100 text-red-700"
                      }`}
                    >
                      {admin.is_active ? "Active" : "Inactive"}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    {admin.is_superadmin ? (
                      <span className="flex items-center gap-1 text-xs px-2 py-1 rounded-full bg-purple/20 text-purple">
                        <ShieldIcon className="w-3 h-3" />
                        Superadmin
                      </span>
                    ) : (
                      <span className="text-xs px-2 py-1 rounded-full bg-gray-100 text-gray-600">
                        Admin
                      </span>
                    )}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-500">
                    {new Date(admin.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex items-center justify-end gap-2">
                      <button
                        onClick={() => openEditModal(admin)}
                        className="p-1 text-gray-400 hover:text-purple transition-colors"
                        title="Edit"
                      >
                        <EditIcon className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => openPasswordModal(admin.id)}
                        className="p-1 text-gray-400 hover:text-blue-600 transition-colors"
                        title="Change Password"
                      >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                        </svg>
                      </button>
                      {admin.id !== currentAdmin?.id && (
                        <button
                          onClick={() => handleDelete(admin.id)}
                          className="p-1 text-gray-400 hover:text-red-600 transition-colors"
                          title="Delete"
                        >
                          <TrashIcon className="w-4 h-4" />
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {admins.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500">No admin users found.</p>
          </div>
        )}
      </div>

      {/* Create/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-gray-900">
                {editingAdmin ? "Edit Admin" : "Add Admin"}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <XIcon className="w-5 h-5" />
              </button>
            </div>

            {error && (
              <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-xl text-red-600 text-sm">
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Name
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  placeholder="Full Name"
                  required
                />
              </div>

              {!editingAdmin && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Username
                  </label>
                  <input
                    type="text"
                    value={formData.username}
                    onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                    placeholder="username"
                    required
                  />
                </div>
              )}

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Email
                </label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  placeholder="admin@example.com"
                  required
                />
              </div>

              {!editingAdmin && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Password
                  </label>
                  <input
                    type="password"
                    value={formData.password}
                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                    placeholder="Min 8 characters"
                    minLength={8}
                    required
                  />
                </div>
              )}

              {editingAdmin && (
                <>
                  <div className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      id="is_active"
                      checked={formData.is_active}
                      onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
                      className="w-4 h-4 text-purple rounded focus:ring-purple"
                    />
                    <label htmlFor="is_active" className="text-sm text-gray-700">
                      Active (can log in)
                    </label>
                  </div>

                  <div className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      id="is_superadmin"
                      checked={formData.is_superadmin}
                      onChange={(e) => setFormData({ ...formData, is_superadmin: e.target.checked })}
                      className="w-4 h-4 text-purple rounded focus:ring-purple"
                    />
                    <label htmlFor="is_superadmin" className="text-sm text-gray-700">
                      Superadmin (full access)
                    </label>
                  </div>
                </>
              )}

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 py-2 text-gray-600 hover:bg-gray-100 rounded-xl transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isSaving}
                  className="flex-1 py-2 bg-purple hover:bg-purple/90 text-white font-medium rounded-xl transition-colors disabled:opacity-50"
                >
                  {isSaving ? (
                    <LoadingSpinner className="w-5 h-5 mx-auto" />
                  ) : editingAdmin ? (
                    "Save Changes"
                  ) : (
                    "Create Admin"
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Change Password Modal */}
      {showPasswordModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-gray-900">Change Password</h2>
              <button
                onClick={() => setShowPasswordModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <XIcon className="w-5 h-5" />
              </button>
            </div>

            {error && (
              <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-xl text-red-600 text-sm">
                {error}
              </div>
            )}

            <form onSubmit={handlePasswordChange} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  New Password
                </label>
                <input
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  placeholder="Min 8 characters"
                  minLength={8}
                  required
                />
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowPasswordModal(false)}
                  className="flex-1 py-2 text-gray-600 hover:bg-gray-100 rounded-xl transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isSaving}
                  className="flex-1 py-2 bg-purple hover:bg-purple/90 text-white font-medium rounded-xl transition-colors disabled:opacity-50"
                >
                  {isSaving ? (
                    <LoadingSpinner className="w-5 h-5 mx-auto" />
                  ) : (
                    "Change Password"
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
