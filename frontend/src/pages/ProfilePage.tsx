import { UserCard } from '../components';
import { lcu } from '../../wailsjs/go/models';

interface ProfilePageProps {
    summoner: lcu.CurrentSummoner | null;
}

export function ProfilePage({ summoner }: ProfilePageProps) {
    return (
        <div className="tab-content">
            <UserCard summoner={summoner} />
        </div>
    );
}

