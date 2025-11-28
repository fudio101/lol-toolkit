interface StatusBarProps {
    isConnected: boolean;
    isConfigured: boolean;
    isPolling: boolean;
}

export function StatusBar({ isConnected, isConfigured, isPolling }: StatusBarProps) {
    return (
        <footer className="status-bar">
            <div className="status-item">
                <span className={`status-dot ${isConnected ? 'online' : 'offline'} ${isPolling ? 'polling' : ''}`} />
                LCU: {isConnected ? 'Connected' : 'Disconnected'}
            </div>
            <div className="status-item">
                <span className={`status-dot ${isConfigured ? 'online' : 'offline'}`} />
                API: {isConfigured ? 'Configured' : 'Not set'}
            </div>
        </footer>
    );
}

