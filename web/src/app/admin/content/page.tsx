"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import { createClient } from "@/lib/api";
import {
  PlusIcon,
  EditIcon,
  TrashIcon,
  ChevronRightIcon,
  ChevronDownIcon,
  XIcon,
  LoadingSpinner,
  VolumeIcon,
} from "@/components/icons";
import type { language, content } from "@/lib/client";

type ModalType = "unit" | "lesson" | "question" | null;

export default function ContentPage() {
  const { token } = useAuth();
  const [languages, setLanguages] = useState<language.Language[]>([]);
  const [selectedLanguage, setSelectedLanguage] = useState<string | null>(null);
  const [units, setUnits] = useState<content.Unit[]>([]);
  const [expandedUnit, setExpandedUnit] = useState<string | null>(null);
  const [lessons, setLessons] = useState<Record<string, content.Lesson[]>>({});
  const [expandedLesson, setExpandedLesson] = useState<string | null>(null);
  const [questions, setQuestions] = useState<Record<string, content.Question[]>>({});
  const [isLoading, setIsLoading] = useState(true);

  // Modal state
  const [modalType, setModalType] = useState<ModalType>(null);
  const [editingItem, setEditingItem] = useState<any>(null);
  const [parentId, setParentId] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  // Form data
  const [unitForm, setUnitForm] = useState({
    title: "",
    description: "",
    sort_order: 0,
  });
  const [lessonForm, setLessonForm] = useState({
    title: "",
    description: "",
    xp_reward: 10,
    sort_order: 0,
  });
  const [questionForm, setQuestionForm] = useState({
    type: "single_choice" as content.QuestionType,
    prompt_text: "",
    use_tts: false,
    correct_answer: "",
    options: ["", "", "", ""],
    hint: "",
    sort_order: 0,
  });

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
      if (res.languages.length > 0) {
        setSelectedLanguage(res.languages[0].id);
        await loadUnits(res.languages[0].id);
      }
    } catch (error) {
      console.error("Failed to load languages:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const loadUnits = async (langId: string) => {
    if (!token) return;

    const client = createClient(token);
    try {
      const res = await client.content.ListUnits(langId);
      setUnits(res.units);
    } catch (error) {
      console.error("Failed to load units:", error);
    }
  };

  const loadLessons = async (unitId: string) => {
    if (!token) return;

    const client = createClient(token);
    try {
      const res = await client.content.ListLessons(unitId);
      setLessons((prev) => ({ ...prev, [unitId]: res.lessons }));
    } catch (error) {
      console.error("Failed to load lessons:", error);
    }
  };

  const loadQuestions = async (lessonId: string) => {
    if (!token) return;

    const client = createClient(token);
    try {
      const res = await client.content.ListQuestions(lessonId);
      setQuestions((prev) => ({ ...prev, [lessonId]: res.questions }));
    } catch (error) {
      console.error("Failed to load questions:", error);
    }
  };

  const handleLanguageChange = async (langId: string) => {
    setSelectedLanguage(langId);
    setExpandedUnit(null);
    setExpandedLesson(null);
    await loadUnits(langId);
  };

  const toggleUnit = async (unitId: string) => {
    if (expandedUnit === unitId) {
      setExpandedUnit(null);
    } else {
      setExpandedUnit(unitId);
      if (!lessons[unitId]) {
        await loadLessons(unitId);
      }
    }
  };

  const toggleLesson = async (lessonId: string) => {
    if (expandedLesson === lessonId) {
      setExpandedLesson(null);
    } else {
      setExpandedLesson(lessonId);
      if (!questions[lessonId]) {
        await loadQuestions(lessonId);
      }
    }
  };

  // Modal handlers
  const openUnitModal = (unit?: content.Unit) => {
    setEditingItem(unit || null);
    setUnitForm({
      title: unit?.title || "",
      description: unit?.description || "",
      sort_order: unit?.sort_order || units.length,
    });
    setModalType("unit");
  };

  const openLessonModal = (unitId: string, lesson?: content.Lesson) => {
    setParentId(unitId);
    setEditingItem(lesson || null);
    const lessonCount = lessons[unitId]?.length || 0;
    setLessonForm({
      title: lesson?.title || "",
      description: lesson?.description || "",
      xp_reward: lesson?.xp_reward || 10,
      sort_order: lesson?.sort_order || lessonCount,
    });
    setModalType("lesson");
  };

  const openQuestionModal = (lessonId: string, question?: content.Question) => {
    setParentId(lessonId);
    setEditingItem(question || null);
    const questionCount = questions[lessonId]?.length || 0;
    setQuestionForm({
      type: question?.type || "single_choice",
      prompt_text: question?.prompt_text || "",
      use_tts: question?.use_tts || false,
      correct_answer: question?.correct_answer || "",
      options: question?.options?.length ? [...question.options] : ["", "", "", ""],
      hint: question?.hint || "",
      sort_order: question?.sort_order || questionCount,
    });
    setModalType("question");
  };

  const handleSaveUnit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || !selectedLanguage) return;

    setIsSaving(true);
    const client = createClient(token);

    try {
      if (editingItem) {
        await client.content.UpdateUnitAdmin(editingItem.id, unitForm);
      } else {
        await client.content.CreateUnitAdmin({
          language_id: selectedLanguage,
          ...unitForm,
        });
      }
      setModalType(null);
      await loadUnits(selectedLanguage);
    } catch (error) {
      console.error("Failed to save unit:", error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleSaveLesson = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || !parentId) return;

    setIsSaving(true);
    const client = createClient(token);

    try {
      if (editingItem) {
        await client.content.UpdateLessonAdmin(editingItem.id, lessonForm);
      } else {
        await client.content.CreateLessonAdmin({
          unit_id: parentId,
          ...lessonForm,
        });
      }
      setModalType(null);
      await loadLessons(parentId);
    } catch (error) {
      console.error("Failed to save lesson:", error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleSaveQuestion = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || !parentId) return;

    setIsSaving(true);
    const client = createClient(token);

    try {
      const filteredOptions = questionForm.options.filter((o) => o.trim());
      if (editingItem) {
        await client.content.UpdateQuestionAdmin(editingItem.id, {
          ...questionForm,
          options: filteredOptions,
        });
      } else {
        await client.content.CreateQuestionAdmin({
          lesson_id: parentId,
          ...questionForm,
          options: filteredOptions,
        });
      }
      setModalType(null);
      await loadQuestions(parentId);
    } catch (error) {
      console.error("Failed to save question:", error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleDeleteUnit = async (unitId: string) => {
    if (!token || !selectedLanguage) return;
    if (!confirm("Delete this unit and all its lessons?")) return;

    const client = createClient(token);
    try {
      await client.content.DeleteUnitAdmin(unitId);
      await loadUnits(selectedLanguage);
    } catch (error) {
      console.error("Failed to delete unit:", error);
    }
  };

  const handleDeleteLesson = async (lessonId: string, unitId: string) => {
    if (!token) return;
    if (!confirm("Delete this lesson and all its questions?")) return;

    const client = createClient(token);
    try {
      await client.content.DeleteLessonAdmin(lessonId);
      await loadLessons(unitId);
    } catch (error) {
      console.error("Failed to delete lesson:", error);
    }
  };

  const handleDeleteQuestion = async (questionId: string, lessonId: string) => {
    if (!token) return;
    if (!confirm("Delete this question?")) return;

    const client = createClient(token);
    try {
      await client.content.DeleteQuestionAdmin(questionId);
      await loadQuestions(lessonId);
    } catch (error) {
      console.error("Failed to delete question:", error);
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
        <h1 className="text-2xl lg:text-3xl font-bold text-gray-900">Content</h1>
      </div>

      {/* Language Selector */}
      <div className="mb-6">
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Select Language
        </label>
        <select
          value={selectedLanguage || ""}
          onChange={(e) => handleLanguageChange(e.target.value)}
          className="px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
        >
          {languages.map((lang) => (
            <option key={lang.id} value={lang.id}>
              {lang.flag_emoji} {lang.name}
            </option>
          ))}
        </select>
      </div>

      {/* Units List */}
      <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
        <div className="flex items-center justify-between p-4 border-b border-gray-100">
          <h2 className="font-semibold text-gray-900">Units</h2>
          <button
            onClick={() => openUnitModal()}
            className="flex items-center gap-1 text-sm text-purple hover:text-purple/80 transition-colors"
          >
            <PlusIcon className="w-4 h-4" />
            Add Unit
          </button>
        </div>

        {units.length === 0 ? (
          <div className="p-8 text-center text-gray-500">
            No units yet. Create your first unit!
          </div>
        ) : (
          <div className="divide-y divide-gray-100">
            {units.map((unit) => (
              <div key={unit.id}>
                {/* Unit Row */}
                <div
                  className="flex items-center justify-between p-4 hover:bg-gray-50 cursor-pointer"
                  onClick={() => toggleUnit(unit.id)}
                >
                  <div className="flex items-center gap-3">
                    {expandedUnit === unit.id ? (
                      <ChevronDownIcon className="w-5 h-5 text-gray-400" />
                    ) : (
                      <ChevronRightIcon className="w-5 h-5 text-gray-400" />
                    )}
                    <div>
                      <p className="font-medium text-gray-900">{unit.title}</p>
                      <p className="text-sm text-gray-500">{unit.description}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2" onClick={(e) => e.stopPropagation()}>
                    <button
                      onClick={() => openUnitModal(unit)}
                      className="p-1 text-gray-400 hover:text-purple"
                    >
                      <EditIcon className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDeleteUnit(unit.id)}
                      className="p-1 text-gray-400 hover:text-red-600"
                    >
                      <TrashIcon className="w-4 h-4" />
                    </button>
                  </div>
                </div>

                {/* Lessons */}
                {expandedUnit === unit.id && (
                  <div className="bg-gray-50 border-t border-gray-100">
                    <div className="flex items-center justify-between px-4 py-2 pl-12">
                      <span className="text-sm font-medium text-gray-600">Lessons</span>
                      <button
                        onClick={() => openLessonModal(unit.id)}
                        className="flex items-center gap-1 text-xs text-purple hover:text-purple/80"
                      >
                        <PlusIcon className="w-3 h-3" />
                        Add Lesson
                      </button>
                    </div>
                    {lessons[unit.id]?.length === 0 && (
                      <div className="px-4 py-4 pl-12 text-sm text-gray-500">
                        No lessons yet
                      </div>
                    )}
                    {lessons[unit.id]?.map((lesson) => (
                      <div key={lesson.id}>
                        <div
                          className="flex items-center justify-between px-4 py-3 pl-12 hover:bg-gray-100 cursor-pointer"
                          onClick={() => toggleLesson(lesson.id)}
                        >
                          <div className="flex items-center gap-3">
                            {expandedLesson === lesson.id ? (
                              <ChevronDownIcon className="w-4 h-4 text-gray-400" />
                            ) : (
                              <ChevronRightIcon className="w-4 h-4 text-gray-400" />
                            )}
                            <div>
                              <p className="text-sm font-medium text-gray-800">{lesson.title}</p>
                              <p className="text-xs text-gray-500">{lesson.xp_reward} XP</p>
                            </div>
                          </div>
                          <div className="flex items-center gap-2" onClick={(e) => e.stopPropagation()}>
                            <button
                              onClick={() => openLessonModal(unit.id, lesson)}
                              className="p-1 text-gray-400 hover:text-purple"
                            >
                              <EditIcon className="w-3 h-3" />
                            </button>
                            <button
                              onClick={() => handleDeleteLesson(lesson.id, unit.id)}
                              className="p-1 text-gray-400 hover:text-red-600"
                            >
                              <TrashIcon className="w-3 h-3" />
                            </button>
                          </div>
                        </div>

                        {/* Questions */}
                        {expandedLesson === lesson.id && (
                          <div className="bg-white border-t border-gray-100">
                            <div className="flex items-center justify-between px-4 py-2 pl-20">
                              <span className="text-xs font-medium text-gray-500">Questions</span>
                              <button
                                onClick={() => openQuestionModal(lesson.id)}
                                className="flex items-center gap-1 text-xs text-purple hover:text-purple/80"
                              >
                                <PlusIcon className="w-3 h-3" />
                                Add
                              </button>
                            </div>
                            {questions[lesson.id]?.length === 0 && (
                              <div className="px-4 py-3 pl-20 text-xs text-gray-500">
                                No questions yet
                              </div>
                            )}
                            {questions[lesson.id]?.map((q, idx) => (
                              <div
                                key={q.id}
                                className="flex items-center justify-between px-4 py-2 pl-20 hover:bg-gray-50"
                              >
                                <div className="flex items-center gap-2">
                                  <span className="text-xs text-gray-400">{idx + 1}.</span>
                                  <span className="text-sm text-gray-700 truncate max-w-xs">
                                    {q.prompt_text}
                                  </span>
                                  {q.prompt_audio_key && (
                                    <VolumeIcon className="w-3 h-3 text-purple" />
                                  )}
                                  <span className="text-xs text-gray-400 bg-gray-100 px-2 py-0.5 rounded">
                                    {q.type}
                                  </span>
                                </div>
                                <div className="flex items-center gap-1">
                                  <button
                                    onClick={() => openQuestionModal(lesson.id, q)}
                                    className="p-1 text-gray-400 hover:text-purple"
                                  >
                                    <EditIcon className="w-3 h-3" />
                                  </button>
                                  <button
                                    onClick={() => handleDeleteQuestion(q.id, lesson.id)}
                                    className="p-1 text-gray-400 hover:text-red-600"
                                  >
                                    <TrashIcon className="w-3 h-3" />
                                  </button>
                                </div>
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Unit Modal */}
      {modalType === "unit" && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-gray-900">
                {editingItem ? "Edit Unit" : "Add Unit"}
              </h2>
              <button onClick={() => setModalType(null)} className="text-gray-400 hover:text-gray-600">
                <XIcon className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleSaveUnit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Title</label>
                <input
                  type="text"
                  value={unitForm.title}
                  onChange={(e) => setUnitForm({ ...unitForm, title: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                <textarea
                  value={unitForm.description}
                  onChange={(e) => setUnitForm({ ...unitForm, description: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  rows={3}
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button type="button" onClick={() => setModalType(null)} className="flex-1 py-2 text-gray-600 hover:bg-gray-100 rounded-xl">
                  Cancel
                </button>
                <button type="submit" disabled={isSaving} className="flex-1 py-2 bg-purple hover:bg-purple/90 text-white font-medium rounded-xl disabled:opacity-50">
                  {isSaving ? <LoadingSpinner className="w-5 h-5 mx-auto" /> : "Save"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Lesson Modal */}
      {modalType === "lesson" && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-gray-900">
                {editingItem ? "Edit Lesson" : "Add Lesson"}
              </h2>
              <button onClick={() => setModalType(null)} className="text-gray-400 hover:text-gray-600">
                <XIcon className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleSaveLesson} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Title</label>
                <input
                  type="text"
                  value={lessonForm.title}
                  onChange={(e) => setLessonForm({ ...lessonForm, title: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                <textarea
                  value={lessonForm.description}
                  onChange={(e) => setLessonForm({ ...lessonForm, description: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  rows={2}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">XP Reward</label>
                <input
                  type="number"
                  value={lessonForm.xp_reward}
                  onChange={(e) => setLessonForm({ ...lessonForm, xp_reward: parseInt(e.target.value) || 0 })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  min={0}
                  required
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button type="button" onClick={() => setModalType(null)} className="flex-1 py-2 text-gray-600 hover:bg-gray-100 rounded-xl">
                  Cancel
                </button>
                <button type="submit" disabled={isSaving} className="flex-1 py-2 bg-purple hover:bg-purple/90 text-white font-medium rounded-xl disabled:opacity-50">
                  {isSaving ? <LoadingSpinner className="w-5 h-5 mx-auto" /> : "Save"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Question Modal */}
      {modalType === "question" && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 overflow-y-auto">
          <div className="bg-white rounded-2xl p-6 w-full max-w-lg my-8">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-gray-900">
                {editingItem ? "Edit Question" : "Add Question"}
              </h2>
              <button onClick={() => setModalType(null)} className="text-gray-400 hover:text-gray-600">
                <XIcon className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleSaveQuestion} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Type</label>
                <select
                  value={questionForm.type}
                  onChange={(e) => setQuestionForm({ ...questionForm, type: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                >
                  <option value="single_choice">Single Choice</option>
                  <option value="multi_choice">Multiple Choice</option>
                  <option value="listen_reply">Listen & Reply</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Prompt Text</label>
                <textarea
                  value={questionForm.prompt_text}
                  onChange={(e) => setQuestionForm({ ...questionForm, prompt_text: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  rows={2}
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Correct Answer</label>
                <input
                  type="text"
                  value={questionForm.correct_answer}
                  onChange={(e) => setQuestionForm({ ...questionForm, correct_answer: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                  required
                />
              </div>
              {(questionForm.type === "single_choice" || questionForm.type === "multi_choice") && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Options</label>
                  {questionForm.options.map((opt, idx) => (
                    <input
                      key={idx}
                      type="text"
                      value={opt}
                      onChange={(e) => {
                        const newOptions = [...questionForm.options];
                        newOptions[idx] = e.target.value;
                        setQuestionForm({ ...questionForm, options: newOptions });
                      }}
                      className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50 mb-2"
                      placeholder={`Option ${idx + 1}`}
                    />
                  ))}
                </div>
              )}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Hint (optional)</label>
                <input
                  type="text"
                  value={questionForm.hint}
                  onChange={(e) => setQuestionForm({ ...questionForm, hint: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple/50"
                />
              </div>
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="use_tts"
                  checked={questionForm.use_tts}
                  onChange={(e) => setQuestionForm({ ...questionForm, use_tts: e.target.checked })}
                  className="w-4 h-4 text-purple rounded"
                />
                <label htmlFor="use_tts" className="text-sm text-gray-700">
                  Use Text-to-Speech
                </label>
              </div>
              <div className="flex gap-3 pt-4">
                <button type="button" onClick={() => setModalType(null)} className="flex-1 py-2 text-gray-600 hover:bg-gray-100 rounded-xl">
                  Cancel
                </button>
                <button type="submit" disabled={isSaving} className="flex-1 py-2 bg-purple hover:bg-purple/90 text-white font-medium rounded-xl disabled:opacity-50">
                  {isSaving ? <LoadingSpinner className="w-5 h-5 mx-auto" /> : "Save"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
