import { z } from "zod";

export const fieldSchema = z.object({
    name: z.string().min(1, "Field name is required"),
    type: z.enum(['string', 'text', 'number', 'boolean', 'date', 'datetime', 'relation']),
    status: z.enum(["draft", "publish"]),
    description: z.string().optional(),
    required: z.boolean(),
    unique: z.boolean(),
});
