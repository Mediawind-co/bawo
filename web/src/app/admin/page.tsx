"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import { createClient } from "@/lib/api";
import {
  UsersIcon,
  GlobeIcon,
  BookIcon,
  FireIcon,
  LoadingSpinner,
} from "@/components/icons";
import type { dashboard } from "@/lib/client";

export default function AdminDashboard() {
  const { token } = useAuth();
  const [overview, setOverview] = useState<dashboard.OverviewStats | null>(null);
  const [userStats, setUserStats] = useState<dashboard.UserStats | null>(null);
  const [contentStats, setContentStats] = useState<dashboard.ContentStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (token) {
      loadData();
    }
  }, [token]);

  const loadData = async () => {
    if (!token) return;

    const client = createClient(token);
    try {
      const [overviewRes, userRes, contentRes] = await Promise.all([
        client.dashboard.GetOverview(),
        client.dashboard.GetUserAnalytics(),
        client.dashboard.GetContentAnalytics(),
      ]);

      setOverview(overviewRes.stats);
      setUserStats(userRes.stats);
      setContentStats(contentRes.stats);
    } catch (error) {
      console.error("Failed to load dashboard data:", error);
    } finally {
      setIsLoading(false);
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
      <h1 className="text-2xl lg:text-3xl font-bold text-gray-900 mb-8">
        Admin Dashboard
      </h1>

      {/* Overview Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div className="bg-white rounded-2xl p-6 shadow-sm">
          <div className="w-10 h-10 bg-purple/20 rounded-xl flex items-center justify-center mb-3">
            <UsersIcon className="w-5 h-5 text-purple" />
          </div>
          <p className="text-2xl font-bold text-gray-900">
            {overview?.total_users || 0}
          </p>
          <p className="text-sm text-gray-500">Total Users</p>
        </div>

        <div className="bg-white rounded-2xl p-6 shadow-sm">
          <div className="w-10 h-10 bg-mint/20 rounded-xl flex items-center justify-center mb-3">
            <GlobeIcon className="w-5 h-5 text-teal" />
          </div>
          <p className="text-2xl font-bold text-gray-900">
            {contentStats?.total_languages || 0}
          </p>
          <p className="text-sm text-gray-500">Languages</p>
        </div>

        <div className="bg-white rounded-2xl p-6 shadow-sm">
          <div className="w-10 h-10 bg-indigo/20 rounded-xl flex items-center justify-center mb-3">
            <BookIcon className="w-5 h-5 text-indigo" />
          </div>
          <p className="text-2xl font-bold text-gray-900">
            {contentStats?.total_lessons || 0}
          </p>
          <p className="text-sm text-gray-500">Total Lessons</p>
        </div>

        <div className="bg-white rounded-2xl p-6 shadow-sm">
          <div className="w-10 h-10 bg-orange-100 rounded-xl flex items-center justify-center mb-3">
            <FireIcon className="w-5 h-5 text-orange-500" />
          </div>
          <p className="text-2xl font-bold text-gray-900">
            {overview?.active_users_today || 0}
          </p>
          <p className="text-sm text-gray-500">Active Today</p>
        </div>
      </div>

      {/* Detailed Stats */}
      <div className="grid lg:grid-cols-2 gap-6">
        {/* User Stats */}
        <div className="bg-white rounded-2xl p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">User Growth</h2>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-gray-600">New users today</span>
              <span className="font-semibold text-gray-900">
                {userStats?.new_users_today || 0}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600">New users this week</span>
              <span className="font-semibold text-gray-900">
                {userStats?.new_users_week || 0}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600">New users this month</span>
              <span className="font-semibold text-gray-900">
                {userStats?.new_users_month || 0}
              </span>
            </div>
            <hr className="my-2" />
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Active this week</span>
              <span className="font-semibold text-gray-900">
                {overview?.active_users_week || 0}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Active this month</span>
              <span className="font-semibold text-gray-900">
                {overview?.active_users_month || 0}
              </span>
            </div>
          </div>
        </div>

        {/* Content Stats */}
        <div className="bg-white rounded-2xl p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Content Overview</h2>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Active languages</span>
              <span className="font-semibold text-gray-900">
                {contentStats?.active_languages || 0}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Total units</span>
              <span className="font-semibold text-gray-900">
                {contentStats?.total_units || 0}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Total lessons</span>
              <span className="font-semibold text-gray-900">
                {contentStats?.total_lessons || 0}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Total questions</span>
              <span className="font-semibold text-gray-900">
                {contentStats?.total_questions || 0}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Questions with audio</span>
              <span className="font-semibold text-gray-900">
                {contentStats?.questions_with_audio || 0}
              </span>
            </div>
          </div>
        </div>

        {/* Auth Providers */}
        <div className="bg-white rounded-2xl p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Sign-up Methods</h2>
          {userStats?.users_by_provider && Object.keys(userStats.users_by_provider).length > 0 ? (
            <div className="space-y-3">
              {Object.entries(userStats.users_by_provider).map(([provider, count]) => (
                <div key={provider} className="flex items-center justify-between">
                  <span className="text-gray-600 capitalize">{provider}</span>
                  <span className="font-semibold text-gray-900">{count}</span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500">No data available</p>
          )}
        </div>

        {/* Activity Stats */}
        <div className="bg-white rounded-2xl p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Today&apos;s Activity</h2>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Lessons completed</span>
              <span className="font-semibold text-gray-900">
                {overview?.lessons_completed_today || 0}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Total enrollments</span>
              <span className="font-semibold text-gray-900">
                {overview?.total_enrollments || 0}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
