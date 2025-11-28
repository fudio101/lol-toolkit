import { useState } from 'react';
import './App.css';
import { useConfig, useLCU } from './contexts';
import { Sidebar, Header } from './components';
import { HomeTab, ProfileTab, ChampionsTab, MatchesTab, SettingsTab, DebugTab } from './pages';

export type TabId = 'home' | 'profile' | 'champions' | 'matches' | 'settings' | 'debug';

const TAB_TITLES: Record<TabId, string> = {
    home: 'Dashboard',
    profile: 'Profile',
    champions: 'Champions',
    matches: 'Match History',
    settings: 'Settings',
    debug: 'Debug',
};

function App() {
    const [activeTab, setActiveTab] = useState<TabId>('home');
    const { isConfigured, loading: configLoading } = useConfig();
    const { status, summoner, loading: lcuLoading, refresh } = useLCU();

    const isConnected = status?.connected ?? false;

    if (configLoading || lcuLoading) {
        return (
            <div className="app-layout">
                <div className="loading-screen">
                    <div className="spinner" />
                    <p>Loading...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="app-layout">
            <Sidebar activeTab={activeTab} onTabChange={(tab) => setActiveTab(tab as TabId)} />
            
            <main className="main-content">
                <Header title={TAB_TITLES[activeTab] || 'Dashboard'} />
                
                <div className="content-area">
                    <TabContent 
                        tab={activeTab} 
                        summoner={summoner} 
                        onRefresh={refresh}
                        isConnected={isConnected}
                        isConfigured={isConfigured}
                    />
                </div>
            </main>
        </div>
    );
}

function NotConnectedMessage({ onRefresh }: { onRefresh: () => void }) {
    return (
        <div className="message-card">
            <div className="message-icon">🔌</div>
            <h2>League Client Not Running</h2>
            <p>Please open the League of Legends client to use this app.</p>
            <button className="btn-primary" onClick={onRefresh}>
                Check Connection
            </button>
        </div>
    );
}

function ApiKeyRequiredMessage() {
    return (
        <div className="message-card">
            <div className="message-icon">⚠️</div>
            <h2>API Key Required</h2>
            <p>Please add your Riot API key to use all features.</p>
            <code className="code-block">internal/config/config.json</code>
            <a href="https://developer.riotgames.com/" target="_blank" className="btn-link">
                Get API Key →
            </a>
        </div>
    );
}

interface TabContentProps {
    tab: TabId;
    summoner: any;
    onRefresh: () => void;
    isConnected: boolean;
    isConfigured: boolean;
}

// Tabs that don't require LCU connection
const CONNECTION_FREE_TABS: TabId[] = ['settings', 'debug'];

function TabContent({ tab, summoner, onRefresh, isConnected, isConfigured }: TabContentProps) {
    // Settings and Debug are always accessible
    if (CONNECTION_FREE_TABS.includes(tab)) {
        switch (tab) {
            case 'settings':
                return <SettingsTab />;
            case 'debug':
                return <DebugTab summoner={summoner} />;
        }
    }

    // Other tabs require connection
    if (!isConnected) {
        return <NotConnectedMessage onRefresh={onRefresh} />;
    }

    if (!isConfigured) {
        return <ApiKeyRequiredMessage />;
    }

    switch (tab) {
        case 'home':
            return <HomeTab summoner={summoner} onRefresh={onRefresh} />;
        case 'profile':
            return <ProfileTab summoner={summoner} onRefresh={onRefresh} />;
        case 'champions':
            return <ChampionsTab />;
        case 'matches':
            return <MatchesTab />;
        default:
            return null;
    }
}

export default App;
