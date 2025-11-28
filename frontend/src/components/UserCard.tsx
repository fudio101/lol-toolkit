import { lcu } from "../../wailsjs/go/models";

interface UserCardProps {
    summoner: lcu.CurrentSummoner | null;
    onRefresh: () => void;
}

const DDRAGON_VERSION = "14.23.1";

export function UserCard({ summoner, onRefresh }: UserCardProps) {
    if (!summoner) {
        return (
            <div className="user-card">
                <div className="user-loading">Loading summoner...</div>
            </div>
        );
    }

    const profileIconUrl = `https://ddragon.leagueoflegends.com/cdn/${DDRAGON_VERSION}/img/profileicon/${summoner.profileIconId}.png`;

    return (
        <div className="user-card">
            <div className="user-avatar">
                <img src={profileIconUrl} alt="Profile Icon" />
                <span className="status-dot online" title="Online" />
            </div>
            <div className="user-info">
                <div className="user-name">
                    {summoner.gameName}
                    <span className="user-tag">#{summoner.tagLine}</span>
                </div>
                <div className="user-level">Level {summoner.summonerLevel}</div>
                <div className="user-xp">
                    <div className="xp-bar">
                        <div
                            className="xp-fill"
                            style={{ width: `${summoner.percentCompleteForNextLevel}%` }}
                        />
                    </div>
                    <span className="xp-text">{summoner.percentCompleteForNextLevel}%</span>
                </div>
            </div>
            <button className="refresh-small" onClick={onRefresh} title="Refresh">
                🔄
            </button>
        </div>
    );
}

