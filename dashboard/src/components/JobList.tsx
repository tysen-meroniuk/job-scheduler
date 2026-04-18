import { useQuery } from '@tanstack/react-query'
import { api } from '../api/client'

export default function JobList() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['jobs'],
    queryFn: api.listJobs,
  })

  if (isLoading) return <div>Loading…</div>
  if (error) return <div className="text-red-600">Error: {String(error)}</div>

  return (
    <section>
      <h2 className="text-lg font-semibold mb-2">Recent Jobs</h2>
      <div className="bg-white rounded border">
        {!data?.length ? (
          <div className="p-4 text-gray-500">No jobs yet.</div>
        ) : (
          <table className="w-full text-sm">
            <thead className="text-left border-b bg-gray-50">
              <tr>
                <th className="p-2">ID</th>
                <th className="p-2">Queue</th>
                <th className="p-2">Type</th>
                <th className="p-2">State</th>
                <th className="p-2">Attempts</th>
              </tr>
            </thead>
            <tbody>
              {data.map(j => (
                <tr key={j.id} className="border-b last:border-0">
                  <td className="p-2">{j.id}</td>
                  <td className="p-2">{j.queue}</td>
                  <td className="p-2">{j.type}</td>
                  <td className="p-2">{j.state}</td>
                  <td className="p-2">
                    {j.attempts}/{j.max_attempts}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </section>
  )
}
