import { useEffect, useRef, useCallback, useState } from 'react';

export interface PollingTask<T> {
    /** Unique task identifier */
    id: string;
    /** Function that performs the polling */
    execute: () => Promise<T>;
    /** Get interval based on current result (dynamic intervals) */
    getInterval: (result: T | null, error: Error | null) => number;
    /** Called when polling completes */
    onResult?: (result: T) => void;
    /** Called on error */
    onError?: (error: Error) => void;
    /** Whether this task is enabled (default: true) */
    enabled?: boolean;
}

/**
 * Hook for managing a single polling task with dynamic intervals.
 */
export function usePolling<T>(task: PollingTask<T>): {
    isPolling: boolean;
    refresh: () => Promise<void>;
} {
    const [isPolling, setIsPolling] = useState(false);
    const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const isExecuting = useRef(false);
    const taskRef = useRef(task);
    
    // Keep task ref updated
    taskRef.current = task;

    const clearTimer = useCallback(() => {
        if (timerRef.current !== null) {
            clearTimeout(timerRef.current);
            timerRef.current = null;
        }
    }, []);

    const poll = useCallback(async (): Promise<void> => {
        // Guard against concurrent execution
        if (isExecuting.current) return;
        
        const currentTask = taskRef.current;
        if (currentTask.enabled === false) return;
        
        isExecuting.current = true;
        setIsPolling(true);

        let result: T | null = null;
        let error: Error | null = null;

        try {
            result = await currentTask.execute();
            currentTask.onResult?.(result);
        } catch (err) {
            error = err instanceof Error ? err : new Error(String(err));
            currentTask.onError?.(error);
        }

        isExecuting.current = false;
        setIsPolling(false);

        // Schedule next poll if still enabled
        if (taskRef.current.enabled !== false) {
            const interval = taskRef.current.getInterval(result, error);
            clearTimer();
            timerRef.current = setTimeout(poll, interval);
        }
    }, [clearTimer]);

    const refresh = useCallback(async () => {
        clearTimer();
        await poll();
    }, [clearTimer, poll]);

    // Lifecycle: start/stop polling based on enabled state
    useEffect(() => {
        const enabled = task.enabled !== false;
        
        if (enabled && !isExecuting.current && timerRef.current === null) {
            // Start polling
            poll();
        } else if (!enabled) {
            // Stop polling
            clearTimer();
        }

        return () => {
            clearTimer();
        };
    }, [task.enabled, poll, clearTimer]);

    return { isPolling, refresh };
}
