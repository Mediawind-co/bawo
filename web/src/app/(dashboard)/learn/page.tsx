"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { createClient } from "@/lib/api";
import {
  FireIcon,
  TargetIcon,
  ChevronRightIcon,
  PlusIcon,
  BookIcon,
  CheckIcon,
  LoadingSpinner,
  GlobeIcon,
} from "@/components/icons";
import type { enrollment, language, streak, content } from "@/lib/client";

export default function LearnPage() {
  const { token, user } = useAuth();
  const [enrollments, setEnrollments] = useState<enrollment.EnrollmentWithLanguage[]>([]);
  const [languages, setLanguages] = useState<language.Language[]>([]);
  const [streakData, setStreakData] = useState<streak.StreakResponse | null>(null);
  const [selectedLanguage, setSelectedLanguage] = useState<string | null>(null);
  const [units, setUnits] = useState<content.Unit[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showLanguageModal, setShowLanguageModal] = useState(false);

  useEffect(() => {
    if (token) {
      loadData();
    }
  }, [token]);

  const loadData = async () => {
    if (!token) return;

    const client = createClient(token);
    try {
      const [enrollmentsRes, streakRes, languagesRes] = await Promise.all([
        client.enrollment.ListEnrollments(),
        client.streak.GetMyStreak(),
        client.language.ListLanguages(),
      ]);

      setEnrollments(enrollmentsRes.enrollments);
      setStreakData(streakRes);
      setLanguages(languagesRes.languages);

      // If user has enrollments, select the first one
      if (enrollmentsRes.enrollments.length > 0) {
        // We need to get the language ID - fetch all languages to match
        const firstEnrollment = enrollmentsRes.enrollments[0];
        const matchedLang = languagesRes.languages.find(
          (l) => l.code === firstEnrollment.language_code
        );
        if (matchedLang) {
          setSelectedLanguage(matchedLang.id);
          await loadUnits(matchedLang.id);
        }
      }
    } catch (error) {
      console.error("Failed to load data:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const loadUnits = async (languageId: string) => {
    if (!token) return;

    const client = createClient(token);
    try {
      const unitsRes = await client.content.ListUnits(languageId);
      setUnits(unitsRes.units);
    } catch (error) {
      console.error("Failed to load units:", error);
    }
  };

  const handleEnroll = async (languageId: string) => {
    if (!token) return;

    const client = createClient(token);
    try {
      await client.enrollment.Enroll({ language_id: languageId });
      setShowLanguageModal(false);
      await loadData();
    } catch (error) {
      console.error("Failed to enroll:", error);
    }
  };

  const handleSelectLanguage = async (langCode: string) => {
    const lang = languages.find((l) => l.code === langCode);
    if (lang) {
      setSelectedLanguage(lang.id);
      await loadUnits(lang.id);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner className="w-8 h-8 text-teal" />
      </div>
    );
  }

  const availableLanguages = languages.filter(
    (l) => !enrollments.some((e) => e.language_code === l.code)
  );

  return (
    <div className="p-4 lg:p-8 pb-24 lg:pb-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-2xl lg:text-3xl font-bold text-teal">
          Welcome back, {user?.name?.split(" ")[0] || "Learner"}!
        </h1>
        <p className="mt-1 text-gray-600">Continue your language journey</p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div className="bg-white rounded-2xl p-4 lg:p-6 shadow-sm border border-gray-100">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-orange-100 rounded-xl flex items-center justify-center">
              <FireIcon className="w-5 h-5 text-orange-500" />
            </div>
            <div>
              <p className="text-2xl font-bold text-teal">
                {streakData?.streak.current_streak || 0}
              </p>
              <p className="text-sm text-gray-500">Day Streak</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-2xl p-4 lg:p-6 shadow-sm border border-gray-100">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-mint/20 rounded-xl flex items-center justify-center">
              <TargetIcon className="w-5 h-5 text-teal" />
            </div>
            <div>
              <p className="text-2xl font-bold text-teal">
                {streakData?.today_xp || 0}/{streakData?.streak.daily_xp_goal || 50}
              </p>
              <p className="text-sm text-gray-500">Daily XP</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-2xl p-4 lg:p-6 shadow-sm border border-gray-100">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-purple/20 rounded-xl flex items-center justify-center">
              <BookIcon className="w-5 h-5 text-purple" />
            </div>
            <div>
              <p className="text-2xl font-bold text-teal">
                {streakData?.today_lessons || 0}
              </p>
              <p className="text-sm text-gray-500">Lessons Today</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-2xl p-4 lg:p-6 shadow-sm border border-gray-100">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-indigo/20 rounded-xl flex items-center justify-center">
              <GlobeIcon className="w-5 h-5 text-indigo" />
            </div>
            <div>
              <p className="text-2xl font-bold text-teal">{enrollments.length}</p>
              <p className="text-sm text-gray-500">Languages</p>
            </div>
          </div>
        </div>
      </div>

      {/* Daily Goal Progress */}
      {streakData && (
        <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-teal">Daily Goal</h2>
            {streakData.goal_met ? (
              <span className="flex items-center gap-1 text-sm text-green-600 bg-green-100 px-3 py-1 rounded-full">
                <CheckIcon className="w-4 h-4" />
                Complete!
              </span>
            ) : (
              <span className="text-sm text-gray-500">
                {streakData.xp_to_goal} XP to go
              </span>
            )}
          </div>
          <div className="w-full bg-gray-200 rounded-full h-3">
            <div
              className="bg-gradient-to-r from-mint to-purple h-3 rounded-full transition-all duration-500"
              style={{
                width: `${Math.min(
                  (streakData.today_xp / streakData.streak.daily_xp_goal) * 100,
                  100
                )}%`,
              }}
            />
          </div>
        </div>
      )}

      {/* Language Selection */}
      {enrollments.length > 0 ? (
        <div className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-teal">Your Languages</h2>
            {availableLanguages.length > 0 && (
              <button
                onClick={() => setShowLanguageModal(true)}
                className="flex items-center gap-1 text-sm text-purple hover:text-purple/80 transition-colors"
              >
                <PlusIcon className="w-4 h-4" />
                Add Language
              </button>
            )}
          </div>
          <div className="flex gap-3 overflow-x-auto pb-2">
            {enrollments.map((enrollment) => {
              const lang = languages.find((l) => l.code === enrollment.language_code);
              const isSelected = lang?.id === selectedLanguage;
              return (
                <button
                  key={enrollment.language_code}
                  onClick={() => handleSelectLanguage(enrollment.language_code)}
                  className={`flex items-center gap-2 px-4 py-2 rounded-xl whitespace-nowrap transition-colors ${
                    isSelected
                      ? "bg-mint text-teal"
                      : "bg-white border border-gray-200 text-gray-700 hover:border-mint"
                  }`}
                >
                  <span className="text-xl">{enrollment.language_emoji}</span>
                  <span className="font-medium">{enrollment.language_name}</span>
                </button>
              );
            })}
          </div>
        </div>
      ) : (
        <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 text-center mb-8">
          <GlobeIcon className="w-12 h-12 text-gray-300 mx-auto mb-4" />
          <h2 className="text-xl font-semibold text-teal mb-2">
            Start Your Journey
          </h2>
          <p className="text-gray-600 mb-6">
            Choose a language to begin learning
          </p>
          <button
            onClick={() => setShowLanguageModal(true)}
            className="bg-mint hover:bg-primary-dark text-teal font-semibold px-6 py-3 rounded-xl transition-colors"
          >
            Choose a Language
          </button>
        </div>
      )}

      {/* Units List */}
      {selectedLanguage && units.length > 0 && (
        <div>
          <h2 className="text-lg font-semibold text-teal mb-4">Units</h2>
          <div className="space-y-4">
            {units.map((unit, index) => (
              <Link
                key={unit.id}
                href={`/learn/${selectedLanguage}/unit/${unit.id}`}
                className="block bg-white rounded-2xl p-6 shadow-sm border border-gray-100 hover:border-mint transition-colors"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-mint to-purple rounded-xl flex items-center justify-center text-white font-bold">
                      {index + 1}
                    </div>
                    <div>
                      <h3 className="font-semibold text-teal">{unit.title}</h3>
                      <p className="text-sm text-gray-500">{unit.description}</p>
                    </div>
                  </div>
                  <ChevronRightIcon className="w-5 h-5 text-gray-400" />
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}

      {selectedLanguage && units.length === 0 && (
        <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 text-center">
          <BookIcon className="w-12 h-12 text-gray-300 mx-auto mb-4" />
          <h2 className="text-xl font-semibold text-teal mb-2">
            Coming Soon
          </h2>
          <p className="text-gray-600">
            Content for this language is being prepared. Check back soon!
          </p>
        </div>
      )}

      {/* Language Selection Modal */}
      {showLanguageModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md max-h-[80vh] overflow-y-auto">
            <h2 className="text-xl font-semibold text-teal mb-4">
              Choose a Language
            </h2>
            <div className="space-y-3">
              {availableLanguages.map((lang) => (
                <button
                  key={lang.id}
                  onClick={() => handleEnroll(lang.id)}
                  className="w-full flex items-center gap-4 p-4 rounded-xl border border-gray-200 hover:border-mint hover:bg-mint/5 transition-colors text-left"
                >
                  <span className="text-3xl">{lang.flag_emoji}</span>
                  <div>
                    <p className="font-semibold text-teal">{lang.name}</p>
                    <p className="text-sm text-gray-500">{lang.description}</p>
                  </div>
                </button>
              ))}
            </div>
            <button
              onClick={() => setShowLanguageModal(false)}
              className="mt-6 w-full py-3 text-gray-600 hover:bg-gray-100 rounded-xl transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
