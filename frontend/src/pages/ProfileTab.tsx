import { useState } from 'react';
import { UserCard } from '../components';
import { lcu } from '../../wailsjs/go/models';

interface ProfileTabProps {
    summoner: lcu.CurrentSummoner | null;
}

export function ProfileTab({ summoner }: ProfileTabProps) {
    const [copiedField, setCopiedField] = useState<string | null>(null);

    const copyToClipboard = async (text: string, field: string) => {
        try {
            await navigator.clipboard.writeText(text);
            setCopiedField(field);
            setTimeout(() => setCopiedField(null), 2000);
        } catch (err) {
            console.error('Failed to copy:', err);
        }
    };

    if (!summoner) {
        return (
            <div className="tab-content">
                <div className="profile-loading">
                    <div className="spinner" />
                    <p>Loading profile data...</p>
                </div>
            </div>
        );
    }

    const totalXpForLevel = (summoner.xpSinceLastLevel || 0) + (summoner.xpUntilNextLevel || 0);
    const reroll = summoner.rerollPoints;

    return (
        <div className="tab-content">
            <UserCard summoner={summoner} />

            {/* Two Column Layout */}
            <div className="profile-two-col">
                {/* Left Column - Stats */}
                <div className="profile-column">
                    {/* Level Progress */}
                    <div className="profile-card">
                        <h3>⬆️ Level Progress</h3>
                        <div className="level-display">
                            <span className="level-number">{summoner.summonerLevel}</span>
                            <div className="level-details">
                                <div className="level-xp-bar">
                                    <div 
                                        className="level-xp-fill" 
                                        style={{ width: `${summoner.percentCompleteForNextLevel}%` }}
                                    />
                                </div>
                                <div className="level-xp-text">
                                    <span>{(summoner.xpSinceLastLevel || 0).toLocaleString()} / {totalXpForLevel.toLocaleString()} XP</span>
                                    <span>{summoner.percentCompleteForNextLevel}%</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* ARAM Rerolls */}
                    <div className="profile-card">
                        <h3>🎲 ARAM Rerolls</h3>
                        <div className="reroll-display">
                            <div className="reroll-rolls">
                                <span className="reroll-current">{reroll?.numberOfRolls || 0}</span>
                                <span className="reroll-max">/ {reroll?.maxRolls || 0}</span>
                            </div>
                            <div className="reroll-details">
                                <div className="level-xp-bar">
                                    <div 
                                        className="level-xp-fill" 
                                        style={{ 
                                            width: `${reroll?.pointsToReroll 
                                                ? Math.min(100, ((reroll.currentPoints || 0) / reroll.pointsToReroll) * 100) 
                                                : 0}%` 
                                        }}
                                    />
                                </div>
                                <div className="level-xp-text">
                                    <span>{reroll?.currentPoints || 0} / {reroll?.pointsToReroll || 0} pts</span>
                                    <span>{reroll?.pointsCostToRoll || 0} per roll</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Right Column - Account Info */}
                <div className="profile-column">
                    <div className="profile-card profile-card-full">
                        <h3>🔑 Account Details</h3>
                        <div className="account-list">
                            <CopyRow 
                                label="Riot ID" 
                                value={`${summoner.gameName}#${summoner.tagLine}`}
                                onCopy={copyToClipboard}
                                copied={copiedField === 'riotId'}
                                field="riotId"
                            />
                            <CopyRow 
                                label="PUUID" 
                                value={summoner.puuid || ''}
                                displayValue={summoner.puuid ? `${summoner.puuid.slice(0, 8)}...${summoner.puuid.slice(-4)}` : 'N/A'}
                                onCopy={copyToClipboard}
                                copied={copiedField === 'puuid'}
                                field="puuid"
                                mono
                            />
                            <CopyRow 
                                label="Summoner ID" 
                                value={String(summoner.summonerId)}
                                onCopy={copyToClipboard}
                                copied={copiedField === 'summonerId'}
                                field="summonerId"
                            />
                            <CopyRow 
                                label="Account ID" 
                                value={String(summoner.accountId)}
                                onCopy={copyToClipboard}
                                copied={copiedField === 'accountId'}
                                field="accountId"
                            />
                            <CopyRow 
                                label="Icon ID" 
                                value={String(summoner.profileIconId)}
                                onCopy={copyToClipboard}
                                copied={copiedField === 'iconId'}
                                field="iconId"
                            />
                            <div className="account-row-simple">
                                <span className="account-label">Name Change</span>
                                <span className={`status-badge ${summoner.nameChangeFlag ? 'status-yes' : 'status-no'}`}>
                                    {summoner.nameChangeFlag ? 'Available' : 'No'}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

interface CopyRowProps {
    label: string;
    value: string;
    displayValue?: string;
    onCopy: (text: string, field: string) => void;
    copied: boolean;
    field: string;
    mono?: boolean;
}

function CopyRow({ label, value, displayValue, onCopy, copied, field, mono }: CopyRowProps) {
    return (
        <div className="account-row-copy">
            <span className="account-label">{label}</span>
            <div className="account-value-wrapper">
                <span className={`account-value ${mono ? 'mono' : ''}`}>
                    {displayValue || value}
                </span>
                <button 
                    className={`copy-btn ${copied ? 'copied' : ''}`}
                    onClick={() => onCopy(value, field)}
                    title="Copy"
                >
                    {copied ? '✓' : '📋'}
                </button>
            </div>
        </div>
    );
}
