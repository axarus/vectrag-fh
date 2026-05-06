import { useState, useEffect } from 'react';
import { apiService } from '../../../shared/services/api';
import type { Model } from '../../../types';

export function useModels() {
  const [models, setModels] = useState<Model[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadModels = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await apiService.getModels();
      setModels(data);
      console.log(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load models');
    } finally {
      setLoading(false);
    }
  };

  const deleteModel = async (id: string): Promise<boolean> => {
    try {
      await apiService.deleteModel(id);
      await loadModels();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete model');
      return false;
    }
  };

  useEffect(() => {
    loadModels();
  }, []);

  return {
    models,
    loading,
    error,
    loadModels,
    deleteModel,
  };
}

