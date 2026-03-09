import React from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Users, 
  Settings, 
  LogOut, 
  Activity,
  User as UserIcon
} from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import styles from './DashboardLayout.module.css';
import { clsx } from 'clsx';

const DashboardLayout: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className={styles.container}>
      {/* Sidebar */}
      <aside className={styles.sidebar}>
        <div className={styles.logo}>
          <Activity size={24} color="var(--color-primary)" />
          <span>MicroCore</span>
        </div>

        <nav className={styles.nav}>
          <NavLink 
            to="/" 
            className={({ isActive }) => clsx(styles.navLink, isActive && styles.active)}
          >
            <LayoutDashboard size={20} />
            <span>Dashboard</span>
          </NavLink>
          <NavLink 
            to="/users" 
            className={({ isActive }) => clsx(styles.navLink, isActive && styles.active)}
          >
            <Users size={20} />
            <span>Users</span>
          </NavLink>
          <NavLink 
            to="/settings" 
            className={({ isActive }) => clsx(styles.navLink, isActive && styles.active)}
          >
            <Settings size={20} />
            <span>Settings</span>
          </NavLink>
        </nav>

        <div className={styles.sidebarFooter}>
          <button onClick={handleLogout} className={styles.logoutBtn}>
            <LogOut size={20} />
            <span>Logout</span>
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className={styles.main}>
        <header className={styles.header}>
          <div className={styles.headerContent}>
            <h1>Dashboard</h1>
            <div className={styles.userProfile}>
              <div className={styles.userInfo}>
                <span className={styles.username}>{user?.username}</span>
                <span className={styles.userEmail}>{user?.email}</span>
              </div>
              <div className={styles.avatar}>
                <UserIcon size={20} />
              </div>
            </div>
          </div>
        </header>

        <div className={styles.content}>
          <Outlet />
        </div>
      </main>
    </div>
  );
};

export default DashboardLayout;
