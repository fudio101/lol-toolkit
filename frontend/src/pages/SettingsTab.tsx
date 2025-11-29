import { useEffect, useRef } from 'react';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { useSettings, useTheme, useLCU } from '../contexts';
import { StartAutoAccept, StopAutoAccept, IsAutoAcceptRunning } from '../../wailsjs/go/app/App';
import { app } from '../../wailsjs/go/models';

export function SettingsTab() {
    const { settings, updateSettings } = useSettings();
    const { theme, setTheme } = useTheme();
    const { status } = useLCU();
    const prevConnectedRef = useRef<boolean | null>(null);
    const wasDisabledByDisconnectRef = useRef<boolean>(false);
    const hasCheckedOnMountRef = useRef<boolean>(false);

    // Check service state on mount
    useEffect(() => {
        if (!hasCheckedOnMountRef.current && status?.connected) {
            hasCheckedOnMountRef.current = true;
            verifyServiceState();
        }
    }, [status?.connected, settings.autoAcceptEnabled, updateSettings]);

    // Listen for auto-accept stopped event from backend
    useEffect(() => {
        return EventsOn('auto-accept-stopped', () => {
            if (settings.autoAcceptEnabled) {
                updateSettings({ autoAcceptEnabled: false });
            }
        });
    }, [settings.autoAcceptEnabled, updateSettings]);

    // Handle connection state changes
    useEffect(() => {
        const wasConnected = prevConnectedRef.current;
        const isNowConnected = status?.connected ?? false;

        handleConnectionTransition(wasConnected, isNowConnected);
        prevConnectedRef.current = isNowConnected;
    }, [status?.connected, settings.autoAcceptEnabled, updateSettings]);

    const verifyServiceState = () => {
        IsAutoAcceptRunning()
            .then((isRunning: boolean) => {
                if (!isRunning && settings.autoAcceptEnabled) {
                    updateSettings({ autoAcceptEnabled: false });
                }
            })
            .catch((err: any) => {
                console.error('Failed to check auto-accept service status:', err);
            });
    };

    const handleConnectionTransition = (wasConnected: boolean | null, isNowConnected: boolean) => {
        // Disconnection
        if (wasConnected === true && !isNowConnected) {
            handleDisconnection();
            return;
        }

        // Reconnection
        if (wasConnected === false && isNowConnected) {
            handleReconnection();
            return;
        }

        // Not connected - ensure toggle is off
        if (!isNowConnected && settings.autoAcceptEnabled) {
            updateSettings({ autoAcceptEnabled: false });
            wasDisabledByDisconnectRef.current = true;
        }
    };

    const handleDisconnection = () => {
        if (settings.autoAcceptEnabled) {
            StopAutoAccept().catch(err => {
                console.error('Failed to stop auto-accept on disconnect:', err);
            });
            updateSettings({ autoAcceptEnabled: false });
            wasDisabledByDisconnectRef.current = true;
        }
    };

    const handleReconnection = () => {
        if (wasDisabledByDisconnectRef.current) {
            wasDisabledByDisconnectRef.current = false;
            if (settings.autoAcceptEnabled) {
                updateSettings({ autoAcceptEnabled: false });
            }
        } else {
            verifyServiceState();
        }
    };

    const handleAutoAcceptChange = async (enabled: boolean) => {
        try {
            updateSettings({ autoAcceptEnabled: enabled });

            if (!enabled) {
                await StopAutoAccept();
                wasDisabledByDisconnectRef.current = false;
                return;
            }

            wasDisabledByDisconnectRef.current = false;
            await StartAutoAccept({
                enabled: true,
                autoAccept: true,
            });
        } catch (err) {
            console.error('Failed to update auto-accept:', err);
            updateSettings({ autoAcceptEnabled: !enabled });
        }
    };

    return (
        <div className="tab-content">
            <div className="settings-card">
                <h3>Appearance</h3>
                <div className="setting-item">
                    <div className="setting-info">
                        <span className="setting-label">Theme</span>
                        <span className="setting-description">Choose your preferred color theme</span>
                    </div>
                    <select 
                        className="setting-select"
                        value={theme}
                        onChange={(e) => setTheme(e.target.value as 'light' | 'dark' | 'system')}
                    >
                        <option value="system">System</option>
                        <option value="light">Light</option>
                        <option value="dark">Dark</option>
                    </select>
                </div>
            </div>

            <div className="settings-card">
                <h3>Developer</h3>
                <div className="setting-item">
                    <div className="setting-info">
                        <span className="setting-label">Show Debug Tab</span>
                        <span className="setting-description">Display debug information panel in sidebar</span>
                    </div>
                    <label className="toggle-switch">
                        <input
                            type="checkbox"
                            checked={settings.showDebug}
                            onChange={(e) => updateSettings({ showDebug: e.target.checked })}
                        />
                        <span className="toggle-slider"></span>
                    </label>
                </div>

                <div className="setting-item">
                    <div className="setting-info">
                        <span className="setting-label">Auto Accept Match</span>
                        <span className="setting-description">
                            Automatically accept match ready checks. Turns off when client closes.
                        </span>
                    </div>
                    <label className="toggle-switch">
                        <input
                            type="checkbox"
                            checked={settings.autoAcceptEnabled}
                            disabled={!status?.connected}
                            onChange={(e) => handleAutoAcceptChange(e.target.checked)}
                        />
                        <span className="toggle-slider"></span>
                    </label>
                </div>
            </div>
        </div>
    );
}

