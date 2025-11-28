import { createContext, useContext, useState, useCallback, ReactNode } from 'react';

export interface ApiLogEntry {
    id: string;
    timestamp: Date;
    type: 'lcu' | 'riot';
    method: string;
    endpoint: string;
    status: 'pending' | 'success' | 'error';
    duration?: number;
    response?: any;
    error?: string;
}

interface ApiLogContextType {
    logs: ApiLogEntry[];
    addLog: (entry: Omit<ApiLogEntry, 'id' | 'timestamp'>) => string;
    updateLog: (id: string, updates: Partial<ApiLogEntry>) => void;
    clearLogs: () => void;
}

const ApiLogContext = createContext<ApiLogContextType | null>(null);

const MAX_LOGS = 50;

export function ApiLogProvider({ children }: { children: ReactNode }) {
    const [logs, setLogs] = useState<ApiLogEntry[]>([]);

    const addLog = useCallback((entry: Omit<ApiLogEntry, 'id' | 'timestamp'>): string => {
        const id = `${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
        const newEntry: ApiLogEntry = {
            ...entry,
            id,
            timestamp: new Date(),
        };

        setLogs(prev => [newEntry, ...prev].slice(0, MAX_LOGS));
        return id;
    }, []);

    const updateLog = useCallback((id: string, updates: Partial<ApiLogEntry>) => {
        setLogs(prev => prev.map(log => 
            log.id === id ? { ...log, ...updates } : log
        ));
    }, []);

    const clearLogs = useCallback(() => {
        setLogs([]);
    }, []);

    return (
        <ApiLogContext.Provider value={{ logs, addLog, updateLog, clearLogs }}>
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

