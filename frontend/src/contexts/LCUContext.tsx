import { createContext, useContext, useState, useEffect, useCallback, useRef, ReactNode } from 'react';
import { GetLCUStatus, GetCurrentSummoner } from "../../wailsjs/go/app/App";
import { lcu, app } from "../../wailsjs/go/models";
import { useApiLog } from './ApiLogContext';

// Polling intervals
const POLL_INTERVAL = {
    STATUS_DISCONNECTED: 10_000,  // Check LCU status every 10s when not connected
    STATUS_CONNECTED: 30_000,     // Check LCU status every 30s when connected
    SUMMONER_INITIAL: 5_000,      // Retry summoner fetch every 5s until we get it
    SUMMONER_REFRESH: 60_000,     // Refresh summoner data every 60s after initial fetch
} as const;

const POLLING_INDICATOR_MS = 500;

interface LCUContextType {
    status: app.LCUStatus | null;
    summoner: lcu.CurrentSummoner | null;
    loading: boolean;
    error: string | null;
    isPolling: boolean;
    refresh: () => Promise<void>;
    refreshSummoner: () => Promise<void>;
}

const LCUContext = createContext<LCUContextType | null>(null);

export function LCUProvider({ children }: { children: ReactNode }) {
    const [status, setStatus] = useState<app.LCUStatus | null>(null);
    const [summoner, setSummoner] = useState<lcu.CurrentSummoner | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isPolling, setIsPolling] = useState(false);
    const isRefreshing = useRef(false);
    const wasConnected = useRef(false);
    const { addLog, updateLog } = useApiLog();

    // Check LCU status only
    const checkStatus = useCallback(async () => {
        if (isRefreshing.current) return;
        isRefreshing.current = true;
        setIsPolling(true);

        const startTime = Date.now();
        const logId = addLog({
            type: 'lcu',
            method: 'GET',
            endpoint: 'GetLCUStatus',
            status: 'pending',
        });

        try {
            setError(null);
            const lcuStatus = await GetLCUStatus();
            updateLog(logId, {
                status: 'success',
                duration: Date.now() - startTime,
                response: lcuStatus,
            });
            setStatus(lcuStatus);

            // If just connected, fetch summoner data
            if (lcuStatus.connected && !wasConnected.current) {
                wasConnected.current = true;
                await fetchSummoner();
            } else if (!lcuStatus.connected) {
                wasConnected.current = false;
                setSummoner(null);
            }
        } catch (err) {
            const errorMsg = err instanceof Error ? err.message : String(err);
            updateLog(logId, {
                status: 'error',
                duration: Date.now() - startTime,
                error: errorMsg,
            });
            setError(errorMsg);
            setStatus({ connected: false, error: errorMsg });
            setSummoner(null);
            wasConnected.current = false;
        } finally {
            setLoading(false);
            isRefreshing.current = false;
            setTimeout(() => setIsPolling(false), POLLING_INDICATOR_MS);
        }
    }, [addLog, updateLog]);

    // Fetch summoner data
    const fetchSummoner = useCallback(async () => {
        if (!status?.connected && !wasConnected.current) return;

        const startTime = Date.now();
        const logId = addLog({
            type: 'lcu',
            method: 'GET',
            endpoint: 'GetCurrentSummoner',
            status: 'pending',
        });

        try {
            const summonerData = await GetCurrentSummoner();
            updateLog(logId, {
                status: 'success',
                duration: Date.now() - startTime,
                response: summonerData,
            });
            setSummoner(summonerData);
        } catch (err) {
            const errorMsg = err instanceof Error ? err.message : String(err);
            updateLog(logId, {
                status: 'error',
                duration: Date.now() - startTime,
                error: errorMsg,
            });
            console.error('Failed to get summoner:', err);
            setSummoner(null);
        }
    }, [status?.connected, addLog, updateLog]);

    // Manual refresh - checks status and fetches summoner
    const refresh = useCallback(async () => {
        await checkStatus();
    }, [checkStatus]);

    // Manual summoner refresh
    const refreshSummoner = useCallback(async () => {
        if (status?.connected) {
            await fetchSummoner();
        }
    }, [status?.connected, fetchSummoner]);

    // Initial fetch
    useEffect(() => {
        checkStatus();
    }, [checkStatus]);

    // Poll status with adaptive interval
    useEffect(() => {
        const interval = status?.connected 
            ? POLL_INTERVAL.STATUS_CONNECTED 
            : POLL_INTERVAL.STATUS_DISCONNECTED;
        const timer = setInterval(checkStatus, interval);
        return () => clearInterval(timer);
    }, [checkStatus, status?.connected]);

    // Poll summoner with adaptive interval
    useEffect(() => {
        if (!status?.connected) return;

        // Use faster interval until we have summoner data
        const interval = summoner 
            ? POLL_INTERVAL.SUMMONER_REFRESH 
            : POLL_INTERVAL.SUMMONER_INITIAL;
        const timer = setInterval(fetchSummoner, interval);
        return () => clearInterval(timer);
    }, [status?.connected, summoner, fetchSummoner]);

    const value: LCUContextType = {
        status,
        summoner,
        loading,
        error,
        isPolling,
        refresh,
        refreshSummoner,
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
