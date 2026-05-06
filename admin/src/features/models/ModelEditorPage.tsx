import { useNavigate, useParams } from 'react-router-dom';
import { useModel } from './hooks/useModel';
import ModelForm from './components/ModelForm';
import LoadingSpinner from '../../shared/components/LoadingSpinner';

export default function ModelEditorPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const {
    model,
    loading,
    saving,
    error,
    validationErrors,
    isNew,
    fieldTypes,
    saveModel,
    updateModel,
    updateModelName,
    addField,
    removeField,
    updateField,
  } = useModel(id);

  const handleSave = async () => {
    const success = await saveModel();
    if (!success && error) {
      alert(error);
    }
  };

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">
          {isNew ? 'Create Model' : 'Edit Model'}
        </h1>
        <div className="flex gap-2">
          <button
            onClick={() => navigate('/models')}
            className="btn btn-ghost"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={saving || validationErrors}
            className="btn btn-primary"
          >
            {saving ? (
              <>
                <span className="loading loading-spinner loading-sm"></span>
                Saving...
              </>
            ) : (
              'Save'
            )}
          </button>
        </div>
      </div>

      <ModelForm
        model={model}
        fieldTypes={fieldTypes}
        onModelChange={updateModel}
        onNameChange={updateModelName}
        onAddField={addField}
        onRemoveField={removeField}
        onFieldChange={updateField}
      />
    </div>
  );
}
