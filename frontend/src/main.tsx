import React from 'react'
import { createRoot } from 'react-dom/client'
import './style.css'
import App from './App'
import { ConfigProvider, LCUProvider, ThemeProvider, SettingsProvider, ApiLogProvider } from './contexts'

const container = document.getElementById('root')
const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <ThemeProvider>
            <SettingsProvider>
                <ApiLogProvider>
                    <ConfigProvider>
                        <LCUProvider>
                            <App />
                        </LCUProvider>
                    </ConfigProvider>
                </ApiLogProvider>
            </SettingsProvider>
        </ThemeProvider>
    </React.StrictMode>
)
