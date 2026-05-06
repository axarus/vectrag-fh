import { Link } from 'react-router-dom';
import { useDashboard } from './hooks/useDashboard';
import LoadingSpinner from '../../shared/components/LoadingSpinner';

export default function DashboardPage() {
  const { stats, loading } = useDashboard();

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-base-content/70 mt-2">
          Welcome to VectraG Development Panel
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="stat bg-base-100 shadow rounded-box">
          <div className="stat-figure text-primary">
            <svg
              className="w-8 h-8"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"
              />
            </svg>
          </div>
          <div className="stat-title">Total Models</div>
          <div className="stat-value text-primary">{stats.total}</div>
          <div className="stat-desc">Models created</div>
        </div>

        <div className="stat bg-base-100 shadow rounded-box">
          <div className="stat-figure text-secondary">
            <svg
              className="w-8 h-8"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
          </div>
          <div className="stat-title">Published</div>
          <div className="stat-value text-secondary">{stats.published}</div>
          <div className="stat-desc">Published models</div>
        </div>

        <div className="stat bg-base-100 shadow rounded-box">
          <div className="stat-figure text-accent">
            <svg
              className="w-8 h-8"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
              />
            </svg>
          </div>
          <div className="stat-title">Draft</div>
          <div className="stat-value text-accent-content">{stats.draft}</div>
          <div className="stat-desc">Draft models</div>
        </div>
      </div>

      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <h2 className="card-title">Quick Actions</h2>
          <div className="card-actions justify-end">
            <Link to="/models/new" className="btn btn-primary">
              Create New Model
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
