"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { createClient } from "@/lib/api";
import {
  ArrowLeftIcon,
  BookIcon,
  CheckIcon,
  LoadingSpinner,
  PlayIcon,
  StarIcon,
} from "@/components/icons";
import type { content, language, tracker } from "@/lib/client";

export default function UnitPage() {
  const params = useParams();
  const router = useRouter();
  const { token } = useAuth();
  const [unit, setUnit] = useState<content.Unit | null>(null);
  const [lessons, setLessons] = useState<content.Lesson[]>([]);
  const [lang, setLang] = useState<language.Language | null>(null);
  const [progress, setProgress] = useState<tracker.LessonProgress[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [startingLesson, setStartingLesson] = useState<string | null>(null);

  const languageId = params.languageId as string;
  const unitId = params.unitId as string;

  useEffect(() => {
    if (token && unitId && languageId) {
      loadData();
    }
  }, [token, unitId, languageId]);

  const loadData = async () => {
    if (!token) return;

    const client = createClient(token);
    try {
      const [unitRes, lessonsRes, langRes, progressRes] = await Promise.all([
        client.content.GetUnit(unitId),
        client.content.ListLessons(unitId),
        client.language.GetLanguage(languageId),
        client.tracker.ListMyProgress(),
      ]);

      setUnit(unitRes.unit);
      setLessons(lessonsRes.lessons);
      setLang(langRes.language);
      setProgress(progressRes.progress);
    } catch (error) {
      console.error("Failed to load data:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const startLesson = async (lessonId: string) => {
    if (!token || startingLesson) return;

    setStartingLesson(lessonId);
    const client = createClient(token);
    try {
      const response = await client.lesson.StartLesson(lessonId);
      router.push(`/learn/${languageId}/unit/${unitId}/lesson/${lessonId}?session=${response.session.id}`);
    } catch (error) {
      console.error("Failed to start lesson:", error);
      setStartingLesson(null);
    }
  };

  const getLessonProgress = (lessonId: string) => {
    return progress.find((p) => p.lesson_id === lessonId);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner className="w-8 h-8 text-teal" />
      </div>
    );
  }

  if (!unit) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-gray-700">Unit not found</h2>
          <Link href="/learn" className="mt-4 text-teal hover:underline">
            Back to Learn
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="p-4 lg:p-8 pb-24 lg:pb-8">
      {/* Back Navigation */}
      <Link
        href="/learn"
        className="inline-flex items-center gap-2 text-gray-600 hover:text-teal mb-6 transition-colors"
      >
        <ArrowLeftIcon className="w-5 h-5" />
        <span>Back to Languages</span>
      </Link>

      {/* Unit Header */}
      <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 mb-8">
        <div className="flex items-start gap-4">
          <div className="w-16 h-16 bg-gradient-to-br from-mint to-purple rounded-xl flex items-center justify-center text-white">
            <BookIcon className="w-8 h-8" />
          </div>
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-1">
              {lang && <span className="text-2xl">{lang.flag_emoji}</span>}
              <span className="text-sm text-gray-500">{lang?.name}</span>
            </div>
            <h1 className="text-2xl font-bold text-teal">{unit.title}</h1>
            <p className="mt-1 text-gray-600">{unit.description}</p>
          </div>
        </div>
      </div>

      {/* Lessons List */}
      <div className="mb-4">
        <h2 className="text-lg font-semibold text-teal">Lessons</h2>
        <p className="text-sm text-gray-500">{lessons.length} lessons in this unit</p>
      </div>

      {lessons.length > 0 ? (
        <div className="space-y-4">
          {lessons.map((lesson, index) => {
            const lessonProgress = getLessonProgress(lesson.id);
            const isCompleted = lessonProgress?.status === "completed";
            const isInProgress = lessonProgress?.status === "in_progress";
            const isStarting = startingLesson === lesson.id;

            return (
              <div
                key={lesson.id}
                className={`bg-white rounded-2xl p-6 shadow-sm border transition-colors ${
                  isCompleted
                    ? "border-green-200 bg-green-50/30"
                    : isInProgress
                    ? "border-mint"
                    : "border-gray-100 hover:border-mint"
                }`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div
                      className={`w-12 h-12 rounded-xl flex items-center justify-center font-bold ${
                        isCompleted
                          ? "bg-green-100 text-green-600"
                          : isInProgress
                          ? "bg-mint text-teal"
                          : "bg-gray-100 text-gray-600"
                      }`}
                    >
                      {isCompleted ? (
                        <CheckIcon className="w-6 h-6" />
                      ) : (
                        index + 1
                      )}
                    </div>
                    <div className="flex-1">
                      <h3 className="font-semibold text-teal">{lesson.title}</h3>
                      <p className="text-sm text-gray-500">{lesson.description}</p>
                      <div className="flex items-center gap-4 mt-2">
                        <span className="flex items-center gap-1 text-xs text-gray-400">
                          <StarIcon className="w-3 h-3" />
                          {lesson.xp_reward} XP
                        </span>
                        {lessonProgress?.best_score !== undefined && lessonProgress.best_score > 0 && (
                          <span className="text-xs text-gray-400">
                            Best: {lessonProgress.best_score}%
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                  <button
                    onClick={() => startLesson(lesson.id)}
                    disabled={isStarting}
                    className={`flex items-center gap-2 px-4 py-2 rounded-xl font-medium transition-colors ${
                      isCompleted
                        ? "bg-green-100 text-green-700 hover:bg-green-200"
                        : isInProgress
                        ? "bg-mint text-teal hover:bg-primary-dark"
                        : "bg-mint text-teal hover:bg-primary-dark"
                    }`}
                  >
                    {isStarting ? (
                      <LoadingSpinner className="w-4 h-4" />
                    ) : (
                      <>
                        <PlayIcon className="w-4 h-4" />
                        {isCompleted ? "Retry" : isInProgress ? "Continue" : "Start"}
                      </>
                    )}
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      ) : (
        <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 text-center">
          <BookIcon className="w-12 h-12 text-gray-300 mx-auto mb-4" />
          <h2 className="text-xl font-semibold text-teal mb-2">Coming Soon</h2>
          <p className="text-gray-600">
            Lessons for this unit are being prepared. Check back soon!
          </p>
        </div>
      )}
    </div>
  );
}
