export function slugify(name: string): string {
  let s = name.trim();
  if (!s) {
    throw new Error("name cannot be empty");
  }

  s = s.toLowerCase();

  s = s.replace(/[^a-z0-9]+/g, "-");

  s = s.replace(/^-+|-+$/g, "");

  if (!s) {
    throw new Error("name results in empty slug");
  }

  if (!/^[a-z0-9-]+$/.test(s)) {
    throw new Error("invalid slug");
  }

  return s;
}