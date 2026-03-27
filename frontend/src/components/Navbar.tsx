'use client';

import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { LogOut, CheckSquare, User } from 'lucide-react';

export default function Navbar() {
  const { user, logout, isAdmin } = useAuth();

  return (
    <nav className="sticky top-0 z-50 border-b border-slate-200/80 bg-white/90 backdrop-blur-sm">
      <div className="page-container">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link href="/tasks" className="flex items-center gap-2 text-xl font-bold text-blue-600 transition-colors hover:text-blue-700">
            <CheckSquare className="h-6 w-6" />
            <span>TaskFlow</span>
          </Link>

          {/* Nav items */}
          {user && (
            <div className="flex items-center gap-4">
              <div className="surface-card flex items-center gap-2 px-3 py-1.5 text-sm text-slate-600">
                <User className="h-4 w-4" />
                <span className="font-medium">{user.username}</span>
                <span className={`status-badge ${
                  isAdmin ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'
                }`}>
                  {user.role}
                </span>
              </div>
              <button
                onClick={logout}
                className="btn-ghost hover:bg-red-50 hover:text-red-600"
              >
                <LogOut className="h-4 w-4" />
                <span>Logout</span>
              </button>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
}

