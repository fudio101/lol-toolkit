import './App.css';
import { useConfig, useLCU } from './contexts';
import { StatusBar, UserCard } from './components';

function App() {
    const { config, isConfigured, loading: configLoading } = useConfig();
    const { status, summoner, loading: lcuLoading, isPolling, refresh } = useLCU();

    const isConnected = status?.connected ?? false;

    // Loading state
    if (configLoading || lcuLoading) {
        return (
            <div className="app">
                <div className="loading-screen">Loading...</div>
            </div>
        );
    }

    // LCU not connected
    if (!isConnected) {
        return (
            <div className="app">
                <div className="container">
                    <h1>🔌 Connect to League Client</h1>
                    <p className="subtitle">
                        Please open the League of Legends client to use this app.
                    </p>
                    <button className="btn-primary" onClick={refresh}>
                        🔄 Check Connection
                    </button>
                </div>
                <StatusBar isConnected={false} isConfigured={isConfigured} isPolling={isPolling} />
            </div>
        );
    }

    // API not configured
    if (!isConfigured) {
        return (
            <div className="app">
                <div className="container">
                    <h1>⚠️ API Key Required</h1>
                    <p className="subtitle">
                        Please add your Riot API key to <code>internal/config/config.json</code> and rebuild the app.
                    </p>
                    <a href="https://developer.riotgames.com/" target="_blank" className="link">
                        Get API Key →
                    </a>
                </div>
                <StatusBar isConnected={isConnected} isConfigured={false} isPolling={isPolling} />
            </div>
        );
    }

    // Main app
    return (
        <div className="app">
            <header className="header">
                <h1>🎮 LoL Toolkit</h1>
                <span className="region-badge">{config?.region?.toUpperCase()}</span>
            </header>

            <UserCard summoner={summoner} onRefresh={refresh} />

            <div className="content">
                {/* Future features */}
            </div>

            <StatusBar isConnected={isConnected} isConfigured={isConfigured} isPolling={isPolling} />
        </div>
    );
}

export default App;
