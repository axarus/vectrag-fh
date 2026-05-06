import { Link } from 'react-router-dom';
import type { Model } from '../../../types';

interface ModelCardProps {
  model: Model;
  onDelete: (id: string) => void;
}

export default function ModelCard({ model, onDelete }: ModelCardProps) {
  const handleDelete = () => {
    if (confirm(`Are you sure you want to delete model "${model.name}"?`)) {
      onDelete(model.name);
    }
  };

  const fieldCount = model.fields?.length ?? 0;
  return (
    <div className="card bg-base-100 shadow">
      <div className="card-body">
        <div className="flex justify-between items-start">
          <div>
            <h2 className="card-title">{model.name}</h2>
            <p className="text-sm text-base-content/70">{model.slug}</p>
            {model.description && (
              <p className="text-sm mt-2">{model.description}</p>
            )}
          </div>
          <div
            className={`badge ${
              model.status === 'publish' ? 'badge-success' : 'badge-warning'
            }`}
          >
            {model.status}
          </div>
        </div>
        <div className="mt-4">
          <p className="text-sm text-base-content/70">
            {fieldCount} field{fieldCount !== 1 ? 's' : ''}
          </p>
        </div>
        <div className="card-actions justify-end mt-4">
          <Link
              to={`/models/${model.name}`}
            className="btn btn-sm btn-primary"
          >
            Edit
          </Link>
          <button
            onClick={handleDelete}
            className="btn btn-sm btn-error"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}

