const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// MVP: 인증 미구현, 하드코딩 user_id
export const USER_ID = "019234ab-0000-7def-8000-000000000001";

export async function api<T>(path: string, options?: RequestInit): Promise<T> {
  const url = path.includes("?")
    ? `${API_URL}${path}&user_id=${USER_ID}`
    : `${API_URL}${path}?user_id=${USER_ID}`;

  const res = await fetch(url, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  if (!res.ok) {
    const error = await res.json().catch(() => ({}));
    throw new Error(error?.error?.message || `API Error: ${res.status}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}
