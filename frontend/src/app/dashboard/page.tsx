"use client";

import { useEffect, useState } from 'react';
import Link from 'next/link';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

async function safeFetchJson(url: string, options?: RequestInit) {
  const res = await fetch(url, options);
  const text = await res.text();
  const data = text ? JSON.parse(text) : {};
  if (!res.ok) {
    throw new Error(data.error || `Server responded with status ${res.status}`);
  }
  return data;
}

interface Enrollment {
  id: number;
  course_id: number;
  enrolled_at: string;
  course: {
    id: number;
    title: string;
    description: string;
    instructor?: { full_name: string };
  };
}

export default function DashboardPage() {
  const [enrollments, setEnrollments] = useState<Enrollment[]>([]);
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const storedUser = localStorage.getItem('lms_user');
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    }
    fetchMyEnrollments();
  }, []);

  const fetchMyEnrollments = async () => {
    const token = localStorage.getItem('lms_token');
    if (!token) {
      setError('Please sign in to view your dashboard');
      setLoading(false);
      return;
    }

    try {
      const data = await safeFetchJson(`${API_BASE_URL}/my-enrollments`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setEnrollments(data.enrollments || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleSignOut = () => {
    localStorage.removeItem('lms_token');
    localStorage.removeItem('lms_user');
    window.location.href = '/auth';
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-50 p-8 font-sans">
      <header className="max-w-6xl mx-auto flex justify-between items-center pb-8 border-b border-slate-800">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Student Dashboard</h1>
          <p className="text-sm text-slate-400 mt-1">
            {user ? (
              <>Welcome back, <span className="text-indigo-400 font-semibold">{user.full_name}</span>!</>
            ) : (
              <>Access your courses and learning progress</>
            )}
          </p>
        </div>
        <div className="flex items-center space-x-4">
          <Link href="/courses" className="bg-indigo-600 hover:bg-indigo-500 text-white px-4 py-2 rounded-lg text-xs font-semibold transition">
            Explore Courses
          </Link>
          {user ? (
            <button
              onClick={handleSignOut}
              className="border border-slate-700 hover:border-slate-500 text-slate-300 px-4 py-2 rounded-lg text-xs font-semibold transition"
            >
              Sign Out
            </button>
          ) : (
            <Link
              href="/auth"
              className="bg-indigo-600 hover:bg-indigo-500 text-white px-4 py-2 rounded-lg text-xs font-semibold transition"
            >
              Sign In
            </Link>
          )}
        </div>
      </header>

      <main className="max-w-6xl mx-auto py-10 space-y-8">
        <div>
          <h2 className="text-xl font-bold mb-4">My Enrolled Courses</h2>

          {loading && <div className="text-center text-slate-400 py-12">Loading your enrolled courses...</div>}
          
          {error && (
            <div className="bg-rose-950/60 border border-rose-800 text-rose-300 p-6 rounded-xl text-center space-y-4 max-w-lg mx-auto">
              <p className="text-sm font-medium">{error}</p>
              <Link
                href="/auth"
                className="inline-block bg-indigo-600 hover:bg-indigo-500 text-white px-5 py-2 rounded-lg text-xs font-semibold transition"
              >
                Sign In Now
              </Link>
            </div>
          )}

          {!loading && enrollments.length === 0 && !error && (
            <div className="bg-slate-900 border border-slate-800 rounded-xl p-8 text-center space-y-4">
              <p className="text-slate-400 text-sm">You have not enrolled in any courses yet.</p>
              <Link href="/courses" className="inline-block bg-indigo-600 hover:bg-indigo-500 text-white px-5 py-2 rounded-lg text-xs font-semibold transition">
                Browse Course Catalog
              </Link>
            </div>
          )}

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {enrollments.map((item) => (
              <div key={item.id} className="bg-slate-900 border border-slate-800 rounded-xl p-6 flex flex-col justify-between space-y-4">
                <div className="space-y-2">
                  <div className="text-xs font-semibold uppercase tracking-wider text-indigo-400">
                    Enrolled: {new Date(item.enrolled_at).toLocaleDateString()}
                  </div>
                  <h3 className="text-xl font-bold text-slate-100">{item.course?.title}</h3>
                  <p className="text-sm text-slate-400 line-clamp-2">{item.course?.description}</p>
                </div>
                <div className="pt-4 border-t border-slate-800/60 flex items-center justify-between">
                  <span className="text-xs text-slate-500">Instructor: {item.course?.instructor?.full_name || 'LMS Faculty'}</span>
                  <Link
                    href={`/courses/${item.course_id}`}
                    className="bg-slate-800 hover:bg-slate-700 text-slate-200 px-3 py-1.5 rounded-lg text-xs font-medium transition"
                  >
                    Continue Learning
                  </Link>
                </div>
              </div>
            ))}
          </div>
        </div>
      </main>
    </div>
  );
}
