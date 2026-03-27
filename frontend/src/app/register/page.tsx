'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import apiClient from '@/lib/axios';
import { AuthUser } from '@/types';
import { CheckSquare, Eye, EyeOff, User, Mail, Lock } from 'lucide-react';
import { AxiosError } from 'axios';
import toast from 'react-hot-toast';
import LoadingSpinner from '@/components/LoadingSpinner';

interface FormErrors {
  username?: string;
  email?: string;
  password?: string;
}

export default function RegisterPage() {
  const [formData, setFormData] = useState({ username: '', email: '', password: '' });
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<FormErrors>({});
  const [apiError, setApiError] = useState('');
  const { login } = useAuth();
  const router = useRouter();

  const validate = (): boolean => {
    const newErrors: FormErrors = {};
    if (!formData.username.trim()) newErrors.username = 'Username is required';
    else if (formData.username.length < 3) newErrors.username = 'Username must be at least 3 characters';
    if (!formData.email.trim()) newErrors.email = 'Email is required';
    else if (!/\S+@\S+\.\S+/.test(formData.email)) newErrors.email = 'Enter a valid email';
    if (!formData.password) newErrors.password = 'Password is required';
    else if (formData.password.length < 6) newErrors.password = 'Password must be at least 6 characters';
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setApiError('');
    if (!validate()) return;
    setIsLoading(true);
    try {
      const { data } = await apiClient.post<AuthUser>('/auth/register', formData);
      login(data);
      toast.success('Account created successfully!');
      router.push('/tasks');
    } catch (err) {
      const axiosError = err as AxiosError<{ message: string; validationErrors?: string[] }>;
      const msgs = axiosError.response?.data?.validationErrors;
      setApiError(msgs ? msgs.join(', ') : axiosError.response?.data?.message || 'Registration failed.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-shell">
      <div className="auth-panel">
        {/* Logo */}
        <div className="auth-brand">
          <div className="auth-brand-icon">
            <CheckSquare className="h-8 w-8 text-white" />
          </div>
          <h1 className="text-2xl font-bold text-slate-900">Join TaskFlow</h1>
          <p className="mt-1 text-sm text-slate-500">Create your account to get started</p>
        </div>

        {/* Card */}
        <div className="surface-card-elevated p-8">
          <h2 className="mb-6 text-xl font-semibold text-slate-900">Create Account</h2>

          {apiError && (
            <div className="error-banner mb-4">
              {apiError}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Username */}
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                <span className="flex items-center gap-1.5"><User className="h-3.5 w-3.5" />Username</span>
              </label>
              <input
                type="text"
                value={formData.username}
                onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                placeholder="Choose a username"
                className={`control-input ${errors.username ? 'border-red-400 focus-visible:ring-red-500/30' : ''}`}
              />
              {errors.username && <p className="text-xs text-red-500 mt-1">{errors.username}</p>}
            </div>

            {/* Email */}
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                <span className="flex items-center gap-1.5"><Mail className="h-3.5 w-3.5" />Email</span>
              </label>
              <input
                type="email"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                placeholder="Enter your email"
                className={`control-input ${errors.email ? 'border-red-400 focus-visible:ring-red-500/30' : ''}`}
              />
              {errors.email && <p className="text-xs text-red-500 mt-1">{errors.email}</p>}
            </div>

            {/* Password */}
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                <span className="flex items-center gap-1.5"><Lock className="h-3.5 w-3.5" />Password</span>
              </label>
              <div className="relative">
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={formData.password}
                  onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                  placeholder="Min. 6 characters"
                  className={`control-input pr-10 ${errors.password ? 'border-red-400 focus-visible:ring-red-500/30' : ''}`}
                />
                <button type="button" onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 transition hover:text-slate-600">
                  {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
              {errors.password && <p className="text-xs text-red-500 mt-1">{errors.password}</p>}
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="btn-primary mt-2 w-full"
            >
              {isLoading ? <><LoadingSpinner size="sm" /> Creating account...</> : 'Create Account'}
            </button>
          </form>
        </div>

        <p className="mt-6 text-center text-sm text-slate-500">
          Already have an account?{' '}
          <Link href="/login" className="auth-link">Sign in</Link>
        </p>
      </div>
    </div>
  );
}

