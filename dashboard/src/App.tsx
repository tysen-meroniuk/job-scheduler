import JobList from './components/JobList'
import EnqueueForm from './components/EnqueueForm'
import Metrics from './components/Metrics'

export default function App() {
  return (
    <div className="min-h-screen bg-gray-50 text-gray-900">
      <header className="bg-white border-b px-6 py-4">
        <h1 className="text-xl font-semibold">Job Queue</h1>
      </header>
      <main className="max-w-5xl mx-auto p-6 space-y-6">
        <Metrics />
        <EnqueueForm />
        <JobList />
      </main>
    </div>
  )
}
