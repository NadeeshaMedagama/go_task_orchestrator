'use client';

import { useState, useEffect, useSyncExternalStore } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { useTasks } from '@/hooks/useTasks';
import { Task, TaskFilters, TaskRequest, TaskStatus } from '@/types';
import Navbar from '@/components/Navbar';
import TaskCard from '@/components/TaskCard';
import TaskForm from '@/components/TaskForm';
import FilterBar from '@/components/FilterBar';
import Pagination from '@/components/Pagination';
import LoadingSpinner from '@/components/LoadingSpinner';
import apiClient from '@/lib/axios';
import toast from 'react-hot-toast';
import { Plus, ListTodo, AlertCircle } from 'lucide-react';

export default function TasksPageClient() {
  const { user, isLoading: authLoading, isAdmin } = useAuth();
  const router = useRouter();
  const isMounted = useSyncExternalStore(
    () => () => {},
    () => true,
    () => false
  );

  const [filters, setFilters] = useState<TaskFilters>({
    status: undefined,
    priority: undefined,
    page: 0,
    size: 12,
    sortBy: 'createdAt',
    direction: 'desc',
  });

  const [showForm, setShowForm] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  const { data, isLoading, error, refetch } = useTasks(filters);

  // Redirect if not logged in
  useEffect(() => {
    if (!authLoading && !user) {
      router.replace('/login');
    }
  }, [user, authLoading, router]);

  if (!isMounted) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (!user) return null;

  const handleFilterChange = (newFilters: Partial<TaskFilters>) => {
    setFilters((prev) => ({ ...prev, ...newFilters }));
  };

  const handleCreate = async (taskData: TaskRequest) => {
    await apiClient.post('/tasks', taskData);
    toast.success('Task created!');
    setShowForm(false);
    refetch();
  };

  const handleUpdate = async (taskData: TaskRequest) => {
    if (!editingTask) return;
    await apiClient.put(`/tasks/${editingTask.id}`, taskData);
    toast.success('Task updated!');
    setEditingTask(null);
    refetch();
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this task?')) return;
    try {
      await apiClient.delete(`/tasks/${id}`);
      toast.success('Task deleted.');
      refetch();
    } catch {
      toast.error('Failed to delete task.');
    }
  };

  const handleStatusChange = async (id: number, status: TaskStatus) => {
    try {
      await apiClient.patch(`/tasks/${id}/status`, { status });
      toast.success(`Moved to ${status.replace('_', ' ')}`);
      refetch();
    } catch {
      toast.error('Failed to update status.');
    }
  };

  const canModifyTask = (task: Task) =>
    isAdmin || task.ownerId === user.userId;

  const tasks = data?.content ?? [];
  const todoCount = tasks.filter((t) => t.status === 'TODO').length;
  const inProgressCount = tasks.filter((t) => t.status === 'IN_PROGRESS').length;
  const doneCount = tasks.filter((t) => t.status === 'DONE').length;

  return (
    <div className="page-shell">
      <Navbar />

      <main className="page-container py-8">
        {/* Page header */}
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="flex items-center gap-2 text-2xl font-bold text-slate-900">
              <ListTodo className="h-6 w-6 text-blue-600" />
              {isAdmin ? 'All Tasks' : 'My Tasks'}
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              {isAdmin
                ? 'Viewing all tasks across all users'
                : 'Manage and track your personal tasks'}
            </p>
          </div>
          <button
            onClick={() => { setEditingTask(null); setShowForm(true); }}
            className="btn-primary shadow-sm"
          >
            <Plus className="h-4 w-4" />
            New Task
          </button>
        </div>

        {/* Stats bar */}
        {data && (
          <div className="grid grid-cols-3 gap-4 mb-6">
            {[
              { label: 'To Do', count: todoCount, color: 'bg-gray-100 text-gray-700' },
              { label: 'In Progress', count: inProgressCount, color: 'bg-yellow-100 text-yellow-700' },
              { label: 'Done', count: doneCount, color: 'bg-green-100 text-green-700' },
            ].map(({ label, count, color }) => (
              <div key={label} className={`surface-card flex items-center justify-between px-4 py-3 ${color}`}>
                <span className="text-sm font-medium">{label}</span>
                <span className="text-lg font-bold">{count}</span>
              </div>
            ))}
          </div>
        )}

        {/* Filters */}
        <div className="mb-6">
          <FilterBar filters={filters} onChange={handleFilterChange} />
        </div>

        {/* Task Grid */}
        {isLoading ? (
          <div className="flex items-center justify-center py-24">
            <LoadingSpinner size="lg" />
          </div>
        ) : error ? (
          <div className="flex flex-col items-center justify-center py-24 text-center">
            <AlertCircle className="h-12 w-12 text-red-400 mb-3" />
            <p className="font-medium text-slate-700">{error}</p>
            <button onClick={refetch} className="mt-3 text-sm text-blue-600 transition hover:text-blue-800 hover:underline">
              Try again
            </button>
          </div>
        ) : tasks.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-24 text-center">
            <div className="w-16 h-16 bg-blue-50 rounded-2xl flex items-center justify-center mb-4">
              <ListTodo className="h-8 w-8 text-blue-400" />
            </div>
            <h3 className="text-lg font-semibold text-slate-700">No tasks found</h3>
            <p className="mt-1 text-sm text-slate-400">
              {filters.status || filters.priority
                ? 'Try adjusting your filters'
                : 'Create your first task to get started'}
            </p>
            {!filters.status && !filters.priority && (
              <button
                onClick={() => { setEditingTask(null); setShowForm(true); }}
                className="btn-primary mt-4"
              >
                <Plus className="h-4 w-4" /> Create Task
              </button>
            )}
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {tasks.map((task) => (
              <TaskCard
                key={task.id}
                task={task}
                onEdit={(t) => { setEditingTask(t); setShowForm(true); }}
                onDelete={handleDelete}
                onStatusChange={handleStatusChange}
                canModify={canModifyTask(task)}
                showOwner={isAdmin}
              />
            ))}
          </div>
        )}

        {/* Pagination */}
        {data && (
          <Pagination
            currentPage={data.number}
            totalPages={data.totalPages}
            totalElements={data.totalElements}
            pageSize={data.size}
            onPageChange={(page) => handleFilterChange({ page })}
          />
        )}
      </main>

      {/* Task Form Modal */}
      {showForm && (
        <TaskForm
          task={editingTask}
          onSubmit={editingTask ? handleUpdate : handleCreate}
          onClose={() => { setShowForm(false); setEditingTask(null); }}
        />
      )}
    </div>
  );
}

