import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { GetConfig, IsConfigured } from "../../wailsjs/go/app/App";
import { config } from "../../wailsjs/go/models";

interface ConfigContextType {
    config: config.Config | null;
    isConfigured: boolean;
    loading: boolean;
}

const ConfigContext = createContext<ConfigContextType | null>(null);

export function ConfigProvider({ children }: { children: ReactNode }) {
    const [appConfig, setAppConfig] = useState<config.Config | null>(null);
    const [isConfigured, setIsConfigured] = useState(false);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function loadConfig() {
            try {
                const [configured, cfg] = await Promise.all([
                    IsConfigured(),
                    GetConfig(),
                ]);
                setIsConfigured(configured);
                setAppConfig(cfg);
            } catch (err) {
                console.error('Failed to load config:', err);
            } finally {
                setLoading(false);
            }
        }

        loadConfig();
    }, []);

    const value: ConfigContextType = {
        config: appConfig,
        isConfigured,
        loading,
    };

    return (
        <ConfigContext.Provider value={value}>
            {children}
        </ConfigContext.Provider>
    );
}

export function useConfig(): ConfigContextType {
    const context = useContext(ConfigContext);
    if (!context) {
        throw new Error('useConfig must be used within ConfigProvider');
    }
    return context;
}
