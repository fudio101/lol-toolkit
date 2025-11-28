import { useSettings, useTheme } from '../contexts';

export function SettingsTab() {
    const { settings, updateSettings } = useSettings();
    const { theme, setTheme } = useTheme();

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
            </div>
        </div>
    );
}

