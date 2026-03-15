import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './contexts/AuthContext';
import DashboardLayout from './layouts/DashboardLayout';
import Login from './features/auth/Login';
import Signup from './features/auth/Signup';
import ForgotPassword from './features/auth/ForgotPassword';
import ResetPassword from './features/auth/ResetPassword';
import UserList from './features/users/UserList';
import Settings from './features/settings/Settings';
import HealthStatus from './features/health/HealthStatus';

const DashboardHome = () => (
  <div>
    <h2 style={{ marginBottom: '1.5rem' }}>System Overview</h2>
    <HealthStatus />
    
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '1.5rem' }}>
      <div className="card">
        <h3>Total Users</h3>
        <p style={{ fontSize: '2rem', fontWeight: 700, marginTop: '0.5rem' }}>1,284</p>
      </div>
      <div className="card">
        <h3>Active Now</h3>
        <p style={{ fontSize: '2rem', fontWeight: 700, marginTop: '0.5rem' }}>42</p>
      </div>
      <div className="card">
        <h3>API Requests</h3>
        <p style={{ fontSize: '2rem', fontWeight: 700, marginTop: '0.5rem', color: 'var(--color-primary)' }}>12.5k</p>
      </div>
    </div>
  </div>
);

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
};

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/reset-password" element={<ResetPassword />} />
        
        <Route path="/" element={
          <ProtectedRoute>
            <DashboardLayout />
          </ProtectedRoute>
        }>
          <Route index element={<DashboardHome />} />
          <Route path="users" element={<UserList />} />
          <Route path="settings" element={<Settings />} />
        </Route>

        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
