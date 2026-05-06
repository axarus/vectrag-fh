import { Link } from 'react-router-dom';
import { useModels } from './hooks/useModels';
import ModelCard from './components/ModelCard';
import LoadingSpinner from '../../shared/components/LoadingSpinner';
import ErrorAlert from '../../shared/components/ErrorAlert';

export default function ModelsListPage() {
  const { models, loading, error, loadModels, deleteModel } = useModels();

  const handleDelete = async (id: string) => {
    const success = await deleteModel(id);
    if (!success) {
      alert('Failed to delete model');
    }
  };

  if (loading) {
    return <LoadingSpinner />;
  }

  if (error) {
    return <ErrorAlert message={error} onRetry={loadModels} />;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Models</h1>
          <p className="text-base-content/70 mt-2">
            Manage your content models
          </p>
        </div>
        <Link to="/models/new" className="btn btn-primary">
          <svg
            className="w-5 h-5 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 4v16m8-8H4"
            />
          </svg>
          Create Model
        </Link>
      </div>

      {models.length === 0 ? (
        <div className="card bg-base-100 shadow">
          <div className="card-body text-center">
            <p className="text-base-content/70">No models found</p>
            <Link to="/models/new" className="btn btn-primary mt-4">
              Create your first model
            </Link>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {models.map((model) => (
            <ModelCard
              key={model.name}
              model={model}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}
    </div>
  );
}
