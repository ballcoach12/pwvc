import { Routes, Route, Navigate } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import ProjectList from './pages/ProjectList/ProjectList';
import ProjectSetup from './pages/ProjectSetup/ProjectSetup';
import AttendeeManagement from './pages/AttendeeManagement/AttendeeManagement';
import FeatureManagement from './pages/FeatureManagement/FeatureManagement';
import PairwiseComparison from './pages/PairwiseComparison';
import FibonacciScoring from './pages/FibonacciScoring';
import Results from './pages/Results';
import ProtectedRoute from './components/auth/ProtectedRoute';

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/projects" element={<ProtectedRoute><ProjectList /></ProtectedRoute>} />
      <Route path="/projects/new" element={<ProtectedRoute><ProjectSetup /></ProtectedRoute>} />
      <Route path="/projects/:id/edit" element={<ProtectedRoute><ProjectSetup /></ProtectedRoute>} />
      <Route path="/projects/:id/attendees" element={<ProtectedRoute><AttendeeManagement /></ProtectedRoute>} />
      <Route path="/projects/:id/features" element={<ProtectedRoute><FeatureManagement /></ProtectedRoute>} />
      <Route path="/projects/:id/pairwise" element={<ProtectedRoute><PairwiseComparison /></ProtectedRoute>} />
      <Route path="/projects/:id/fibonacci" element={<ProtectedRoute><FibonacciScoring /></ProtectedRoute>} />
      <Route path="/projects/:id/results" element={<ProtectedRoute><Results /></ProtectedRoute>} />
      <Route path="/" element={<Navigate to="/projects" replace />} />
      <Route path="*" element={<Navigate to="/projects" replace />} />
    </Routes>
  );
}

export default App
