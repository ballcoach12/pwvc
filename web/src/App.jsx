import { Navigate, Route, Routes } from 'react-router-dom'
import Layout from './components/Layout/Layout.jsx'
import { NotificationProvider } from './components/NotificationProvider.jsx'
import AttendeeManagement from './pages/AttendeeManagement/AttendeeManagement.jsx'
import FeatureManagement from './pages/FeatureManagement/FeatureManagement.jsx'
import ProjectList from './pages/ProjectList/ProjectList.jsx'
import ProjectSetup from './pages/ProjectSetup/ProjectSetup.jsx'

function App() {
  return (
    <NotificationProvider>
      <Layout>
        <Routes>
          <Route path="/" element={<Navigate to="/projects" replace />} />
          <Route path="/projects" element={<ProjectList />} />
          <Route path="/projects/new" element={<ProjectSetup />} />
          <Route path="/projects/:id/edit" element={<ProjectSetup />} />
          <Route path="/projects/:id/attendees" element={<AttendeeManagement />} />
          <Route path="/projects/:id/features" element={<FeatureManagement />} />
          {/* TODO: Add more routes for pairwise comparison and results */}
          <Route path="*" element={<Navigate to="/projects" replace />} />
        </Routes>
      </Layout>
    </NotificationProvider>
  )
}

export default App