import { useState } from 'react';
import { useConfig, useLCU, useTheme, useApiLog } from '../contexts';
import { useWindowSize } from '../hooks';
import { lcu } from '../../wailsjs/go/models';

interface DebugTabProps {
    summoner: lcu.CurrentSummoner | null;
}

export function DebugTab({ summoner }: DebugTabProps) {
    const { config } = useConfig();
    const { status, isPolling } = useLCU();
    const { theme, resolvedTheme } = useTheme();
    const { logs, clearLogs } = useApiLog();
    const [expandedLog, setExpandedLog] = useState<string | null>(null);
    const [copiedId, setCopiedId] = useState<string | null>(null);
    const windowSize = useWindowSize();

    const debugInfo = {
        'App Version': '1.0.0',
        'Window Size': `${windowSize.width} × ${windowSize.height}`,
        'Theme': `${theme} (resolved: ${resolvedTheme})`,
        'LCU Connected': status?.connected ? 'Yes' : 'No',
        'LCU Port': status?.port || 'N/A',
        'LCU Polling': isPolling ? 'Yes' : 'No',
        'Region': config?.region || 'N/A',
        'API Key': config?.riot_api_key ? `${config.riot_api_key.slice(0, 12)}...` : 'Not set',
        'Summoner PUUID': summoner?.puuid ? `${summoner.puuid.slice(0, 20)}...` : 'N/A',
        'Summoner ID': summoner?.summonerId || 'N/A',
        'Account ID': summoner?.accountId || 'N/A',
    };

    const formatTime = (date: Date) => {
        return date.toLocaleTimeString('en-US', { 
            hour12: false, 
            hour: '2-digit', 
            minute: '2-digit', 
            second: '2-digit',
            fractionalSecondDigits: 3 
        });
    };

    const copyToClipboard = async (text: string, id: string) => {
        try {
            await navigator.clipboard.writeText(text);
            setCopiedId(id);
            setTimeout(() => setCopiedId(null), 2000);
        } catch (err) {
            console.error('Failed to copy:', err);
        }
    };

    const getEndpointPath = (endpoint: string) => {
        const paths: Record<string, string> = {
            'GetLCUStatus': '/lol-gameflow/v1/gameflow-phase',
            'GetCurrentSummoner': '/lol-summoner/v1/current-summoner',
        };
        return paths[endpoint] || `/${endpoint}`;
    };

    const getBase64Auth = () => {
        if (status?.authToken) {
            return btoa(`riot:${status.authToken}`);
        }
        return 'BASE64_AUTH_TOKEN';
    };

    const generateCurl = (log: typeof logs[0]) => {
        const port = status?.port || 'PORT';
        const auth = getBase64Auth();
        const path = getEndpointPath(log.endpoint);
        
        if (log.type === 'lcu') {
            return `curl -X ${log.method} "https://127.0.0.1:${port}${path}" -H "Authorization: Basic ${auth}" -H "Accept: application/json" --insecure`;
        }
        const apiKey = config?.riot_api_key || 'YOUR_API_KEY';
        const region = config?.region || 'na1';
        return `curl -X ${log.method} "https://${region}.api.riotgames.com${path}" -H "X-Riot-Token: ${apiKey}" -H "Accept: application/json"`;
    };

    const generatePowerShell = (log: typeof logs[0]) => {
        const port = status?.port || 'PORT';
        const auth = getBase64Auth();
        const path = getEndpointPath(log.endpoint);
        
        if (log.type === 'lcu') {
            return `Invoke-RestMethod -Uri "https://127.0.0.1:${port}${path}" -Method ${log.method} -Headers @{"Authorization"="Basic ${auth}";"Accept"="application/json"} -SkipCertificateCheck`;
        }
        const apiKey = config?.riot_api_key || 'YOUR_API_KEY';
        const region = config?.region || 'na1';
        return `Invoke-RestMethod -Uri "https://${region}.api.riotgames.com${path}" -Method ${log.method} -Headers @{"X-Riot-Token"="${apiKey}";"Accept"="application/json"}`;
    };

    return (
        <div className="tab-content">
            <div className="debug-card">
                <div className="debug-card-header">
                    <h3>📡 API Call Log</h3>
                    <div className="debug-actions">
                        <span className="log-count">{logs.length} calls</span>
                        <button className="btn-small btn-danger" onClick={clearLogs}>Clear</button>
                    </div>
                </div>
                
                <div className="api-log-list">
                    {logs.length === 0 ? (
                        <div className="api-log-empty">No API calls logged yet.</div>
                    ) : (
                        logs.map((log) => (
                            <div 
                                key={log.id} 
                                className={`api-log-item ${log.status}`}
                            >
                                <div 
                                    className="api-log-header"
                                    onClick={() => setExpandedLog(expandedLog === log.id ? null : log.id)}
                                >
                                    <span className={`api-log-status ${log.status}`}>
                                        {log.status === 'pending' ? '⏳' : log.status === 'success' ? '✅' : '❌'}
                                    </span>
                                    <span className="api-log-type">{log.type.toUpperCase()}</span>
                                    <span className="api-log-method">{log.method}</span>
                                    <span className="api-log-endpoint">{log.endpoint}</span>
                                    <span className="api-log-time">{formatTime(log.timestamp)}</span>
                                    {log.duration && (
                                        <span className="api-log-duration">{log.duration}ms</span>
                                    )}
                                    <span className="api-log-expand">{expandedLog === log.id ? '▼' : '▶'}</span>
                                </div>
                                {expandedLog === log.id && (
                                    <div className="api-log-details">
                                        <div className="api-log-copy-buttons">
                                            <button 
                                                className={`btn-copy ${copiedId === `${log.id}-curl` ? 'copied' : ''}`}
                                                onClick={(e) => { e.stopPropagation(); copyToClipboard(generateCurl(log), `${log.id}-curl`); }}
                                            >
                                                {copiedId === `${log.id}-curl` ? '✓ Copied' : '📋 cURL'}
                                            </button>
                                            <button 
                                                className={`btn-copy ${copiedId === `${log.id}-ps` ? 'copied' : ''}`}
                                                onClick={(e) => { e.stopPropagation(); copyToClipboard(generatePowerShell(log), `${log.id}-ps`); }}
                                            >
                                                {copiedId === `${log.id}-ps` ? '✓ Copied' : '📋 PowerShell'}
                                            </button>
                                            {log.response && (
                                                <button 
                                                    className={`btn-copy ${copiedId === `${log.id}-response` ? 'copied' : ''}`}
                                                    onClick={(e) => { e.stopPropagation(); copyToClipboard(JSON.stringify(log.response, null, 2), `${log.id}-response`); }}
                                                >
                                                    {copiedId === `${log.id}-response` ? '✓ Copied' : '📋 Response'}
                                                </button>
                                            )}
                                        </div>
                                        {log.error ? (
                                            <pre className="api-log-error">{log.error}</pre>
                                        ) : log.response ? (
                                            <pre className="api-log-response">{JSON.stringify(log.response, null, 2)}</pre>
                                        ) : null}
                                    </div>
                                )}
                            </div>
                        ))
                    )}
                </div>
            </div>

            <div className="debug-card">
                <h3>🐛 Debug Information</h3>
                <div className="debug-grid">
                    {Object.entries(debugInfo).map(([key, value]) => (
                        <div key={key} className="debug-row">
                            <span className="debug-key">{key}</span>
                            <span className="debug-value">{value}</span>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}

