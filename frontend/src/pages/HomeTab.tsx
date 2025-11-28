import { UserCard } from '../components';
import { lcu } from '../../wailsjs/go/models';

interface HomeTabProps {
    summoner: lcu.CurrentSummoner | null;
    onRefresh: () => void;
}

export function HomeTab({ summoner, onRefresh }: HomeTabProps) {
    return (
        <div className="tab-content">
            <UserCard summoner={summoner} onRefresh={onRefresh} />
            <div className="quick-stats">
                <div className="stat-card">
                    <span className="stat-icon">🎮</span>
                    <span className="stat-label">Ready to Play</span>
                </div>
                <div className="stat-card">
                    <span className="stat-icon">📈</span>
                    <span className="stat-label">Stats Available</span>
                </div>
            </div>
        </div>
    );
}

