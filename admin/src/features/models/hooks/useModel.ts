import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiService } from '../../../shared/services/api';
import { modelSchema } from '../../../schemas/model.schema';
import type { Model, Field, FieldType } from '../../../types';

const FIELD_TYPES: FieldType[] = [
  'string',
  'text',
  'number',
  'boolean',
  'date',
  'datetime',
  'relation',
];

const createEmptyModel = (): Model => ({
  name: '',
  slug: '',
  description: '',
  fields: [],
  status: 'draft',
  schemaVersion: 1,
});

const createEmptyField = (): Field => ({
  name: '',
  type: 'string',
  description: '',
  unique: false,
  required: false,
  status: 'draft',
});

export function useModel(id: string | undefined) {
  const navigate = useNavigate();
  const isNew = id === 'new' || !id;

  const [model, setModel] = useState<Model>(createEmptyModel());
  const [loading, setLoading] = useState(!isNew);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [validationErrors, setValidationErrors] = useState<any>(null);

  useEffect(() => {
    if (!isNew && id) {
      loadModel(id);
    }
  }, [id, isNew]);

  useEffect(() => {
    if (!isNew) {
      validateModel(model);
    }
  }, [model]);

  const loadModel = async (modelId: string) => {
    try {
      setLoading(true);
      setError(null);
      const data = await apiService.getModel(modelId);
      setModel(data);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to load model';
      setError(message);
      navigate('/models');
    } finally {
      setLoading(false);
    }
  };

  const validateModel = (modelToValidate: Model) => {
    try {
      const result = modelSchema.safeParse(modelToValidate);
      if (result.success) {
        setValidationErrors(null);
        return true;
      } else {
        setValidationErrors(result.error.format());
        return false;
      }
    } catch (err) {
      setValidationErrors({ _errors: ['Validation failed'] });
      return false;
    }
  };

  const saveModel = async (): Promise<boolean> => {
    if (!validateModel(model)) {
      return false;
    }

    try {
      setSaving(true);
      setError(null);
      if (isNew) {
        console.log(model);
        await apiService.createModel(model);
      } else {
        await apiService.updateModel(id!, model);
      }
      navigate('/models');
      return true;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to save model';
      setError(message);
      return false;
    } finally {
      setSaving(false);
    }
  };

  const updateModel = (updates: Partial<Model>) => {
    setModel((prev) => ({ ...prev, ...updates }));
  };

  const updateModelName = (name: string) => {
    const slug = name
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
    setModel((prev) => ({ ...prev, name, slug }));
  };

  const addField = () => {
    setModel((prev) => ({
      ...prev,
      fields: [...prev.fields, createEmptyField()],
    }));
  };

  const removeField = (index: number) => {
    setModel((prev) => ({
      ...prev,
      fields: prev.fields.filter((_, i) => i !== index),
    }));
  };

  const updateField = (index: number, updates: Partial<Field>) => {
    setModel((prev) => {
      const newFields = [...prev.fields];
      newFields[index] = { ...newFields[index], ...updates };
      return { ...prev, fields: newFields };
    });
  };

  return {
    model,
    loading,
    saving,
    error,
    validationErrors,
    isNew,
    fieldTypes: FIELD_TYPES,
    saveModel,
    updateModel,
    updateModelName,
    addField,
    removeField,
    updateField,
  };
}

