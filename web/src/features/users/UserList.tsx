import React from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Mail, Calendar, User as UserIcon, Loader2, AlertCircle, Trash2, Edit, Shield } from 'lucide-react';
import apiClient from '../../api/client';
import type { User } from '../../types/api';
import { useAuth } from '../../contexts/AuthContext';
import styles from './UserList.module.css';

interface PaginatedUsers {
  data: User[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

const UserList: React.FC = () => {
  const queryClient = useQueryClient();
  const { user: currentUser } = useAuth();
  
  const { data: usersData, isLoading, error } = useQuery({
    queryKey: ['users'],
    queryFn: async () => {
      // The backend returns { data, total, ... }
      const response = await apiClient.get<PaginatedUsers>('/users');
      return response.data;
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      await apiClient.delete(`/users/${id}`);
    },
    onSuccess: () => {
      // Invalidate and refetch users list
      void queryClient.invalidateQueries({ queryKey: ['users'] });
    },
    onError: (err: any) => {
       alert(err.response?.data?.error || "Failed to delete user. You might not have permission.");
    }
  });

  const handleDelete = (id: number, username: string) => {
    // Prevent deleting self
    if (currentUser?.id === id) {
      alert("You cannot delete your own account.");
      return;
    }

    if (window.confirm(`Are you sure you want to delete user "${username}"?`)) {
      deleteMutation.mutate(id);
    }
  };

  if (isLoading) {
    return (
      <div className={styles.center} style={{ height: '50vh', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
        <Loader2 className={styles.spin} size={32} />
        <p style={{ marginTop: '1rem', color: 'var(--color-text-muted)' }}>Loading users...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.center} style={{ height: '50vh', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
        <AlertCircle color="var(--color-danger)" size={48} />
        <h3 style={{ marginTop: '1rem' }}>Access Denied</h3>
        <p style={{ marginTop: '0.5rem', color: 'var(--color-text-muted)' }}>
          You do not have permission to view the user list.
        </p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
        <div>
          <h2>Users Management</h2>
          <p style={{ color: 'var(--color-text-muted)' }}>View and manage registered users.</p>
        </div>
        {currentUser?.role === 'admin' && (
           <button className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
             <UserIcon size={16} /> Add User
           </button>
        )}
      </div>

      <div className="card" style={{ padding: 0, overflow: 'hidden' }}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>User</th>
              <th>Email</th>
              <th>Role</th>
              <th>Joined Date</th>
              <th style={{ textAlign: 'right' }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {usersData?.data.map((user) => (
              <tr key={user.id} style={{ opacity: deleteMutation.isPending && deleteMutation.variables === user.id ? 0.5 : 1 }}>
                <td>
                  <div className={styles.userCell} style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                    <div className={styles.avatar} style={{ width: '32px', height: '32px', borderRadius: '50%', background: 'var(--color-bg-secondary)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                      <UserIcon size={16} color="var(--color-text-muted)" />
                    </div>
                    <span className={styles.username} style={{ fontWeight: 500 }}>{user.username}</span>
                  </div>
                </td>
                <td>
                  <div className={styles.emailCell} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--color-text-muted)' }}>
                    <Mail size={14} />
                    <span>{user.email}</span>
                  </div>
                </td>
                <td>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                    {user.role === 'admin' ? (
                      <span className={styles.badge} style={{ background: '#dbeafe', color: '#1e40af', padding: '0.25rem 0.5rem', borderRadius: '9999px', fontSize: '0.75rem', fontWeight: 600, display: 'inline-flex', alignItems: 'center', gap: '0.25rem' }}>
                        <Shield size={12} /> Admin
                      </span>
                    ) : (
                      <span className={styles.badge} style={{ background: '#f3f4f6', color: '#374151', padding: '0.25rem 0.5rem', borderRadius: '9999px', fontSize: '0.75rem', fontWeight: 600 }}>
                        User
                      </span>
                    )}
                  </div>
                </td>
                <td>
                  <div className={styles.dateCell} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--color-text-muted)' }}>
                    <Calendar size={14} />
                    <span>{new Date(user.created_at).toLocaleDateString()}</span>
                  </div>
                </td>
                <td style={{ textAlign: 'right' }}>
                  <div className={styles.actionsCell} style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.5rem' }}>
                    {currentUser?.role === 'admin' && (
                      <button 
                        className={styles.deleteBtn} 
                        title="Delete User"
                        onClick={() => handleDelete(user.id, user.username)}
                        disabled={(deleteMutation.isPending && deleteMutation.variables === user.id) || user.id === currentUser.id}
                        style={{ border: 'none', background: 'transparent', cursor: 'pointer', color: 'var(--color-danger)', padding: '0.25rem' }}
                      >
                        {deleteMutation.isPending && deleteMutation.variables === user.id ? (
                          <Loader2 size={16} className={styles.spin} />
                        ) : (
                          <Trash2 size={16} />
                        )}
                      </button>
                    )}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        
        {usersData?.data.length === 0 && (
          <div className={styles.empty} style={{ padding: '2rem', textAlign: 'center', color: 'var(--color-text-muted)' }}>
            No users found.
          </div>
        )}
      </div>
    </div>
  );
};

export default UserList;
