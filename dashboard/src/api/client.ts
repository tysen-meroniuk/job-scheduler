export type JobState = 'queued' | 'running' | 'done' | 'failed'

export interface Job {
  id: number
  queue: string
  type: string
  payload: unknown
  state: JobState
  priority: number
  run_at: string
  attempts: number
  max_attempts: number
  last_error?: string
  created_at: string
  started_at?: string
  completed_at?: string
}

export interface EnqueueBody {
  queue?: string
  type: string
  payload: unknown
  priority?: number
  run_at?: string
  max_attempts?: number
  idempotency_key?: string
}

export interface Metrics {
  by_state: Record<JobState, number>
  by_queue: Record<string, number>
}

const API = '/api'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...init?.headers },
  })
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
  return res.json() as Promise<T>
}

export const api = {
  listJobs: () => request<Job[]>('/jobs'),
  getJob: (id: number) => request<Job>(`/jobs/${id}`),
  enqueue: (body: EnqueueBody) =>
    request<{ id: number }>('/jobs', {
      method: 'POST',
      body: JSON.stringify(body),
    }),
  metrics: () => request<Metrics>('/metrics'),
}
