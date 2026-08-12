import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-50 flex flex-col justify-between p-8 font-sans">
      {/* Header / Navbar */}
      <header className="flex justify-between items-center max-w-6xl mx-auto w-full py-4 border-b border-slate-800">
        <Link href="/" className="flex items-center space-x-2">
          <div className="w-8 h-8 rounded-lg bg-indigo-600 flex items-center justify-center font-bold text-lg text-white">
            L
          </div>
          <span className="text-xl font-bold tracking-tight text-white">Mini LMS Enterprise</span>
        </Link>
        <nav className="flex items-center space-x-6 text-sm font-medium text-slate-300">
          <Link href="/courses" className="hover:text-white transition">Courses</Link>
          <Link href="/dashboard" className="hover:text-white transition">Dashboard</Link>
          <Link href="/auth" className="bg-indigo-600 hover:bg-indigo-500 text-white px-4 py-2 rounded-lg font-semibold text-xs tracking-wider uppercase transition">
            Sign In
          </Link>
        </nav>
      </header>

      {/* Hero Section */}
      <main className="max-w-4xl mx-auto text-center my-16 space-y-6">
        <div className="inline-block bg-indigo-950/60 border border-indigo-800/50 text-indigo-300 px-3 py-1 rounded-full text-xs font-semibold uppercase tracking-widest">
          Go + Next.js Enterprise Architecture
        </div>
        <h1 className="text-5xl md:text-6xl font-extrabold tracking-tight text-white leading-tight">
          Scalable Learning Management System
        </h1>
        <p className="text-lg text-slate-400 max-w-2xl mx-auto">
          High-performance RESTful API powered by Go (Gin + GORM + PostgreSQL) coupled with a modern dynamic Next.js App Router frontend.
        </p>
        <div className="flex justify-center space-x-4 pt-4">
          <Link href="/courses" className="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-3 rounded-lg font-semibold transition shadow-lg shadow-indigo-600/30">
            Browse Courses
          </Link>
          <Link href="/auth" className="border border-slate-700 hover:border-slate-500 text-slate-300 px-6 py-3 rounded-lg font-semibold transition">
            Sign In / Register
          </Link>
        </div>
      </main>

      {/* Features Grid */}
      <section className="max-w-5xl mx-auto grid md:grid-cols-3 gap-6 my-12">
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 space-y-3">
          <div className="text-indigo-400 font-bold text-lg">Gin Web Framework</div>
          <p className="text-slate-400 text-sm">
            High performance routing with middleware support for CORS, RBAC, and error handling.
          </p>
        </div>
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 space-y-3">
          <div className="text-indigo-400 font-bold text-lg">JWT & Role Auth</div>
          <p className="text-slate-400 text-sm">
            Secure token-based authentication supporting Student, Instructor, and Admin roles.
          </p>
        </div>
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 space-y-3">
          <div className="text-indigo-400 font-bold text-lg">PostgreSQL & GORM</div>
          <p className="text-slate-400 text-sm">
            Clean ORM models with auto-migrations and structured relational domain mapping.
          </p>
        </div>
      </section>

      {/* Footer */}
      <footer className="max-w-6xl mx-auto w-full text-center border-t border-slate-800 pt-6 text-slate-500 text-xs">
        © 2026 Mini LMS Enterprise Project. All rights reserved.
      </footer>
    </div>
  );
}
