import { createContext, useContext, useState, useCallback, ReactNode, useEffect } from 'react';
import { EventsOn } from '../../wailsjs/runtime/runtime';

export interface ApiLogEntry {
    id: string;
    timestamp: Date;
    type: 'lcu' | 'riot';
    method: string;
    endpoint: string;
    status: 'pending' | 'success' | 'error';
    duration?: number;
    headers?: Record<string, string>;
    response?: any;
    error?: string;
}

interface ApiLogContextType {
    logs: ApiLogEntry[];
}

const ApiLogContext = createContext<ApiLogContextType | null>(null);

const MAX_LOGS = 50;

export function ApiLogProvider({ children }: { children: ReactNode }) {
    const [logs, setLogs] = useState<ApiLogEntry[]>([]);

    // Subscribe to backend API call events
    useEffect(() => {
        const off = EventsOn('api-call', (data: any) => {
            if (!data) return;

            const id = `${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
            const type = (data.type === 'lcu' || data.type === 'riot') ? data.type : 'lcu';
            const statusCode: number | undefined = typeof data.statusCode === 'number' ? data.statusCode : undefined;
            const durationMs: number | undefined =
                typeof data.duration === 'number' ? data.duration : (typeof data.duration === 'string' ? parseInt(data.duration, 10) : undefined);

            const status: ApiLogEntry['status'] =
                data.error ? 'error' : (statusCode && statusCode >= 200 && statusCode < 300 ? 'success' : 'pending');

            const newEntry: ApiLogEntry = {
                id,
                timestamp: new Date(),
                type,
                method: data.method || 'GET',
                endpoint: data.endpoint || '',
                status,
                duration: durationMs,
                headers: data.headers || {},
                response: data.response,
                error: data.error,
            };

            setLogs(prev => [newEntry, ...prev].slice(0, MAX_LOGS));
        });

        return () => {
            if (typeof off === 'function') {
                off();
            }
        };
    }, []);

    return (
        <ApiLogContext.Provider value={{ logs }}>
            {children}
        </ApiLogContext.Provider>
    );
}

export function useApiLog(): ApiLogContextType {
    const context = useContext(ApiLogContext);
    if (!context) {
        throw new Error('useApiLog must be used within ApiLogProvider');
    }
    return context;
}

