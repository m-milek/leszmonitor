import {createFileRoute, Navigate} from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  // From the root route, I redirect to the dashboard route.
  return (
    <Navigate to="/dashboard" replace={true} />
  )
}