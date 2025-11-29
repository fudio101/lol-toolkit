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
    const { status, summoner, loading: lcuLoading } = useLCU();

    const isLoading = configLoading || lcuLoading;
    const isConnected = status?.connected ?? false;

    if (isLoading) {
        return <LoadingScreen />;
    }

    return (
        <div className="app-layout">
            <Sidebar 
                activeTab={activeTab} 
                onTabChange={handleTabChange(setActiveTab)} 
            />
            
            <main className="main-content">
                <Header title={getTabTitle(activeTab)} />
                
                <div className="content-area">
                    <TabContent 
                        tab={activeTab} 
                        summoner={summoner} 
                        isConnected={isConnected}
                        isConfigured={isConfigured}
                    />
                </div>
            </main>
        </div>
    );
}

function LoadingScreen() {
    return (
        <div className="app-layout">
            <div className="loading-screen">
                <div className="spinner" />
                <p>Loading...</p>
            </div>
        </div>
    );
}

function handleTabChange(setActiveTab: (tab: TabId) => void) {
    return (tab: string) => setActiveTab(tab as TabId);
}

function getTabTitle(tab: TabId): string {
    return TAB_TITLES[tab] || 'Dashboard';
}

function NotConnectedMessage() {
    return (
        <div className="message-card">
            <div className="message-icon">🔌</div>
            <h2>League Client Not Running</h2>
            <p>Please open the League of Legends client to use this app.</p>
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
    isConnected: boolean;
    isConfigured: boolean;
}

// Tabs that don't require LCU connection
const CONNECTION_FREE_TABS: TabId[] = ['settings', 'debug'];

function TabContent({ tab, summoner, isConnected, isConfigured }: TabContentProps) {
    if (isConnectionFreeTab(tab)) {
        return renderConnectionFreeTab(tab, summoner);
    }

    if (!isConnected) {
        return <NotConnectedMessage />;
    }

    if (!isConfigured) {
        return <ApiKeyRequiredMessage />;
    }

    return renderConnectedTab(tab, summoner);
}

function isConnectionFreeTab(tab: TabId): boolean {
    return CONNECTION_FREE_TABS.includes(tab);
}

function renderConnectionFreeTab(tab: TabId, summoner: any) {
    switch (tab) {
        case 'settings':
            return <SettingsTab />;
        case 'debug':
            return <DebugTab summoner={summoner} />;
        default:
            return null;
    }
}

function renderConnectedTab(tab: TabId, summoner: any) {
    switch (tab) {
        case 'home':
            return <HomeTab summoner={summoner} />;
        case 'profile':
            return <ProfileTab summoner={summoner} />;
        case 'champions':
            return <ChampionsTab />;
        case 'matches':
            return <MatchesTab />;
        default:
            return null;
    }
}

export default App;
