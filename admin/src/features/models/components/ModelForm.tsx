import { useForm, useFieldArray } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import type { Model, Field, FieldType } from '../../../types';
import { modelSchema } from '../../../schemas';
import FieldEditor from './FieldEditor';
import { slugify } from '../../../utils/slugify';

type ModelFormData = z.infer<typeof modelSchema>;

interface ModelFormProps {
  model: Model;
  fieldTypes: string[];
  onModelChange: (updates: Partial<Model>) => void;
  onNameChange: (name: string) => void;
  onAddField: () => void;
  onRemoveField: (index: number) => void;
  onFieldChange: (index: number, updates: Partial<Field>) => void;
}

export default function ModelForm({
  model,
  fieldTypes,
  onModelChange,
  onNameChange,
  onAddField,
  onRemoveField,
  onFieldChange,
}: ModelFormProps) {
  const {
    register,
    handleSubmit,
    control,
    watch,
    setValue,
    formState: { errors },
  } = useForm<ModelFormData>({
    resolver: zodResolver(modelSchema),
    defaultValues: {
      name: model.name,
      slug: model.slug,
      description: model.description || '',
      status: model.status === 'delete' ? 'draft' : model.status,
      fields: model.fields.map(field => ({
        ...field,
        status: field.status === 'delete' ? 'draft' : field.status,
      })),
      schemaVersion: model.schemaVersion || 1,
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'fields',
  });

  const watchedFields = watch('fields');

  const handleNameChange = (value: string) => {
    setValue('name', value);
    setValue('slug', slugify(value));
    onNameChange(value);
    onModelChange({ slug: slugify(value) });
  };

  const handleAddField = () => {
    const newField = {
      name: '',
      type: 'string' as FieldType,
      description: '',
      unique: false,
      required: false,
      status: 'draft' as const,
    };
    append(newField);
    onAddField();
  };

  const handleRemoveField = (index: number) => {
    remove(index);
    onRemoveField(index);
  };

  const handleFieldChange = (index: number, updates: Partial<Field>) => {
    const statusValue = updates.status === 'delete' ? 'draft' : updates.status;
    setValue(`fields.${index}.name`, updates.name || watchedFields[index]?.name);
    setValue(`fields.${index}.type`, updates.type || watchedFields[index]?.type);
    setValue(`fields.${index}.description`, updates.description || watchedFields[index]?.description);
    setValue(`fields.${index}.required`, updates.required ?? watchedFields[index]?.required);
    setValue(`fields.${index}.unique`, updates.unique ?? watchedFields[index]?.unique);
    if (statusValue) {
      setValue(`fields.${index}.status`, statusValue);
    }
    onFieldChange(index, updates);
  };
  return (
    <form onSubmit={handleSubmit((data) => console.log('Form submitted:', data))}>
      <div className="card bg-base-100 shadow">
        <div className="card-body space-y-4">
          <div className="form-control">
            <label className="label">
              <span className="label-text">Name</span>
            </label>
            <input
              type="text"
              {...register('name')}
              onChange={(e) => handleNameChange(e.target.value)}
              placeholder="User"
              className={`input input-bordered ${errors.name ? 'input-error' : ''}`}
              required
            />
            {errors.name && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {errors.name.message}
                </span>
              </label>
            )}
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Slug</span>
            </label>
            <input
              type="text"
              {...register('slug')}
              placeholder="user"
              className="input input-bordered"
              disabled={true}
              required
            />
            {errors.slug && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {errors.slug.message}
                </span>
              </label>
            )}
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Description</span>
            </label>
            <textarea
              {...register('description')}
              placeholder="Model description"
              className={`textarea textarea-bordered ${errors.description ? 'input-error' : ''}`}
              rows={3}
            />
            {errors.description && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {errors.description.message}
                </span>
              </label>
            )}
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Status</span>
            </label>
            <select
              {...register('status')}
              className={`select select-bordered ${errors.status ? 'input-error' : ''}`}
            >
              <option value="draft">Draft</option>
              <option value="publish">Publish</option>
            </select>
            {errors.status && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {errors.status.message}
                </span>
              </label>
            )}
          </div>
        </div>
      </div>

      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h2 className="card-title">Fields</h2>
            <button type="button" onClick={handleAddField} className="btn btn-sm btn-primary">
              Add Field
            </button>
          </div>

          {fields.length === 0 ? (
            <p className="text-base-content/70 text-center py-8">
              No fields yet. Add your first field.
            </p>
          ) : (
            <div className="space-y-4">
              {fields.map((field, index) => (
                <FieldEditor
                  key={field.id}
                  field={watchedFields[index]}
                  index={index}
                  fieldTypes={fieldTypes as FieldType[]}
                  onChange={(updates) => handleFieldChange(index, updates)}
                  onRemove={() => handleRemoveField(index)}
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </form>
  );
}

