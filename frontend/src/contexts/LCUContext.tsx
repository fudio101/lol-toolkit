import { createContext, useContext, useState, useCallback, useRef, ReactNode, useEffect } from 'react';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { GetLCUStatus, GetCurrentSummoner } from "../../wailsjs/go/app/App";
import { lcu, app } from "../../wailsjs/go/models";
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
    
    const statusRef = useRef(status);
    const prevConnectedRef = useRef<boolean | null>(null);
    const summonerRefreshRef = useRef<(() => Promise<void>) | null>(null);
    const handleStatusResultRef = useRef<((result: app.LCUStatus) => void) | null>(null);

    useEffect(() => {
        statusRef.current = status;
    }, [status]);

    // ============================================================================
    // Status Polling
    // ============================================================================

    const fetchStatus = useCallback(async (): Promise<app.LCUStatus> => {
        return await GetLCUStatus();
    }, []);

    const getStatusInterval = useCallback((result: app.LCUStatus | null): number => {
        return result?.connected 
            ? INTERVAL.STATUS_CONNECTED 
            : INTERVAL.STATUS_DISCONNECTED;
    }, []);

    const handleStatusResult = useCallback((result: app.LCUStatus) => {
        statusRef.current = result;
        
        const wasConnected = prevConnectedRef.current;
        const isNowConnected = result.connected === true;
        const isNewlyConnected = (wasConnected === null || wasConnected === false) && isNowConnected;
        
        prevConnectedRef.current = isNowConnected;
        
        setStatus(result);
        setError(null);
        setLoading(false);
        
        // Handle summoner based on connection
        if (!result.connected) {
            setSummoner(null);
        } else if (isNewlyConnected && summonerRefreshRef.current) {
            summonerRefreshRef.current();
        }
    }, []);

    // Store handleStatusResult in ref for use in fetchSummoner
    useEffect(() => {
        handleStatusResultRef.current = handleStatusResult;
    }, [handleStatusResult]);

    // ============================================================================
    // Summoner Functions
    // ============================================================================

    const fetchSummoner = useCallback(async (): Promise<lcu.CurrentSummoner | null> => {
        if (!statusRef.current?.connected) {
            return null;
        }

        try {
            return await GetCurrentSummoner();
        } catch {
            // Backend handles errors and emits events
            return null;
        }
    }, []);

    const handleSummonerResult = useCallback((result: lcu.CurrentSummoner | null) => {
        if (result) {
            setSummoner(result);
        }
    }, []);

    const handleStatusError = useCallback((err: Error) => {
        setError(err.message);
        const errorStatus: app.LCUStatus = { connected: false, error: err.message };
        setStatus(errorStatus);
        prevConnectedRef.current = false;
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

    const getSummonerInterval = useCallback((result: lcu.CurrentSummoner | null): number => {
        if (!statusRef.current?.connected) {
            return INTERVAL.STATUS_DISCONNECTED;
        }
        // Faster polling until we have summoner data, then slower refresh
        return result ? INTERVAL.SUMMONER_REFRESH : INTERVAL.SUMMONER_INITIAL;
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

    // Backup handler for connection state changes
    useEffect(() => {
        if (!status) return;

        const isNowConnected = status.connected === true;
        const wasConnected = prevConnectedRef.current;
        const isNewlyConnected = (wasConnected === null || wasConnected === false) && isNowConnected;

        if (isNewlyConnected && summonerRefreshRef.current) {
            summonerRefreshRef.current();
        }
    }, [status]);

    // ============================================================================
    // Listen for backend status change events
    // ============================================================================

    useEffect(() => {
        let debounceTimer: ReturnType<typeof setTimeout> | null = null;
        const DEBOUNCE_MS = 500;

        const unsubscribe = EventsOn('lcu-status-changed', (data: any) => {
            if (!data || typeof data.connected !== 'boolean') return;

            if (debounceTimer) {
                clearTimeout(debounceTimer);
            }

            debounceTimer = setTimeout(() => {
                const currentConnected = statusRef.current?.connected ?? false;

                if (currentConnected !== data.connected) {
                    const newStatus: app.LCUStatus = {
                        connected: data.connected,
                        error: data.connected ? undefined : 'League client connection refused',
                    };
                    handleStatusResult(newStatus);
                }
            }, DEBOUNCE_MS);
        });

        return () => {
            if (debounceTimer) {
                clearTimeout(debounceTimer);
            }
            unsubscribe();
        };
    }, [handleStatusResult]);

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
