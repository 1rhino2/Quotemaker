# Random Quote Generator

[![Go](https://img.shields.io/badge/Go-1.22-blue?logo=go)](https://golang.org/) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE) [![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](#)

A simple Go web app that fetches a random quote from public APIs and credits [1rhino2](https://github.com/1rhino2). Features a modern UI, dark mode, author avatars, and a robust refresh button.

## Features

- Fetches a random quote from multiple free APIs (no API key needed)
- Minimal, responsive UI
- Dark mode toggle (🌙, persists across reloads)
- Author avatars (Gravatar)
- Refresh button for new quotes (AJAX, disables while loading)
- `/api/quote` endpoint for programmatic access
- Favicon support
- 404 page for unknown routes

## Usage

1. Install Go ([download](https://golang.org/dl/))
2. Clone this repo or copy the files
3. In the project folder, run:
   ```sh
   go run main.go
   ```
4. Open [http://localhost:8080](http://localhost:8080) in your browser

## License

MIT
