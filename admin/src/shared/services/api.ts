import type { Model } from '../../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:51987/api';

class ApiService {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return response.json();
  }

  async getModels(): Promise<Model[]> {
    return this.request<Model[]>('/models');
  }

  async getModel(id: string): Promise<Model> {
    return this.request<Model>(`/models/${id}`);
  }

  async createModel(model: Model): Promise<Model> {
    return this.request<Model>('/models', {
      method: 'POST',
      body: JSON.stringify(model),
    });
  }

  async updateModel(id: string, model: Model): Promise<Model> {
    return this.request<Model>(`/models/${id}`, {
      method: 'PUT',
      body: JSON.stringify(model),
    });
  }

  async deleteModel(id: string): Promise<void> {
    return this.request<void>(`/models/${id}`, {
      method: 'DELETE',
    });
  }
}

export const apiService = new ApiService();

