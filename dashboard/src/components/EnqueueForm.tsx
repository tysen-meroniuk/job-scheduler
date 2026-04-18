import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '../api/client'

export default function EnqueueForm() {
  const [type, setType] = useState('echo')
  const [payload, setPayload] = useState('{"msg":"hello"}')
  const qc = useQueryClient()

  const enqueue = useMutation({
    mutationFn: () => api.enqueue({ type, payload: JSON.parse(payload) }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['jobs'] })
      qc.invalidateQueries({ queryKey: ['metrics'] })
    },
  })

  return (
    <section className="bg-white rounded border p-4">
      <h2 className="text-lg font-semibold mb-2">Enqueue a Job</h2>
      <form
        onSubmit={e => {
          e.preventDefault()
          enqueue.mutate()
        }}
        className="flex gap-2 items-start"
      >
        <input
          className="border rounded px-2 py-1 w-32"
          value={type}
          onChange={e => setType(e.target.value)}
          placeholder="type"
        />
        <textarea
          className="border rounded px-2 py-1 flex-1 font-mono text-sm"
          rows={2}
          value={payload}
          onChange={e => setPayload(e.target.value)}
        />
        <button
          type="submit"
          disabled={enqueue.isPending}
          className="bg-blue-600 text-white px-3 py-1 rounded disabled:opacity-50"
        >
          Enqueue
        </button>
      </form>
      {enqueue.error && (
        <div className="text-red-600 text-sm mt-2">
          {String(enqueue.error)}
        </div>
      )}
    </section>
  )
}
