import { useState, useEffect } from 'react';
import { apiService } from '../../../shared/services/api';
// import type { Model } from '../../../types';

interface DashboardStats {
  total: number;
  published: number;
  draft: number;
}

export function useDashboard() {
  const [stats, setStats] = useState<DashboardStats>({
    total: 0,
    published: 0,
    draft: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      setLoading(true);
      const models = await apiService.getModels();
      const published = models.filter((m) => m.status === 'publish').length;
      const draft = models.filter((m) => m.status === 'draft').length;

      setStats({
        total: models.length,
        published,
        draft,
      });
    } catch (err) {
      console.error('Failed to load dashboard stats:', err);
    } finally {
      setLoading(false);
    }
  };

  return {
    stats,
    loading,
    refresh: loadStats,
  };
}

