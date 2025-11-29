import { lcu } from "../../wailsjs/go/models";

interface UserCardProps {
    summoner: lcu.CurrentSummoner | null;
}

const DDRAGON_VERSION = "14.23.1";

export function UserCard({ summoner }: UserCardProps) {
    if (!summoner) {
        return <LoadingSummonerCard />;
    }

    return (
        <div className="user-card">
            <UserAvatar iconId={summoner.profileIconId} />
            <UserInfo summoner={summoner} />
        </div>
    );
}

function LoadingSummonerCard() {
    return (
        <div className="user-card">
            <div className="user-loading">Loading summoner data...</div>
        </div>
    );
}

function UserAvatar({ iconId }: { iconId: number }) {
    const profileIconUrl = getProfileIconUrl(iconId);
    
    return (
        <div className="user-avatar">
            <img src={profileIconUrl} alt="Profile Icon" />
            <span className="status-dot" />
        </div>
    );
}

function UserInfo({ summoner }: { summoner: lcu.CurrentSummoner }) {
    return (
        <div className="user-info">
            <UserName gameName={summoner.gameName} tagLine={summoner.tagLine} />
            <UserLevel level={summoner.summonerLevel} />
            <UserXP percentComplete={summoner.percentCompleteForNextLevel} />
        </div>
    );
}

function UserName({ gameName, tagLine }: { gameName: string; tagLine: string }) {
    return (
        <div className="user-name">
            {gameName}
            <span className="user-tag">#{tagLine}</span>
        </div>
    );
}

function UserLevel({ level }: { level: number }) {
    return <div className="user-level">Level {level}</div>;
}

function UserXP({ percentComplete }: { percentComplete: number }) {
    return (
        <div className="user-xp">
            <div className="xp-bar">
                <div
                    className="xp-fill"
                    style={{ width: `${percentComplete}%` }}
                />
            </div>
            <span className="xp-text">{percentComplete}% to next level</span>
        </div>
    );
}

function getProfileIconUrl(iconId: number): string {
    return `https://ddragon.leagueoflegends.com/cdn/${DDRAGON_VERSION}/img/profileicon/${iconId}.png`;
}
