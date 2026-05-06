import type { Field, FieldType, Status } from '../../../types';

interface FieldEditorProps {
  field: Field;
  index: number;
  fieldTypes: FieldType[];
  onChange: (updates: Partial<Field>) => void;
  onRemove: () => void;
}

export default function FieldEditor({
  field,
  index,
  fieldTypes,
  onChange,
  onRemove,
}: FieldEditorProps) {
  return (
    <div className="card bg-base-200">
      <div className="card-body">
        <div className="flex justify-between items-start mb-4">
          <h3 className="font-semibold">Field {index + 1}</h3>
          <button onClick={onRemove} className="btn btn-sm btn-error">
            Remove
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="form-control">
            <label className="label">
              <span className="label-text">Name</span>
            </label>
            <input
              type="text"
              value={field.name}
              onChange={(e) => onChange({ name: e.target.value })}
              placeholder="Email"
              className="input input-bordered"
            />
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Type</span>
            </label>
            <select
              value={field.type}
              onChange={(e) => onChange({ type: e.target.value as FieldType })}
              className="select select-bordered"
            >
              {fieldTypes.map((type) => (
                <option key={type} value={type}>
                  {type}
                </option>
              ))}
            </select>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Status</span>
            </label>
            <select
              value={field.status}
              onChange={(e) => onChange({ status: e.target.value as Status })}
              className="select select-bordered"
            >
              <option value="draft">Draft</option>
              <option value="publish">Publish</option>
            </select>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Description</span>
            </label>
            <input
              type="text"
              value={field.description || ''}
              onChange={(e) => onChange({ description: e.target.value })}
              placeholder="Field description"
              className="input input-bordered"
            />
          </div>
        </div>

        <div className="flex gap-4 mt-4">
          <label className="label cursor-pointer">
            <span className="label-text mr-2">Required</span>
            <input
              type="checkbox"
              checked={field.required}
              onChange={(e) => onChange({ required: e.target.checked })}
              className="checkbox"
            />
          </label>

          <label className="label cursor-pointer">
            <span className="label-text mr-2">Unique</span>
            <input
              type="checkbox"
              checked={field.unique}
              onChange={(e) => onChange({ unique: e.target.checked })}
              className="checkbox"
            />
          </label>
        </div>
      </div>
    </div>
  );
}

