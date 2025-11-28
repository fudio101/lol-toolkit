import { UserCard } from '../components';
import { lcu } from '../../wailsjs/go/models';

interface ProfileTabProps {
    summoner: lcu.CurrentSummoner | null;
    onRefresh: () => void;
}

export function ProfileTab({ summoner, onRefresh }: ProfileTabProps) {
    return (
        <div className="tab-content">
            <UserCard summoner={summoner} onRefresh={onRefresh} />
        </div>
    );
}

