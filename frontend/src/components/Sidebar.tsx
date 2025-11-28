import { useLCU, useSettings } from '../contexts';

interface NavItem {
    id: string;
    icon: string;
    label: string;
    hidden?: boolean;
}

interface SidebarProps {
    activeTab: string;
    onTabChange: (tab: string) => void;
}

export function Sidebar({ activeTab, onTabChange }: SidebarProps) {
    const { status, isPolling } = useLCU();
    const { settings } = useSettings();

    const navItems: NavItem[] = [
        { id: 'home', icon: '🏠', label: 'Home' },
        { id: 'profile', icon: '👤', label: 'Profile' },
        { id: 'champions', icon: '⚔️', label: 'Champions' },
        { id: 'matches', icon: '📊', label: 'Matches' },
        { id: 'settings', icon: '⚙️', label: 'Settings' },
        { id: 'debug', icon: '🐛', label: 'Debug', hidden: !settings.showDebug },
    ];

    const visibleItems = navItems.filter(item => !item.hidden);

    return (
        <aside className="sidebar">
            <div className="sidebar-header">
                <div className="logo">
                    <span className="logo-icon">⚡</span>
                    <span className="logo-text">LoL Toolkit</span>
                </div>
            </div>

            <nav className="sidebar-nav">
                {visibleItems.map((item) => (
                    <button
                        key={item.id}
                        className={`nav-item ${activeTab === item.id ? 'active' : ''}`}
                        onClick={() => onTabChange(item.id)}
                    >
                        <span className="nav-icon">{item.icon}</span>
                        <span className="nav-label">{item.label}</span>
                    </button>
                ))}
            </nav>

            <div className="sidebar-footer">
                <div className="connection-status">
                    <span className={`status-indicator ${status?.connected ? 'connected' : 'disconnected'} ${isPolling ? 'polling' : ''}`} />
                    <span className="status-text">
                        {status?.connected ? 'Connected' : 'Disconnected'}
                    </span>
                </div>
            </div>
        </aside>
    );
}
