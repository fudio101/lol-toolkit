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
    const { logs } = useApiLog();
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

    const isLCUConnected = (log: any): boolean => {
        return log.type === 'lcu' && !!status?.port;
    };

    const buildLCUUrl = (endpoint: string): string => {
        return `https://127.0.0.1:${status?.port}${endpoint}`;
    };

    const formatCurlHeaders = (headers: Record<string, string>): string => {
        return Object.entries(headers)
            .map(([key, value]) => ` \\\n  -H "${key}: ${value}"`)
            .join('');
    };

    const formatPowerShellHeaders = (headers: Record<string, string>): string => {
        const headerEntries = Object.entries(headers)
            .map(([key, value]) => `\n    "${key}" = "${value}"`)
            .join(';');
        return ` \`\n  -Headers @{${headerEntries}\n  }`;
    };

    const generateCurlCommand = (log: any): string => {
        if (!isLCUConnected(log)) {
            return 'N/A - LCU not connected';
        }

        const url = buildLCUUrl(log.endpoint);
        let command = `curl -k -X ${log.method} "${url}"`;

        if (log.headers) {
            command += formatCurlHeaders(log.headers);
        }

        return command;
    };

    const generatePowerShellCommand = (log: any): string => {
        if (!isLCUConnected(log)) {
            return 'N/A - LCU not connected';
        }

        const url = buildLCUUrl(log.endpoint);
        let command = `Invoke-WebRequest -Uri "${url}" \`\n  -Method ${log.method} \`\n  -SkipCertificateCheck`;

        if (log.headers && Object.keys(log.headers).length > 0) {
            command += formatPowerShellHeaders(log.headers);
        }

        return command;
    };


    return (
        <div className="tab-content">
            <div className="debug-card">
                <div className="debug-card-header">
                    <h3>🛰 API Calls</h3>
                    <div className="debug-actions">
                        <span className="log-count">{logs.length} calls</span>
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
                                    {log.duration !== undefined && (
                                        <span className="api-log-duration">{log.duration}ms</span>
                                    )}
                                    <span className="api-log-expand">{expandedLog === log.id ? '▼' : '▶'}</span>
                                </div>
                                {expandedLog === log.id && (
                                    <div className="api-log-details">
                                        {log.type === 'lcu' && log.endpoint !== 'GetLCUStatus' && (
                                            <div className="api-log-section" style={{ marginBottom: '12px' }}>
                                                <div style={{ display: 'flex', gap: '8px' }}>
                                                    <button 
                                                        className={`api-copy-btn ${copiedId === `${log.id}-curl` ? 'copied' : ''}`}
                                                        onClick={(e) => { 
                                                            e.stopPropagation(); 
                                                            copyToClipboard(generateCurlCommand(log), `${log.id}-curl`); 
                                                        }}
                                                        title="Copy cURL command"
                                                    >
                                                        {copiedId === `${log.id}-curl` ? '✓' : '📋'} cURL
                                                    </button>
                                                    <button 
                                                        className={`api-copy-btn ${copiedId === `${log.id}-pwsh` ? 'copied' : ''}`}
                                                        onClick={(e) => { 
                                                            e.stopPropagation(); 
                                                            copyToClipboard(generatePowerShellCommand(log), `${log.id}-pwsh`); 
                                                        }}
                                                        title="Copy PowerShell command"
                                                    >
                                                        {copiedId === `${log.id}-pwsh` ? '✓' : '📋'} PowerShell
                                                    </button>
                                                </div>
                                            </div>
                                        )}
                                        {log.error && (
                                            <div className="api-log-section">
                                                <div className="api-log-code-block">
                                                    <button 
                                                        className={`btn-copy-icon ${copiedId === `${log.id}-error` ? 'copied' : ''}`}
                                                        onClick={(e) => { e.stopPropagation(); copyToClipboard(log.error!, `${log.id}-error`); }}
                                                        title="Copy error"
                                                    >
                                                        {copiedId === `${log.id}-error` ? '✓' : '📋'}
                                                    </button>
                                                    <pre className="api-log-error">{log.error}</pre>
                                                </div>
                                            </div>
                                        )}
                                        {!log.error && log.response && (
                                            <div className="api-log-section">
                                                <div className="api-log-section-title">Response</div>
                                                <div className="api-log-code-block">
                                                    <button 
                                                        className={`btn-copy-icon ${copiedId === `${log.id}-response` ? 'copied' : ''}`}
                                                        onClick={(e) => { 
                                                            e.stopPropagation(); 
                                                            const responseText = typeof log.response === 'string'
                                                                ? log.response
                                                                : JSON.stringify(log.response, null, 2);
                                                            copyToClipboard(responseText, `${log.id}-response`); 
                                                        }}
                                                        title="Copy response"
                                                    >
                                                        {copiedId === `${log.id}-response` ? '✓' : '📋'}
                                                    </button>
                                                    <pre className="api-log-response">
                                                        {typeof log.response === 'string'
                                                            ? log.response
                                                            : JSON.stringify(log.response, null, 2)}
                                                    </pre>
                                                </div>
                                            </div>
                                        )}
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

