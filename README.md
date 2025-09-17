
# GhostifyBot

![Build](https://github.com/DoniLite/GhostifyBot/actions/workflows/build.yml/badge.svg)
![Test](https://github.com/DoniLite/GhostifyBot/actions/workflows/test.yml/badge.svg)

GhostifyBot is a powerful Go-based streaming automation tool that downloads media content from torrents, processes it using `ffmpeg`, and streams it directly to Telegram channels.

---

## Overview

GhostifyBot aims to provide an autonomous backend solution for distributing multimedia content using a combination of:

- ✅ Torrent downloading
- ✅ Media processing via `ffmpeg`
- ✅ Automated Telegram streaming

This project is designed with flexibility, automation, and performance in mind — targeting developers who want to streamline content delivery pipelines.

---

## Objectives

- **Fetch media from torrents**: Download videos or audio files through magnet links or `.torrent` files.
- **Process and transcode with ffmpeg**: Clean, cut, convert or reformat media using powerful ffmpeg commands.
- **Distribute via Telegram**: Stream or upload the processed content automatically to Telegram channels or groups via bots.
- **Keep everything asynchronous**: Efficient task management with goroutines and wait groups for concurrency.

---

## Technologies

> Core stack and packages used in this project:

| Purpose               | Tech                                                                 |
|-----------------------|----------------------------------------------------------------------|
| Programming Language  | [Go](https://golang.org)                                              |
| Torrent Handling      | [`anacrolix/torrent`](https://github.com/anacrolix/torrent)          |
| Telegram API          | [`go-telegram-bot-api`](https://github.com/go-telegram-bot-api/telegram-bot-api) |
| Media Processing      | [`ffmpeg`](https://ffmpeg.org) (CLI)                                 |
| Crawling/Automation   | [`Rod`](https://github.com/go-rod/rod)      |
| Concurrency Handling  | `sync.WaitGroup`, goroutines                                         |

You can see the full list in the [`go.mod`](https://github.com/DoniLite/GhostifyBot/blob/main/go.mod).

---

## Installation

> **Requirements:**

- Go 1.24 or higher
- `ffmpeg` installed and available in your `$PATH`
- A Telegram bot token and channel ID
- Make installed on your system

### Clone the repo

```bash
git clone https://github.com/DoniLite/GhostifyBot.git
cd GhostifyBot
````

### Install Dependencies

```bash
make install-deps
```

### Build The project

```bash
make build
```

### Run the bot

```bash
make run
```

> Make sure to configure your Telegram credentials and ffmpeg settings in a config file or environment variables (WIP).

---

## Features Roadmap

- [x] Torrent downloading via magnet or .torrent (Processing...)
- [x] ffmpeg integration for media processing (Processing...)
- [x] Telegram channel media delivery (Processing...)
- [ ] Rod integration for site crawling
- [ ] Web dashboard or CLI interface
- [ ] Playlist or bulk torrent handling
- [ ] Custom transcoding profiles

---

## 🤝 Contributing

Contributions are welcome! Whether you're fixing bugs, improving performance, or adding features — your help is appreciated.

### How to contribute

1. Fork the repo
2. Create a new branch: `git checkout -b feature/my-new-feature`
3. Make your changes
4. Commit: `git commit -am 'Add my feature'`
5. Push: `git push origin feature/my-new-feature`
6. Open a Pull Request targeting the `develop` branch

---

## 📂 Project Structure (WIP)

```bash
GhostifyBot/
├── cmd/               # CLI or entrypoint (future)
├── utils/             # Utilities
├── services/          # Event system, App logic (torrent, telegram, ffmpeg) etc.
├── assets/            # Media files (optional)
├── downloads/         # Downloading contents
├── main.go            # Application entrypoint
└── go.mod             # Module dependencies
```

---

---

## 🔐 Environment Variables

GhostifyBot requires a few environment variables to be set for proper operation:

| Variable Name         | Description                              |
|-----------------------|------------------------------------------|
| `TELEGRAM_BOT_TOKEN`  | Your Telegram bot token                  |
| `TELEGRAM_CHANNEL_ID` | The target channel ID (e.g., `@mychannel`) |
| `FFMPEG_PATH`         | (Optional) Custom path to ffmpeg binary |
| `TORRENT_TMP_DIR`     | (Optional) Temp directory for torrent data |

You can create a `.env` file at the root of your project:

```bash
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHANNEL_ID=@your_channel
FFMPEG_PATH=/usr/bin/ffmpeg
TORRENT_TMP_DIR=./downloads
```

## 📄 License

This project is open-source and under the MIT License.

---

## 🌐 Links

- GitHub: [GhostifyBot](https://github.com/DoniLite/GhostifyBot)

- ffmpeg: [https://ffmpeg.org](https://ffmpeg.org)

- Rod (crawler): [https://github.com/go-rod/rod](https://github.com/go-rod/rod)

- Telegram Bot API: [https://core.telegram.org/bots/api](https://core.telegram.org/bots/api)

---

> Feel free to reach out for ideas, suggestions, or contributions!
