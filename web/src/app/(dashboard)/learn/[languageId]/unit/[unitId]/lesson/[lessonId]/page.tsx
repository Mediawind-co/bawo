"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { createClient } from "@/lib/api";
import {
  CheckIcon,
  XIcon,
  LoadingSpinner,
  VolumeIcon,
  StarIcon,
} from "@/components/icons";
import type { lesson } from "@/lib/client";

export default function LessonPage() {
  const params = useParams();
  const router = useRouter();
  const searchParams = useSearchParams();
  const { token } = useAuth();

  const languageId = params.languageId as string;
  const unitId = params.unitId as string;
  const lessonId = params.lessonId as string;
  const sessionId = searchParams.get("session");

  const [session, setSession] = useState<lesson.LessonSession | null>(null);
  const [questions, setQuestions] = useState<lesson.QuestionInfo[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [selectedAnswer, setSelectedAnswer] = useState<string>("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [feedback, setFeedback] = useState<{
    isCorrect: boolean;
    correctAnswer: string;
    hint: string;
    xpEarned: number;
  } | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [completionData, setCompletionData] = useState<lesson.CompleteSessionResponse | null>(null);

  useEffect(() => {
    if (token && lessonId) {
      loadData();
    }
  }, [token, lessonId, sessionId]);

  const loadData = async () => {
    if (!token) return;

    const client = createClient(token);
    try {
      // Start or resume the session
      const sessionRes = await client.lesson.StartLesson(lessonId);
      setSession(sessionRes.session);
      setQuestions(sessionRes.questions);
      setCurrentIndex(sessionRes.session.current_index);
    } catch (error) {
      console.error("Failed to load lesson:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmitAnswer = async () => {
    if (!token || !session || !selectedAnswer || isSubmitting) return;

    setIsSubmitting(true);
    const client = createClient(token);

    try {
      const currentQuestion = questions[currentIndex];
      const response = await client.lesson.SubmitAnswer(session.id, {
        question_id: currentQuestion.id,
        answer: selectedAnswer,
      });

      setFeedback({
        isCorrect: response.is_correct,
        correctAnswer: response.correct_answer,
        hint: response.hint,
        xpEarned: response.xp_earned,
      });
    } catch (error) {
      console.error("Failed to submit answer:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleContinue = async () => {
    if (!token || !session) return;

    const client = createClient(token);

    // Check if this was the last question
    if (currentIndex >= questions.length - 1) {
      // Complete the lesson
      try {
        const response = await client.lesson.CompleteLesson(session.id);
        setCompletionData(response);
      } catch (error) {
        console.error("Failed to complete lesson:", error);
      }
    } else {
      // Move to next question
      setCurrentIndex((prev) => prev + 1);
      setSelectedAnswer("");
      setFeedback(null);
    }
  };

  const handleExit = () => {
    router.push(`/learn/${languageId}/unit/${unitId}`);
  };

  const playAudio = async (questionId: string) => {
    if (!token) return;

    const client = createClient(token);
    try {
      const response = await client.content.GetQuestionAudioURL(questionId);
      const audio = new Audio(response.url);
      audio.play();
    } catch (error) {
      console.error("Failed to play audio:", error);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner className="w-8 h-8 text-teal" />
      </div>
    );
  }

  // Completion screen
  if (completionData) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4">
        <div className="bg-white rounded-2xl p-8 shadow-lg border border-gray-100 text-center max-w-md w-full">
          <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <CheckIcon className="w-10 h-10 text-green-600" />
          </div>
          <h1 className="text-2xl font-bold text-teal mb-2">Lesson Complete!</h1>
          <p className="text-gray-600 mb-6">
            Great job completing this lesson!
          </p>

          <div className="grid grid-cols-2 gap-4 mb-8">
            <div className="bg-gray-50 rounded-xl p-4">
              <div className="flex items-center justify-center gap-2 mb-1">
                <StarIcon className="w-5 h-5 text-yellow-500" />
                <span className="text-2xl font-bold text-teal">
                  {completionData.total_xp}
                </span>
              </div>
              <p className="text-sm text-gray-500">XP Earned</p>
            </div>
            <div className="bg-gray-50 rounded-xl p-4">
              <div className="text-2xl font-bold text-teal mb-1">
                {Math.round(completionData.accuracy)}%
              </div>
              <p className="text-sm text-gray-500">Accuracy</p>
            </div>
          </div>

          <div className="text-sm text-gray-500 mb-6">
            {completionData.correct_count} / {completionData.total_count} correct
          </div>

          <button
            onClick={handleExit}
            className="w-full bg-mint hover:bg-primary-dark text-teal font-semibold py-3 rounded-xl transition-colors"
          >
            Continue
          </button>
        </div>
      </div>
    );
  }

  const currentQuestion = questions[currentIndex];
  const progress = ((currentIndex + 1) / questions.length) * 100;

  if (!currentQuestion) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-gray-700">No questions available</h2>
          <button
            onClick={handleExit}
            className="mt-4 text-teal hover:underline"
          >
            Back to Unit
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 p-4">
        <div className="max-w-2xl mx-auto flex items-center gap-4">
          <button
            onClick={handleExit}
            className="p-2 text-gray-400 hover:text-gray-600 transition-colors"
          >
            <XIcon className="w-6 h-6" />
          </button>

          {/* Progress bar */}
          <div className="flex-1">
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div
                className="bg-gradient-to-r from-mint to-purple h-2 rounded-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
          </div>

          <span className="text-sm text-gray-500">
            {currentIndex + 1} / {questions.length}
          </span>
        </div>
      </div>

      {/* Question content */}
      <div className="flex-1 p-4 lg:p-8 pb-32">
        <div className="max-w-2xl mx-auto">
          {/* Question prompt */}
          <div className="mb-8">
            <h2 className="text-lg font-semibold text-gray-700 mb-4">
              {currentQuestion.type === "multiple_choice"
                ? "Select the correct answer"
                : currentQuestion.type === "translation"
                ? "Translate this phrase"
                : currentQuestion.type === "listening"
                ? "What do you hear?"
                : "Answer the question"}
            </h2>

            <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
              <div className="flex items-start gap-4">
                {currentQuestion.has_audio && (
                  <button
                    onClick={() => playAudio(currentQuestion.id)}
                    className="w-12 h-12 bg-mint rounded-xl flex items-center justify-center text-teal hover:bg-primary-dark transition-colors shrink-0"
                  >
                    <VolumeIcon className="w-6 h-6" />
                  </button>
                )}
                <p className="text-xl text-teal font-medium flex-1">
                  {currentQuestion.prompt_text}
                </p>
              </div>
            </div>
          </div>

          {/* Answer options */}
          {currentQuestion.options && currentQuestion.options.length > 0 ? (
            <div className="space-y-3 mb-8">
              {currentQuestion.options.map((option, index) => {
                const isSelected = selectedAnswer === option;
                const showFeedback = feedback !== null;
                const isCorrectAnswer = feedback?.correctAnswer === option;
                const isWrongSelection = showFeedback && isSelected && !feedback.isCorrect;

                return (
                  <button
                    key={index}
                    onClick={() => !showFeedback && setSelectedAnswer(option)}
                    disabled={showFeedback}
                    className={`w-full text-left p-4 rounded-xl border-2 transition-all ${
                      showFeedback
                        ? isCorrectAnswer
                          ? "border-green-500 bg-green-50"
                          : isWrongSelection
                          ? "border-red-500 bg-red-50"
                          : "border-gray-200 bg-white opacity-50"
                        : isSelected
                        ? "border-mint bg-mint/10"
                        : "border-gray-200 bg-white hover:border-gray-300"
                    }`}
                  >
                    <div className="flex items-center justify-between">
                      <span
                        className={`font-medium ${
                          showFeedback && isCorrectAnswer
                            ? "text-green-700"
                            : showFeedback && isWrongSelection
                            ? "text-red-700"
                            : "text-gray-700"
                        }`}
                      >
                        {option}
                      </span>
                      {showFeedback && isCorrectAnswer && (
                        <CheckIcon className="w-5 h-5 text-green-600" />
                      )}
                      {isWrongSelection && (
                        <XIcon className="w-5 h-5 text-red-600" />
                      )}
                    </div>
                  </button>
                );
              })}
            </div>
          ) : (
            // Text input for non-multiple choice questions
            <div className="mb-8">
              <input
                type="text"
                value={selectedAnswer}
                onChange={(e) => !feedback && setSelectedAnswer(e.target.value)}
                disabled={feedback !== null}
                placeholder="Type your answer..."
                className="w-full p-4 rounded-xl border-2 border-gray-200 focus:border-mint focus:ring-0 outline-none text-lg"
              />
            </div>
          )}

          {/* Feedback display */}
          {feedback && (
            <div
              className={`p-4 rounded-xl mb-8 ${
                feedback.isCorrect
                  ? "bg-green-100 border border-green-200"
                  : "bg-red-100 border border-red-200"
              }`}
            >
              <div className="flex items-center gap-2 mb-2">
                {feedback.isCorrect ? (
                  <>
                    <CheckIcon className="w-5 h-5 text-green-600" />
                    <span className="font-semibold text-green-700">Correct!</span>
                    {feedback.xpEarned > 0 && (
                      <span className="ml-auto text-sm text-green-600">
                        +{feedback.xpEarned} XP
                      </span>
                    )}
                  </>
                ) : (
                  <>
                    <XIcon className="w-5 h-5 text-red-600" />
                    <span className="font-semibold text-red-700">Not quite</span>
                  </>
                )}
              </div>
              {!feedback.isCorrect && feedback.hint && (
                <p className="text-sm text-red-600">{feedback.hint}</p>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Bottom action bar */}
      <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 p-4 lg:pl-64">
        <div className="max-w-2xl mx-auto">
          {feedback ? (
            <button
              onClick={handleContinue}
              className={`w-full py-4 rounded-xl font-semibold transition-colors ${
                feedback.isCorrect
                  ? "bg-green-500 hover:bg-green-600 text-white"
                  : "bg-red-500 hover:bg-red-600 text-white"
              }`}
            >
              Continue
            </button>
          ) : (
            <button
              onClick={handleSubmitAnswer}
              disabled={!selectedAnswer || isSubmitting}
              className={`w-full py-4 rounded-xl font-semibold transition-colors ${
                selectedAnswer
                  ? "bg-mint hover:bg-primary-dark text-teal"
                  : "bg-gray-200 text-gray-400 cursor-not-allowed"
              }`}
            >
              {isSubmitting ? (
                <LoadingSpinner className="w-5 h-5 mx-auto" />
              ) : (
                "Check"
              )}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
