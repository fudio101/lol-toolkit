# LoL Toolkit

A League of Legends toolkit desktop application built with Go + React.

## Features

- 🔍 Summoner search by Riot ID
- 📊 Ranked stats viewer
- 🏆 Champion mastery viewer
- 🎮 League leaderboards

## Prerequisites

- [Go 1.23+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

## Setup

### 1. Create Config File (Required)

Create `internal/config/config.json`:

```json
{
  "riot_api_key": "",
  "region": "vn2"
}
```

### 2. Add Your Riot API Key

Get your API key from [Riot Developer Portal](https://developer.riotgames.com/).

```json
{
  "riot_api_key": "RGAPI-your-key-here",
  "region": "vn2"
}
```

> ⚠️ This file is gitignored. Your API key will be embedded in the built executable.
> 
> ⚠️ Development API keys expire every **24 hours**. Regenerate at [developer.riotgames.com](https://developer.riotgames.com/).

### Available Regions

| Code | Region |
|------|--------|
| `vn2` | Vietnam |
| `na1` | North America |
| `euw1` | Europe West |
| `kr` | Korea |
| `jp1` | Japan |
| `sea` | Southeast Asia |

## Development

```bash
wails dev
```

## Building

```bash
wails build
```

## Project Structure

```
lol-toolkit/
├── main.go
├── internal/
│   ├── app/                     # App logic (exposed to frontend)
│   ├── config/
│   │   ├── config.go
│   │   └── config.json          # ← CREATE THIS FILE
│   └── lol/                     # Riot API client
├── frontend/                    # React + TypeScript
└── wails.json
```

## License

MIT
