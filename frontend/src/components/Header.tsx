import { useConfig } from '../contexts';

interface HeaderProps {
    title: string;
}

export function Header({ title }: HeaderProps) {
    const { config } = useConfig();

    return (
        <header className="main-header">
            <h1 className="page-title">{title}</h1>
            <div className="header-actions">
                <span className="region-badge">{config?.region?.toUpperCase() || 'N/A'}</span>
            </div>
        </header>
    );
}

