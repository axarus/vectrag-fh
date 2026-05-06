import {z} from 'zod';
import {fieldSchema} from "./field.schema.ts";

export const modelSchema = z.object({
    name: z.string().min(1, "Model name is required"),
    slug: z.string().min(1, "Slug is required"),
    description: z.string().optional(),
    status: z.enum(["draft", "publish"]),
    fields: z.array(fieldSchema).min(0),
    schemaVersion: z.number().min(1),
});
