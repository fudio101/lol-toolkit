import { UserCard } from '../components';
import { lcu } from '../../wailsjs/go/models';

interface ProfilePageProps {
    summoner: lcu.CurrentSummoner | null;
    onRefresh: () => void;
}

export function ProfilePage({ summoner, onRefresh }: ProfilePageProps) {
    return (
        <div className="tab-content">
            <UserCard summoner={summoner} onRefresh={onRefresh} />
        </div>
    );
}

