import { Navigate, Route, Routes } from 'react-router-dom'
import { ErrorBoundary, ErrorDisplay, NetworkStatus } from './components/ErrorDisplay/ErrorDisplay.jsx'
import Layout from './components/Layout/Layout.jsx'
import { NotificationProvider } from './components/NotificationProvider.jsx'
import { ErrorProvider } from './contexts/ErrorContext.jsx'
import AttendeeManagement from './pages/AttendeeManagement/AttendeeManagement.jsx'
import FeatureManagement from './pages/FeatureManagement/FeatureManagement.jsx'
import FibonacciScoring from './pages/FibonacciScoring.jsx'
import PairwiseComparison from './pages/PairwiseComparison.jsx'
import ProjectList from './pages/ProjectList/ProjectList.jsx'
import ProjectSetup from './pages/ProjectSetup/ProjectSetup.jsx'
import Results from './pages/Results.jsx'

function App() {
  return (
    <ErrorBoundary>
      <ErrorProvider>
        <NotificationProvider>
          <Layout>
            <Routes>
              <Route path="/" element={<Navigate to="/projects" replace />} />
              <Route path="/projects" element={<ProjectList />} />
              <Route path="/projects/new" element={<ProjectSetup />} />
              <Route path="/projects/:id/edit" element={<ProjectSetup />} />
              <Route path="/projects/:id/attendees" element={<AttendeeManagement />} />
              <Route path="/projects/:id/features" element={<FeatureManagement />} />
              <Route path="/projects/:projectId/comparison" element={<PairwiseComparison />} />
              <Route path="/projects/:projectId/scoring/value" element={<FibonacciScoring />} />
              <Route path="/projects/:projectId/scoring/complexity" element={<FibonacciScoring />} />
              <Route path="/projects/:projectId/results" element={<Results />} />
              <Route path="*" element={<Navigate to="/projects" replace />} />
            </Routes>
          </Layout>
          <ErrorDisplay />
          <NetworkStatus />
        </NotificationProvider>
      </ErrorProvider>
    </ErrorBoundary>
  )
}

export default App