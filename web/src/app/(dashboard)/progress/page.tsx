"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import { createClient } from "@/lib/api";
import {
  FireIcon,
  TrophyIcon,
  BookIcon,
  TargetIcon,
  CheckIcon,
  XIcon,
  LoadingSpinner,
} from "@/components/icons";
import type { tracker, streak } from "@/lib/client";

export default function ProgressPage() {
  const { token } = useAuth();
  const [stats, setStats] = useState<tracker.UserStatsResponse | null>(null);
  const [streakData, setStreakData] = useState<streak.StreakResponse | null>(null);
  const [activities, setActivities] = useState<streak.DailyActivity[]>([]);
  const [mistakes, setMistakes] = useState<tracker.MistakeWithQuestion[]>([]);
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
      const [statsRes, streakRes, activityRes, mistakesRes] = await Promise.all([
        client.tracker.GetMyStats(),
        client.streak.GetMyStreak(),
        client.streak.GetActivityHistory(),
        client.tracker.ListMyMistakes({ Limit: 10, Offset: 0 }),
      ]);

      setStats(statsRes);
      setStreakData(streakRes);
      setActivities(activityRes.activities);
      setMistakes(mistakesRes.mistakes);
    } catch (error) {
      console.error("Failed to load progress data:", error);
    } finally {
      setIsLoading(false);
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
      <h1 className="text-2xl lg:text-3xl font-bold text-teal mb-8">Your Progress</h1>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
          <div className="w-10 h-10 bg-purple/20 rounded-xl flex items-center justify-center mb-3">
            <TrophyIcon className="w-5 h-5 text-purple" />
          </div>
          <p className="text-2xl font-bold text-teal">{stats?.total_xp || 0}</p>
          <p className="text-sm text-gray-500">Total XP</p>
        </div>

        <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
          <div className="w-10 h-10 bg-mint/20 rounded-xl flex items-center justify-center mb-3">
            <BookIcon className="w-5 h-5 text-teal" />
          </div>
          <p className="text-2xl font-bold text-teal">{stats?.lessons_completed || 0}</p>
          <p className="text-sm text-gray-500">Lessons Done</p>
        </div>

        <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
          <div className="w-10 h-10 bg-orange-100 rounded-xl flex items-center justify-center mb-3">
            <FireIcon className="w-5 h-5 text-orange-500" />
          </div>
          <p className="text-2xl font-bold text-teal">
            {streakData?.streak.longest_streak || 0}
          </p>
          <p className="text-sm text-gray-500">Best Streak</p>
        </div>

        <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
          <div className="w-10 h-10 bg-green-100 rounded-xl flex items-center justify-center mb-3">
            <TargetIcon className="w-5 h-5 text-green-600" />
          </div>
          <p className="text-2xl font-bold text-teal">
            {stats?.average_score ? `${Math.round(stats.average_score)}%` : "N/A"}
          </p>
          <p className="text-sm text-gray-500">Avg. Score</p>
        </div>
      </div>

      {/* Activity Calendar */}
      <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 mb-8">
        <h2 className="text-lg font-semibold text-teal mb-4">Activity (Last 30 Days)</h2>
        <div className="grid grid-cols-10 gap-2">
          {activities.length > 0 ? (
            activities.slice(0, 30).map((activity, index) => (
              <div
                key={index}
                className={`w-full aspect-square rounded-md ${
                  activity.goal_met
                    ? "bg-mint"
                    : activity.xp_earned > 0
                    ? "bg-mint/40"
                    : "bg-gray-100"
                }`}
                title={`${new Date(activity.date).toLocaleDateString()}: ${activity.xp_earned} XP`}
              />
            ))
          ) : (
            Array.from({ length: 30 }).map((_, index) => (
              <div key={index} className="w-full aspect-square rounded-md bg-gray-100" />
            ))
          )}
        </div>
        <div className="flex items-center gap-4 mt-4 text-sm text-gray-500">
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 rounded bg-gray-100" />
            <span>No activity</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 rounded bg-mint/40" />
            <span>Some XP</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 rounded bg-mint" />
            <span>Goal met</span>
          </div>
        </div>
      </div>

      {/* Streak Info */}
      {streakData && (
        <div className="bg-gradient-to-r from-orange-400 to-orange-500 rounded-2xl p-6 text-white mb-8">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-orange-100 text-sm">Current Streak</p>
              <p className="text-4xl font-bold mt-1">
                {streakData.streak.current_streak} days
              </p>
            </div>
            <FireIcon className="w-16 h-16 text-white/20" />
          </div>
          {streakData.streak_at_risk && (
            <div className="mt-4 bg-white/20 rounded-lg p-3">
              <p className="text-sm">
                Your streak is at risk! Complete your daily goal to keep it going.
              </p>
            </div>
          )}
          {streakData.streak.streak_freezes > 0 && (
            <p className="mt-4 text-sm text-orange-100">
              You have {streakData.streak.streak_freezes} streak freeze(s) available
            </p>
          )}
        </div>
      )}

      {/* Mistakes to Review */}
      <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-teal">Review Mistakes</h2>
          {stats && stats.unreviewed_mistakes > 0 && (
            <span className="bg-red-100 text-red-600 text-sm px-3 py-1 rounded-full">
              {stats.unreviewed_mistakes} to review
            </span>
          )}
        </div>
        {mistakes.length > 0 ? (
          <div className="space-y-3">
            {mistakes.slice(0, 5).map((mistake, index) => (
              <div
                key={index}
                className="flex items-start gap-4 p-4 bg-gray-50 rounded-xl"
              >
                <div className="w-8 h-8 bg-red-100 rounded-lg flex items-center justify-center flex-shrink-0">
                  <XIcon className="w-4 h-4 text-red-500" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-gray-700">{mistake.prompt_text}</p>
                  {mistake.hint && (
                    <p className="text-sm text-gray-500 mt-1">Hint: {mistake.hint}</p>
                  )}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <CheckIcon className="w-12 h-12 text-green-400 mx-auto mb-3" />
            <p className="text-gray-600">No mistakes to review. Great job!</p>
          </div>
        )}
      </div>
    </div>
  );
}
