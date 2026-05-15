# NetForge Production Documentation

This document provides a comprehensive technical overview, architectural details, and production deployment guide for the NetForge CLI toolkit.

## 1. Architectural Overview
NetForge is built with a modular architecture to ensure scalability, ease of maintenance, and professional standards.

### Core Components
- **CLI Layer (`cmd/`)**: Built on the **Cobra** framework. It manages command routing, flag parsing, and console output.
- **Logic Layer (`internal/`)**: Contains the business logic for each feature. This code is isolated and cannot be imported by external packages (Go `internal` convention).
- **Utility Layer (`internal/utils/`)**: Provides shared services like JSON/Table formatting.

### Design Principles
- **Concurrency**: High-performance commands (like `scan` and `http benchmark`) utilize Go's goroutines and worker pools.
- **Cross-Platform**: The codebase is platform-agnostic, with specialized handlers for Windows-specific features (e.g., WiFi SSID detection via `netsh`).
- **Structured Output**: Support for JSON enables NetForge to be part of automated CI/CD pipelines.

---

## 2. Production Build & Deployment

### Compilation
NetForge uses a specialized build script to ensure correct environment settings.
- **Windows**: `./build.ps1`
- **Linux/Unix**: `make`

### Binaries
The build process generates optimized binaries in the `bin/` directory. For production, ensure you distribute the correct architecture:
- `netforge-windows-amd64.exe` (Standard 64-bit Windows)
- `netforge-linux-amd64` (Standard 64-bit Linux)

### Environment Variables
NetForge respects standard Go environment variables for cross-compilation:
- `GOOS`: Target Operating System
- `GOARCH`: Target Architecture

---

## 3. Command Reference & Examples

### Networking Group
| Command | Description | Example |
|---------|-------------|---------|
| `dns`   | DNS resolution | `netforge dns google.com` |
| `ping`  | ICMP latency check | `netforge ping 8.8.8.8` |
| `scan`  | Port scanner | `netforge scan localhost` |
| `tcp`   | TCP connectivity | `netforge tcp host port` |
| `whois` | Domain registry | `netforge whois google.com` |

### Monitoring Group
| Command | Description | Example |
|---------|-------------|---------|
| `monitor system` | OS & Memory info | `netforge monitor system` |
| `monitor network`| Interface & WiFi | `netforge monitor network` |
| `monitor cpu`    | Real-time CPU % | `netforge monitor cpu` |
| `monitor bandwidth`| Network speed | `netforge monitor bandwidth` |

### HTTP Group
| Command | Description | Example |
|---------|-------------|---------|
| `http get` | Header inspection | `netforge http get url` |
| `http benchmark`| Performance test | `netforge http benchmark url` |

---

## 4. Troubleshooting
- **Ping/Socket Errors**: Ensure you are running as **Administrator** on Windows.
- **Mousetrap Message**: This tool is designed for the terminal. Avoid double-clicking the `.exe` in Explorer.
- **WiFi SSID "Disconnected"**: Ensure your WiFi is turned on and connected to an Access Point.

---

## 5. Future Extensibility
To add a new feature:
1. Create a package in `internal/[feature]`.
2. Implement the logic.
3. Create a command file in `cmd/[feature].go` and register it in `init()`.
