import { UserCard } from '../components';
import { lcu } from '../../wailsjs/go/models';

interface HomePageProps {
    summoner: lcu.CurrentSummoner | null;
}

const QUICK_STATS = [
    { icon: '🎮', label: 'Ready to Play' },
    { icon: '📈', label: 'Stats Available' },
] as const;

export function HomePage({ summoner }: HomePageProps) {
    return (
        <div className="tab-content">
            <UserCard summoner={summoner} />
            <QuickStats stats={QUICK_STATS} />
        </div>
    );
}

function QuickStats({ stats }: { stats: readonly { icon: string; label: string }[] }) {
    return (
        <div className="quick-stats">
            {stats.map((stat, index) => (
                <StatCard key={index} icon={stat.icon} label={stat.label} />
            ))}
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

