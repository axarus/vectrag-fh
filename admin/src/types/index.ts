export type FieldType = 'string' | 'text' | 'number' | 'boolean' | 'date' | 'datetime' | 'relation';

export type Status = 'draft' | 'publish' | 'delete';

export interface Field {
  name: string;
  type: FieldType;
  description?: string;
  unique: boolean;
  required: boolean;
  status: Status;
  createdAt?: string;
  updatedAt?: string;
}

export interface Model {
  name: string;
  slug: string;
  description?: string;
  fields: Field[];
  status: Status;
  schemaVersion: number;
}

