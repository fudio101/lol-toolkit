import React from 'react'
import { createRoot } from 'react-dom/client'
import './style.css'
import App from './App'
import { ConfigProvider, LCUProvider } from './contexts'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <ConfigProvider>
            <LCUProvider>
                <App />
            </LCUProvider>
        </ConfigProvider>
    </React.StrictMode>
)
