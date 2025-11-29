import { useSettings, useTheme, useLCU } from '../contexts';
import { StartAutoPick, StopAutoPick } from '../../wailsjs/go/app/App';
import { app } from '../../wailsjs/go/models';

export function SettingsTab() {
    const { settings, updateSettings } = useSettings();
    const { theme, setTheme } = useTheme();
    const { status } = useLCU();

    const handleAutoAcceptChange = async (enabled: boolean) => {
        try {
            // Update settings first to reflect the change immediately
            updateSettings({ autoAcceptEnabled: enabled });

            if (!enabled) {
                await StopAutoPick();
                return;
            }

            const config: app.AutoPickConfig = {
                enabled: true,
                autoAccept: true,
                autoPick: false,
                autoLock: false,
                championId: 0,
                championName: '',
            };

            await StartAutoPick(config);
        } catch (err) {
            console.error('Failed to update auto-accept:', err);
            // Revert on error
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
                        <span className="setting-label">Auto Accept Match (Test)</span>
                        <span className="setting-description">
                            Temporarily sends DECLINE instead of ACCEPT for safe testing.
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

