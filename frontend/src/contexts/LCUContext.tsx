import { createContext, useContext, useState, useCallback, useRef, ReactNode, useEffect } from 'react';
import { GetLCUStatus, GetCurrentSummoner } from "../../wailsjs/go/app/App";
import { lcu, app } from "../../wailsjs/go/models";
import { useApiLog } from './ApiLogContext';
import { usePolling } from '../hooks';

// Polling intervals (ms)
const INTERVAL = {
    STATUS_DISCONNECTED: 10_000,  // 10s - check if client started
    STATUS_CONNECTED: 30_000,     // 30s - normal connected polling
    SUMMONER_INITIAL: 5_000,      // 5s - waiting for summoner data
    SUMMONER_REFRESH: 60_000,     // 60s - periodic summoner refresh
} as const;

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
    
    // Refs for tracking state and function references
    const statusRef = useRef(status);
    statusRef.current = status;
    
    const prevConnectedRef = useRef<boolean | null>(null);
    const summonerRefreshRef = useRef<(() => Promise<void>) | null>(null);
    
    const { addLog, updateLog } = useApiLog();

    // ============================================================================
    // Status Polling
    // ============================================================================

    const fetchStatus = useCallback(async (): Promise<app.LCUStatus> => {
        const startTime = Date.now();
        const logId = addLog({
            type: 'lcu',
            method: 'GET',
            endpoint: 'GetLCUStatus',
            status: 'pending',
        });

        try {
            const result = await GetLCUStatus();
            updateLog(logId, {
                status: 'success',
                duration: Date.now() - startTime,
                response: result,
            });
            return result;
        } catch (err) {
            const errorMsg = err instanceof Error ? err.message : String(err);
            updateLog(logId, {
                status: 'error',
                duration: Date.now() - startTime,
                error: errorMsg,
            });
            throw err;
        }
    }, [addLog, updateLog]);

    const getStatusInterval = useCallback((result: app.LCUStatus | null): number => {
        return result?.connected 
            ? INTERVAL.STATUS_CONNECTED 
            : INTERVAL.STATUS_DISCONNECTED;
    }, []);

    const handleStatusResult = useCallback((result: app.LCUStatus) => {
        setStatus(result);
        setError(null);
        setLoading(false);
        
        // Clear summoner if disconnected
        if (!result.connected) {
            setSummoner(null);
        }
    }, []);

    const handleStatusError = useCallback((err: Error) => {
        setError(err.message);
        setStatus({ connected: false, error: err.message });
        setSummoner(null);
        setLoading(false);
    }, []);

    const { isPolling: isPollingStatus, refresh: refreshStatus } = usePolling({
        id: 'lcu-status',
        enabled: true,
        execute: fetchStatus,
        getInterval: getStatusInterval,
        onResult: handleStatusResult,
        onError: handleStatusError,
    });

    // ============================================================================
    // Summoner Polling
    // ============================================================================

    const fetchSummoner = useCallback(async (): Promise<lcu.CurrentSummoner | null> => {
        // Only fetch if connected
        if (!statusRef.current?.connected) {
            return null;
        }

        const startTime = Date.now();
        const logId = addLog({
            type: 'lcu',
            method: 'GET',
            endpoint: 'GetCurrentSummoner',
            status: 'pending',
        });

        try {
            const result = await GetCurrentSummoner();
            updateLog(logId, {
                status: 'success',
                duration: Date.now() - startTime,
                response: result,
            });
            return result;
        } catch (err) {
            const errorMsg = err instanceof Error ? err.message : String(err);
            updateLog(logId, {
                status: 'error',
                duration: Date.now() - startTime,
                error: errorMsg,
            });
            return null;
        }
    }, [addLog, updateLog]);

    const getSummonerInterval = useCallback((result: lcu.CurrentSummoner | null): number => {
        if (!statusRef.current?.connected) {
            return INTERVAL.STATUS_DISCONNECTED;
        }
        // Faster polling until we have summoner data, then slower refresh
        return result ? INTERVAL.SUMMONER_REFRESH : INTERVAL.SUMMONER_INITIAL;
    }, []);

    const handleSummonerResult = useCallback((result: lcu.CurrentSummoner | null) => {
        if (result) {
            setSummoner(result);
        }
    }, []);

    const { isPolling: isPollingSummoner, refresh: refreshSummoner } = usePolling({
        id: 'lcu-summoner',
        enabled: true,
        execute: fetchSummoner,
        getInterval: getSummonerInterval,
        onResult: handleSummonerResult,
    });

    // ============================================================================
    // Connection Transition Handler
    // ============================================================================

    // Store refresh function in ref for transition handler
    useEffect(() => {
        summonerRefreshRef.current = refreshSummoner;
    }, [refreshSummoner]);

    // Watch for connection state changes and immediately fetch summoner on connect
    useEffect(() => {
        const isNowConnected = status?.connected ?? false;
        const wasConnected = prevConnectedRef.current;

        // Detect transition: (null/false) → true
        // This handles both first-time connection and reconnection
        const isNewlyConnected = (wasConnected === null || wasConnected === false) && isNowConnected === true;

        if (isNewlyConnected && summonerRefreshRef.current) {
            // Immediately fetch summoner and reset polling timer
            summonerRefreshRef.current();
        }

        // Update ref for next check
        prevConnectedRef.current = isNowConnected;
    }, [status?.connected]);

    // ============================================================================
    // Public API
    // ============================================================================

    const refresh = useCallback(async () => {
        await Promise.all([refreshStatus(), refreshSummoner()]);
    }, [refreshStatus, refreshSummoner]);

    const value: LCUContextType = {
        status,
        summoner,
        loading,
        error,
        isPolling: isPollingStatus || isPollingSummoner,
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
