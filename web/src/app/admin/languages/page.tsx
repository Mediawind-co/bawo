"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import { createClient } from "@/lib/api";
import {
  PlusIcon,
  EditIcon,
  TrashIcon,
  XIcon,
  LoadingSpinner,
} from "@/components/icons";
import type { language } from "@/lib/client";

export default function LanguagesPage() {
  const { token } = useAuth();
  const [languages, setLanguages] = useState<language.Language[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingLanguage, setEditingLanguage] = useState<language.Language | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    code: "",
    description: "",
    flag_emoji: "",
    is_active: true,
  });
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (token) {
      loadLanguages();
    }
  }, [token]);

  const loadLanguages = async () => {
    if (!token) return;

    const client = createClient(token);
    try {
      const res = await client.language.ListAllLanguages();
      setLanguages(res.languages);
    } catch (error) {
      console.error("Failed to load languages:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const openCreateModal = () => {
    setEditingLanguage(null);
    setFormData({
      name: "",
      code: "",
      description: "",
      flag_emoji: "",
      is_active: true,
    });
    setShowModal(true);
  };

  const openEditModal = (lang: language.Language) => {
    setEditingLanguage(lang);
    setFormData({
      name: lang.name,
      code: lang.code,
      description: lang.description,
      flag_emoji: lang.flag_emoji,
      is_active: lang.is_active,
    });
    setShowModal(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;

    setIsSaving(true);
    const client = createClient(token);

    try {
      if (editingLanguage) {
        await client.language.UpdateLanguage(editingLanguage.id, formData);
      } else {
        await client.language.CreateLanguage({
          name: formData.name,
          code: formData.code,
          description: formData.description,
          flag_emoji: formData.flag_emoji,
        });
      }
      setShowModal(false);
      await loadLanguages();
    } catch (error) {
      console.error("Failed to save language:", error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!token) return;
    if (!confirm("Are you sure you want to delete this language?")) return;

    const client = createClient(token);
    try {
      await client.language.DeleteLanguage(id);
      await loadLanguages();
    } catch (error) {
      console.error("Failed to delete language:", error);
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
        <h1 className="text-2xl lg:text-3xl font-bold text-gray-900">Languages</h1>
        <button
          onClick={openCreateModal}
          className="flex items-center gap-2 bg-purple hover:bg-purple/90 text-white font-medium px-4 py-2 rounded-xl transition-colors"
        >
          <PlusIcon className="w-5 h-5" />
          Add Language
        </button>
      </div>

      {/* Languages Grid */}
      <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {languages.map((lang) => (
          <div
            key={lang.id}
            className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100"
          >
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-3">
                <span className="text-4xl">{lang.flag_emoji}</span>
                <div>
                  <h3 className="font-semibold text-gray-900">{lang.name}</h3>
                  <p className="text-sm text-gray-500">Code: {lang.code}</p>
                </div>
              </div>
              <span
                className={`text-xs px-2 py-1 rounded-full ${
                  lang.is_active
                    ? "bg-green-100 text-green-700"
                    : "bg-gray-100 text-gray-500"
                }`}
              >
                {lang.is_active ? "Active" : "Inactive"}
              </span>
            </div>
            <p className="text-sm text-gray-600 mb-4 line-clamp-2">
              {lang.description || "No description"}
            </p>
            <div className="flex items-center gap-2">
              <button
                onClick={() => openEditModal(lang)}
                className="flex items-center gap-1 text-sm text-gray-600 hover:text-purple transition-colors"
              >
                <EditIcon className="w-4 h-4" />
                Edit
              </button>
              <button
                onClick={() => handleDelete(lang.id)}
                className="flex items-center gap-1 text-sm text-gray-600 hover:text-red-600 transition-colors"
              >
                <TrashIcon className="w-4 h-4" />
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>

      {languages.length === 0 && (
        <div className="text-center py-12 bg-white rounded-2xl shadow-sm">
          <p className="text-gray-500">No languages yet. Add your first language!</p>
        </div>
      )}

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-gray-900">
                {editingLanguage ? "Edit Language" : "Add Language"}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <XIcon className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Name
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  placeholder="e.g., Yoruba"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Code (ISO 639-1)
                </label>
                <input
                  type="text"
                  value={formData.code}
                  onChange={(e) =>
                    setFormData({ ...formData, code: e.target.value.toLowerCase() })
                  }
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  placeholder="e.g., yo"
                  maxLength={3}
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Flag Emoji
                </label>
                <input
                  type="text"
                  value={formData.flag_emoji}
                  onChange={(e) =>
                    setFormData({ ...formData, flag_emoji: e.target.value })
                  }
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  placeholder="e.g., 🇳🇬"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  placeholder="Brief description of the language"
                  rows={3}
                />
              </div>

              {editingLanguage && (
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    id="is_active"
                    checked={formData.is_active}
                    onChange={(e) =>
                      setFormData({ ...formData, is_active: e.target.checked })
                    }
                    className="w-4 h-4 text-purple rounded focus:ring-purple"
                  />
                  <label htmlFor="is_active" className="text-sm text-gray-700">
                    Active (visible to learners)
                  </label>
                </div>
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
                  ) : editingLanguage ? (
                    "Save Changes"
                  ) : (
                    "Create Language"
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
