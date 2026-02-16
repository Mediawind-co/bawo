"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import { createClient } from "@/lib/api";
import {
  TargetIcon,
  CheckIcon,
  LoadingSpinner,
} from "@/components/icons";
import type { streak } from "@/lib/client";

const dailyGoalOptions = [
  { value: 10, label: "Casual", description: "10 XP/day" },
  { value: 20, label: "Regular", description: "20 XP/day" },
  { value: 50, label: "Serious", description: "50 XP/day" },
  { value: 100, label: "Intense", description: "100 XP/day" },
];

export default function SettingsPage() {
  const { token, user, logout } = useAuth();
  const [streakData, setStreakData] = useState<streak.StreakResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState(false);

  useEffect(() => {
    if (token) {
      loadData();
    }
  }, [token]);

  const loadData = async () => {
    if (!token) return;

    const client = createClient(token);
    try {
      const streakRes = await client.streak.GetMyStreak();
      setStreakData(streakRes);
    } catch (error) {
      console.error("Failed to load settings:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdateGoal = async (goal: number) => {
    if (!token) return;

    setIsSaving(true);
    setSaveSuccess(false);

    const client = createClient(token);
    try {
      const res = await client.streak.UpdateDailyGoal({ daily_xp_goal: goal });
      setStreakData(res);
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 2000);
    } catch (error) {
      console.error("Failed to update goal:", error);
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner className="w-8 h-8 text-teal" />
      </div>
    );
  }

  return (
    <div className="p-4 lg:p-8 pb-24 lg:pb-8">
      <h1 className="text-2xl lg:text-3xl font-bold text-teal mb-8">Settings</h1>

      {/* Profile Section */}
      <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 mb-6">
        <h2 className="text-lg font-semibold text-teal mb-4">Profile</h2>
        <div className="flex items-center gap-4">
          {user?.avatar_url ? (
            <img
              src={user.avatar_url}
              alt={user.name}
              className="w-16 h-16 rounded-full"
            />
          ) : (
            <div className="w-16 h-16 rounded-full bg-mint flex items-center justify-center text-teal text-2xl font-semibold">
              {user?.name?.charAt(0) || "U"}
            </div>
          )}
          <div>
            <p className="font-semibold text-gray-900">{user?.name}</p>
            <p className="text-sm text-gray-500">{user?.email}</p>
            <p className="text-xs text-gray-400 mt-1">
              Signed in with {user?.provider}
            </p>
          </div>
        </div>
      </div>

      {/* Daily Goal Section */}
      <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 mb-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-mint/20 rounded-xl flex items-center justify-center">
              <TargetIcon className="w-5 h-5 text-teal" />
            </div>
            <div>
              <h2 className="text-lg font-semibold text-teal">Daily Goal</h2>
              <p className="text-sm text-gray-500">
                Current: {streakData?.streak.daily_xp_goal || 50} XP/day
              </p>
            </div>
          </div>
          {saveSuccess && (
            <span className="flex items-center gap-1 text-sm text-green-600 bg-green-100 px-3 py-1 rounded-full">
              <CheckIcon className="w-4 h-4" />
              Saved!
            </span>
          )}
        </div>

        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
          {dailyGoalOptions.map((option) => {
            const isSelected = streakData?.streak.daily_xp_goal === option.value;
            return (
              <button
                key={option.value}
                onClick={() => handleUpdateGoal(option.value)}
                disabled={isSaving}
                className={`p-4 rounded-xl border-2 transition-colors text-left ${
                  isSelected
                    ? "border-mint bg-mint/10"
                    : "border-gray-200 hover:border-mint/50"
                } ${isSaving ? "opacity-50 cursor-not-allowed" : ""}`}
              >
                <p className="font-semibold text-teal">{option.label}</p>
                <p className="text-sm text-gray-500">{option.description}</p>
              </button>
            );
          })}
        </div>
      </div>

      {/* Streak Freezes */}
      {streakData && (
        <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 mb-6">
          <h2 className="text-lg font-semibold text-teal mb-4">Streak Freezes</h2>
          <p className="text-gray-600">
            You have{" "}
            <span className="font-semibold text-teal">
              {streakData.streak.streak_freezes}
            </span>{" "}
            streak freeze(s) available. Streak freezes protect your streak when you
            miss a day.
          </p>
        </div>
      )}

      {/* Account Actions */}
      <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
        <h2 className="text-lg font-semibold text-teal mb-4">Account</h2>
        <button
          onClick={logout}
          className="text-red-600 hover:text-red-700 font-medium"
        >
          Sign Out
        </button>
      </div>
    </div>
  );
}
