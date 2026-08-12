"use client";

import { useEffect, useState } from 'react';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

interface Lesson {
  id: number;
  title: string;
  sequence: number;
}

interface Course {
  id: number;
  title: string;
  description: string;
  instructor?: { full_name: string };
  lessons?: Lesson[];
}

export default function CoursesPage() {
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [enrolledCourseIds, setEnrolledCourseIds] = useState<number[]>([]);

  useEffect(() => {
    fetchCourses();
    fetchMyEnrollments();
  }, []);

  const fetchCourses = async () => {
    try {
      const res = await fetch(`${API_BASE_URL}/courses`);
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Failed to fetch courses');
      setCourses(data.courses || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchMyEnrollments = async () => {
    const token = localStorage.getItem('lms_token');
    if (!token) return;

    try {
      const res = await fetch(`${API_BASE_URL}/my-enrollments`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const data = await res.json();
      if (res.ok && data.enrollments) {
        const ids = data.enrollments.map((e: any) => e.course_id);
        setEnrolledCourseIds(ids);
      }
    } catch (err) {
      // User not logged in or token expired
    }
  };

  const handleEnroll = async (courseId: number) => {
    const token = localStorage.getItem('lms_token');
    if (!token) {
      alert('Please sign in first to enroll in courses');
      window.location.href = '/auth';
      return;
    }

    try {
      const res = await fetch(`${API_BASE_URL}/courses/${courseId}/enroll`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });

      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Enrollment failed');

      alert('Enrolled successfully!');
      setEnrolledCourseIds([...enrolledCourseIds, courseId]);
    } catch (err: any) {
      alert(err.message);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-50 p-8 font-sans">
      <header className="max-w-6xl mx-auto flex justify-between items-center pb-8 border-b border-slate-800">
        <h1 className="text-3xl font-bold tracking-tight">Available Courses</h1>
        <a href="/" className="text-sm font-medium text-slate-400 hover:text-indigo-400 transition">
          ← Back to Home
        </a>
      </header>

      <main className="max-w-6xl mx-auto py-10">
        {loading && <div className="text-center text-slate-400 py-12">Loading courses...</div>}
        {error && <div className="bg-rose-950/60 border border-rose-800 text-rose-300 text-sm p-4 rounded-lg text-center mb-8">{error}</div>}

        {!loading && courses.length === 0 && (
          <div className="text-center text-slate-500 py-16">No courses available at the moment.</div>
        )}

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {courses.map((course) => {
            const isEnrolled = enrolledCourseIds.includes(course.id);
            return (
              <div key={course.id} className="bg-slate-900 border border-slate-800 rounded-xl p-6 flex flex-col justify-between space-y-4 hover:border-slate-700 transition">
                <div className="space-y-2">
                  <div className="text-xs font-semibold uppercase tracking-wider text-indigo-400">
                    Instructor: {course.instructor?.full_name || 'LMS Faculty'}
                  </div>
                  <h3 className="text-xl font-bold text-slate-100">{course.title}</h3>
                  <p className="text-sm text-slate-400 line-clamp-3">{course.description || 'No description provided.'}</p>
                </div>

                <div className="pt-4 border-t border-slate-800/60 flex items-center justify-between">
                  <span className="text-xs text-slate-500 font-medium">
                    {course.lessons?.length || 0} Lessons
                  </span>
                  {isEnrolled ? (
                    <span className="bg-emerald-950/80 border border-emerald-800 text-emerald-300 px-3 py-1.5 rounded-lg text-xs font-semibold">
                      Enrolled
                    </span>
                  ) : (
                    <button
                      onClick={() => handleEnroll(course.id)}
                      className="bg-indigo-600 hover:bg-indigo-500 text-white px-4 py-1.5 rounded-lg text-xs font-semibold transition"
                    >
                      Enroll Now
                    </button>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </main>
    </div>
  );
}
