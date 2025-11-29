import { UserCard } from '../components';
import { lcu } from '../../wailsjs/go/models';

interface HomeTabProps {
    summoner: lcu.CurrentSummoner | null;
}

export function HomeTab({ summoner }: HomeTabProps) {
    return (
        <div className="tab-content">
            <UserCard summoner={summoner} />
            <QuickStats />
        </div>
    );
}

function QuickStats() {
    return (
        <div className="quick-stats">
            <StatCard icon="🎮" label="Ready to Play" />
            <StatCard icon="📈" label="Stats Available" />
        </div>
    );
}

function StatCard({ icon, label }: { icon: string; label: string }) {
    return (
        <div className="stat-card">
            <span className="stat-icon">{icon}</span>
            <span className="stat-label">{label}</span>
        </div>
    );
}

