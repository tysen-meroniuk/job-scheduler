import { useQuery } from '@tanstack/react-query'
import { api, type JobState } from '../api/client'

const STATES: JobState[] = ['queued', 'running', 'done', 'failed']

export default function Metrics() {
  const { data } = useQuery({
    queryKey: ['metrics'],
    queryFn: api.metrics,
  })

  return (
    <section className="grid grid-cols-4 gap-3">
      {STATES.map(state => (
        <div key={state} className="bg-white rounded border p-3">
          <div className="text-xs uppercase text-gray-500">{state}</div>
          <div className="text-2xl font-semibold">
            {data?.by_state?.[state] ?? 0}
          </div>
        </div>
      ))}
    </section>
  )
}
