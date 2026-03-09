import React from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Mail, Calendar, User as UserIcon, Loader2, AlertCircle, Trash2, Edit } from 'lucide-react';
import apiClient from '../../api/client';
import type { User } from '../../types/api';
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
  
  const { data, isLoading, error } = useQuery({
    queryKey: ['users'],
    queryFn: async () => {
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
  });

  const handleDelete = (id: number, username: string) => {
    if (globalThis.confirm(`Are you sure you want to delete user "${username}"?`)) {
      deleteMutation.mutate(id);
    }
  };

  if (isLoading) {
    return (
      <div className={styles.center}>
        <Loader2 className={styles.spin} size={32} />
        <p>Loading users...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.center}>
        <AlertCircle color="var(--color-danger)" size={32} />
        <p>Failed to load users. Please try again later.</p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div>
          <h2>Users</h2>
          <p>Manage and view all registered users in the system.</p>
        </div>
        <button className="btn-primary">Add User</button>
      </div>

      <div className="card" style={{ padding: 0, overflow: 'hidden' }}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>User</th>
              <th>Email</th>
              <th>Joined Date</th>
              <th>Status</th>
              <th style={{ textAlign: 'right' }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {data?.data.map((user) => (
              <tr key={user.id} className={deleteMutation.variables === user.id ? styles.deleting : ''}>
                <td>
                  <div className={styles.userCell}>
                    <div className={styles.avatar}>
                      <UserIcon size={16} />
                    </div>
                    <span className={styles.username}>{user.username}</span>
                  </div>
                </td>
                <td>
                  <div className={styles.emailCell}>
                    <Mail size={14} />
                    <span>{user.email}</span>
                  </div>
                </td>
                <td>
                  <div className={styles.dateCell}>
                    <Calendar size={14} />
                    <span>{new Date(user.created_at).toLocaleDateString()}</span>
                  </div>
                </td>
                <td>
                  <span className={styles.badge}>Active</span>
                </td>
                <td style={{ textAlign: 'right' }}>
                  <div className={styles.actionsCell}>
                    <button className={styles.editBtn} title="Edit User">
                      <Edit size={16} />
                    </button>
                    <button 
                      className={styles.deleteBtn} 
                      title="Delete User"
                      onClick={() => handleDelete(user.id, user.username)}
                      disabled={deleteMutation.isPending && deleteMutation.variables === user.id}
                    >
                      {deleteMutation.isPending && deleteMutation.variables === user.id ? (
                        <Loader2 size={16} className={styles.spin} />
                      ) : (
                        <Trash2 size={16} />
                      )}
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        
        {data?.data.length === 0 && (
          <div className={styles.empty}>
            No users found.
          </div>
        )}
      </div>
    </div>
  );
};

export default UserList;
