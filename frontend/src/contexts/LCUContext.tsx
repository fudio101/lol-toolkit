import { createContext, useContext, useState, useCallback, useRef, ReactNode, useEffect } from 'react';
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
    
    // Refs for tracking state and function references
    const statusRef = useRef(status);
    // Keep statusRef updated
    useEffect(() => {
        statusRef.current = status;
    }, [status]);
    
    const prevConnectedRef = useRef<boolean | null>(null);
    const summonerRefreshRef = useRef<(() => Promise<void>) | null>(null);

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

    // ============================================================================
    // Summoner Functions (defined early so they can be used in handleStatusResult)
    // ============================================================================

    const fetchSummoner = useCallback(async (): Promise<lcu.CurrentSummoner | null> => {
        // Only fetch if connected
        if (!statusRef.current?.connected) {
            return null;
        }

        try {
            return await GetCurrentSummoner();
        } catch (err) {
            // Errors are logged by the backend, just return null
            return null;
        }
    }, []);

    const handleSummonerResult = useCallback((result: lcu.CurrentSummoner | null) => {
        if (result) {
            setSummoner(result);
        }
    }, []);

    const handleStatusResult = useCallback((result: app.LCUStatus) => {
        // Update statusRef immediately so fetchSummoner can use it
        statusRef.current = result;
        
        const wasConnected = prevConnectedRef.current;
        const isNowConnected = result.connected === true;
        
        // Detect transition: (null/false) → true
        const wasNotConnected = wasConnected === null || wasConnected === false;
        const isNewlyConnected = wasNotConnected && isNowConnected;
        
        // Update ref before calling setStatus
        prevConnectedRef.current = isNowConnected;
        
        setStatus(result);
        setError(null);
        setLoading(false);
        
        // Clear summoner if disconnected
        if (!result.connected) {
            setSummoner(null);
        } else if (isNewlyConnected) {
            // Immediately fetch summoner when connection is established
            // refreshSummoner will execute immediately and reset the polling timer
            if (summonerRefreshRef.current) {
                summonerRefreshRef.current();
            }
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

    // Backup handler: Watch for connection state changes in case handleStatusResult misses it
    useEffect(() => {
        if (status === null) {
            return;
        }

        const isNowConnected = status.connected === true;
        const wasConnected = prevConnectedRef.current;
        const wasNotConnected = wasConnected === null || wasConnected === false;
        const isNewlyConnected = wasNotConnected && isNowConnected;

        if (isNewlyConnected && summonerRefreshRef.current) {
            summonerRefreshRef.current();
        }
    }, [status]);

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
