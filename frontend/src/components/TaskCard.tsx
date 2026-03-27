'use client';

import { Task, TaskStatus, Priority } from '@/types';
import { Calendar, User, Edit, Trash2, CheckCircle2, Clock, Circle } from 'lucide-react';
import { format } from 'date-fns';

interface TaskCardProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (id: number) => void;
  onStatusChange: (id: number, status: TaskStatus) => void;
  canModify: boolean;
  showOwner?: boolean;
}

const statusConfig: Record<TaskStatus, { label: string; color: string; icon: React.ReactNode }> = {
  TODO: { label: 'To Do', color: 'bg-gray-100 text-gray-700', icon: <Circle className="h-3.5 w-3.5" /> },
  IN_PROGRESS: { label: 'In Progress', color: 'bg-yellow-100 text-yellow-700', icon: <Clock className="h-3.5 w-3.5" /> },
  DONE: { label: 'Done', color: 'bg-green-100 text-green-700', icon: <CheckCircle2 className="h-3.5 w-3.5" /> },
};

const priorityConfig: Record<Priority, { label: string; color: string; dot: string }> = {
  LOW: { label: 'Low', color: 'bg-blue-50 text-blue-600', dot: 'bg-blue-400' },
  MEDIUM: { label: 'Medium', color: 'bg-orange-50 text-orange-600', dot: 'bg-orange-400' },
  HIGH: { label: 'High', color: 'bg-red-50 text-red-600', dot: 'bg-red-500' },
};

export default function TaskCard({ task, onEdit, onDelete, onStatusChange, canModify, showOwner }: TaskCardProps) {
  const status = statusConfig[task.status];
  const priority = priorityConfig[task.priority];
  const isOverdue = task.dueDate && new Date(task.dueDate) < new Date() && task.status !== 'DONE';

  const nextStatus: Record<TaskStatus, TaskStatus> = {
    TODO: 'IN_PROGRESS',
    IN_PROGRESS: 'DONE',
    DONE: 'TODO',
  };

  return (
    <div className={`surface-card p-5 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md flex flex-col gap-3 ${
      task.status === 'DONE' ? 'opacity-75' : ''
    }`}>
      {/* Header */}
      <div className="flex items-start justify-between gap-2">
        <h3 className={`flex-1 text-sm font-semibold leading-snug text-slate-900 ${task.status === 'DONE' ? 'text-slate-500 line-through' : ''}`}>
          {task.title}
        </h3>
        {canModify && (
          <div className="flex items-center gap-1 shrink-0">
            <button
              onClick={() => onEdit(task)}
              className="rounded-lg p-1.5 text-slate-400 transition-colors hover:bg-blue-50 hover:text-blue-600"
              title="Edit task"
            >
              <Edit className="h-3.5 w-3.5" />
            </button>
            <button
              onClick={() => onDelete(task.id)}
              className="rounded-lg p-1.5 text-slate-400 transition-colors hover:bg-red-50 hover:text-red-600"
              title="Delete task"
            >
              <Trash2 className="h-3.5 w-3.5" />
            </button>
          </div>
        )}
      </div>

      {/* Description */}
      {task.description && (
        <p className="line-clamp-2 text-xs leading-relaxed text-slate-500">{task.description}</p>
      )}

      {/* Badges */}
      <div className="flex flex-wrap items-center gap-2">
        <span className={`status-badge ${status.color}`}>
          {status.icon}
          {status.label}
        </span>
        <span className={`status-badge gap-1.5 ${priority.color}`}>
          <span className={`h-1.5 w-1.5 rounded-full ${priority.dot}`} />
          {priority.label}
        </span>
      </div>

      {/* Footer */}
      <div className="flex items-center justify-between border-t border-slate-100 pt-1">
        <div className="flex items-center gap-3 text-xs text-slate-400">
          {task.dueDate && (
            <span className={`flex items-center gap-1 ${isOverdue ? 'font-medium text-red-500' : ''}`}>
              <Calendar className="h-3.5 w-3.5" />
              {format(new Date(task.dueDate), 'MMM d, yyyy')}
              {isOverdue && ' (Overdue)'}
            </span>
          )}
          {showOwner && (
            <span className="flex items-center gap-1">
              <User className="h-3.5 w-3.5" />
              {task.ownerUsername}
            </span>
          )}
        </div>

        {/* Quick status toggle */}
        {canModify && task.status !== 'DONE' && (
          <button
            onClick={() => onStatusChange(task.id, nextStatus[task.status])}
            className="text-xs font-medium text-blue-600 transition-colors hover:text-blue-800 hover:underline"
          >
            → {statusConfig[nextStatus[task.status]].label}
          </button>
        )}
      </div>
    </div>
  );
}

