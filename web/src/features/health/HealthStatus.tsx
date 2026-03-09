import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Activity, Database, Zap, Loader2, CheckCircle2, XCircle } from 'lucide-react';
import apiClient from '../../api/client';
import styles from './HealthStatus.module.css';

interface HealthResponse {
  status: 'UP' | 'DOWN';
  details: {
    database: 'ok' | 'disconnected';
    redis: 'ok' | 'disconnected';
  };
  time: string;
}

const HealthStatus: React.FC = () => {
  const { data, isLoading } = useQuery({
    queryKey: ['health'],
    queryFn: async () => {
      const response = await apiClient.get<HealthResponse>('/health', {
        baseURL: 'http://localhost:3003' // Health is outside /api/v1 usually
      });
      return response.data;
    },
    refetchInterval: 10000, // Poll every 10 seconds
  });

  if (isLoading) {
    return (
      <div className={styles.container}>
        <Loader2 className={styles.spin} size={20} />
        <span>Checking system health...</span>
      </div>
    );
  }

  const isDBUp = data?.details.database === 'ok';
  const isRedisUp = data?.details.redis === 'ok';

  return (
    <div className={styles.grid}>
      <div className={styles.card}>
        <div className={styles.cardHeader}>
          <Activity size={18} />
          <span>API Server</span>
        </div>
        <div className={styles.status}>
          {data?.status === 'UP' ? (
            <CheckCircle2 color="var(--color-success)" size={20} />
          ) : (
            <XCircle color="var(--color-danger)" size={20} />
          )}
          <span className={data?.status === 'UP' ? styles.up : styles.down}>
            {data?.status || 'UNKNOWN'}
          </span>
        </div>
      </div>

      <div className={styles.card}>
        <div className={styles.cardHeader}>
          <Database size={18} />
          <span>PostgreSQL</span>
        </div>
        <div className={styles.status}>
          {isDBUp ? (
            <CheckCircle2 color="var(--color-success)" size={20} />
          ) : (
            <XCircle color="var(--color-danger)" size={20} />
          )}
          <span className={isDBUp ? styles.up : styles.down}>
            {isDBUp ? 'CONNECTED' : 'DISCONNECTED'}
          </span>
        </div>
      </div>

      <div className={styles.card}>
        <div className={styles.cardHeader}>
          <Zap size={18} />
          <span>Redis Cache</span>
        </div>
        <div className={styles.status}>
          {isRedisUp ? (
            <CheckCircle2 color="var(--color-success)" size={20} />
          ) : (
            <XCircle color="var(--color-danger)" size={20} />
          )}
          <span className={isRedisUp ? styles.up : styles.down}>
            {isRedisUp ? 'CONNECTED' : 'DISCONNECTED'}
          </span>
        </div>
      </div>
    </div>
  );
};

export default HealthStatus;
