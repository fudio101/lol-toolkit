import { createContext, useContext, useState, useEffect, useCallback, useRef, ReactNode } from 'react';
import { GetLCUStatus, GetCurrentSummoner } from "../../wailsjs/go/app/App";
import { lcu, app } from "../../wailsjs/go/models";

// Polling intervals in milliseconds
const POLL_INTERVAL = {
    DISCONNECTED: 10_000,  // 10 seconds when not connected
    CONNECTED: 30_000,     // 30 seconds when connected
} as const;

// Polling indicator duration
const POLLING_INDICATOR_MS = 500;

interface LCUContextType {
    status: app.LCUStatus | null;
    summoner: lcu.CurrentSummoner | null;
    loading: boolean;
    error: string | null;
    isPolling: boolean;
    refresh: () => Promise<void>;
}

const LCUContext = createContext<LCUContextType | null>(null);

export function LCUProvider({ children }: { children: ReactNode }) {
    const [status, setStatus] = useState<app.LCUStatus | null>(null);
    const [summoner, setSummoner] = useState<lcu.CurrentSummoner | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isPolling, setIsPolling] = useState(false);
    const isRefreshing = useRef(false);

    const refresh = useCallback(async () => {
        // Prevent concurrent refreshes
        if (isRefreshing.current) return;
        isRefreshing.current = true;
        setIsPolling(true);

        try {
            setError(null);
            const lcuStatus = await GetLCUStatus();
            setStatus(lcuStatus);

            if (lcuStatus.connected) {
                try {
                    const summonerData = await GetCurrentSummoner();
                    setSummoner(summonerData);
                } catch (err) {
                    console.error('Failed to get summoner:', err);
                    setSummoner(null);
                }
            } else {
                setSummoner(null);
            }
        } catch (err) {
            const errorMsg = err instanceof Error ? err.message : String(err);
            setError(errorMsg);
            setStatus({ connected: false, error: errorMsg });
            setSummoner(null);
        } finally {
            setLoading(false);
            isRefreshing.current = false;
            setTimeout(() => setIsPolling(false), POLLING_INDICATOR_MS);
        }
    }, []);

    // Initial fetch
    useEffect(() => {
        refresh();
    }, [refresh]);

    // Adaptive polling
    useEffect(() => {
        const interval = status?.connected 
            ? POLL_INTERVAL.CONNECTED 
            : POLL_INTERVAL.DISCONNECTED;

        const timer = setInterval(refresh, interval);
        return () => clearInterval(timer);
    }, [refresh, status?.connected]);

    const value: LCUContextType = {
        status,
        summoner,
        loading,
        error,
        isPolling,
        refresh,
    };

    return (
        <LCUContext.Provider value={value}>
            {children}
        </LCUContext.Provider>
    );
}

export function useLCU(): LCUContextType {
    const context = useContext(LCUContext);
    if (!context) {
        throw new Error('useLCU must be used within LCUProvider');
    }
    return context;
}
