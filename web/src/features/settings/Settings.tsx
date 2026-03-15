import React from 'react';
import { User, Mail, Shield, Calendar, LogOut } from 'lucide-react';
import { useAuth } from '../../contexts/AuthContext';
import styles from './Settings.module.css';

const Settings: React.FC = () => {
  const { user, logout } = useAuth();

  if (!user) return null;

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2>Account Settings</h2>
        <p>Manage your profile and preferences.</p>
      </div>

      <div className={styles.grid}>
        {/* Profile Card */}
        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <h3>Profile Information</h3>
          </div>
          <div className={styles.profileSection}>
            <div className={styles.avatar}>
              <User size={48} color="#666" />
            </div>
            <div className={styles.profileInfo}>
              <div className={styles.infoGroup}>
                <label>Username</label>
                <div className={styles.valueWithIcon}>
                  <User size={16} />
                  <span>{user.username}</span>
                </div>
              </div>
              <div className={styles.infoGroup}>
                <label>Email Address</label>
                <div className={styles.valueWithIcon}>
                  <Mail size={16} />
                  <span>{user.email}</span>
                </div>
              </div>
              <div className={styles.infoGroup}>
                <label>Role</label>
                <div className={styles.valueWithIcon}>
                  <Shield size={16} />
                  <span className={styles.roleBadge}>{user.role}</span>
                </div>
              </div>
              <div className={styles.infoGroup}>
                <label>Member Since</label>
                <div className={styles.valueWithIcon}>
                  <Calendar size={16} />
                  <span>{new Date(user.created_at).toLocaleDateString()}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Security / Actions Card */}
        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <h3>Account Actions</h3>
          </div>
          <div className={styles.actionsList}>
            <button className={styles.actionBtn} onClick={logout}>
              <LogOut size={18} />
              <span>Sign Out</span>
            </button>
            {/* Future: Change Password, etc. */}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Settings;
