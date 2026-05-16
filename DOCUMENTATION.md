# NetForge v1.0.1 Production Documentation

This document provides a comprehensive technical overview, architectural details, and production deployment guide for the NetForge CLI toolkit.

## 1. Architectural Overview
NetForge is built with a modular architecture to ensure scalability, ease of maintenance, and professional standards.

### Core Components
- **CLI Layer (`cmd/`)**: Built on the **Cobra** framework. It manages command routing, flag parsing, and console output.
- **Logic Layer (`internal/`)**: Contains the business logic for each feature. This code is isolated and cannot be imported by external packages (Go `internal` convention).
- **Utility Layer (`internal/utils/`)**: Provides shared services like JSON/Table formatting.

### Design Principles
- **Concurrency**: High-performance commands (like `scan`, `dnsfuzz`, `http benchmark`) utilize Go's goroutines and worker pools.
- **Cross-Platform**: The codebase is platform-agnostic, with specialized handlers for Windows-specific features (e.g., WiFi SSID detection via `netsh`, socket diagnostics via `netstat`).
- **Structured Output**: All commands support `--output json` or `--output table` for automated pipelines.
- **Build Tags**: Phase 3 CGO features (`pcap`, `topology`, `fingerprint`, `replay`) are gated behind the `//go:build pcap` tag to prevent breaking standard Windows builds.

---

## 2. Production Build & Deployment

### Compilation
NetForge uses a specialized build script to ensure correct environment settings.
- **Windows (standard)**: `./build.ps1` or `go build -o bin/netforge.exe .`
- **Windows (pcap features)**: `go build -tags pcap -o bin/netforge-pcap.exe .` (requires Npcap SDK + CGO)
- **Linux/Unix**: `make`

### Binaries
The build process generates optimized binaries in the `bin/` directory:
- `netforge-windows-amd64.exe` — Standard 64-bit Windows (no libpcap needed)
- `netforge-linux-amd64` — Standard 64-bit Linux
- `netforge-darwin-arm64` — Apple Silicon macOS

### Build Tags
| Tag | Effect |
|-----|--------|
| *(none)* | Default build — all pure-Go modules included |
| `-tags pcap` | Adds packet capture, ARP topology, OS fingerprint, and PCAP replay modules |

### Environment Variables
NetForge respects standard Go environment variables for cross-compilation:
- `GOOS`: Target Operating System
- `GOARCH`: Target Architecture
- `CGO_ENABLED=1`: Required for `-tags pcap` builds

---

## 3. Full Command Reference

### 🌐 Networking
| Command | Description | Example |
|---------|-------------|---------|
| `dns` | DNS resolution | `netforge dns google.com` |
| `ping` | ICMP latency check | `netforge ping 8.8.8.8` |
| `scan` | Concurrent port scanner | `netforge scan host --banner` |
| `tcp` | TCP connectivity test | `netforge tcp host 80` |
| `whois` | Domain WHOIS registry | `netforge whois google.com` |
| `traceroute` | Network path diagnostics | `netforge traceroute google.com` |
| `proxy` | Reverse HTTP proxy | `netforge proxy --from :8080 --to http://...` |

### 🧰 Advanced Utilities
| Command | Description | Example |
|---------|-------------|---------|
| `ip info` | Geolocation & ISP lookup | `netforge ip info 8.8.8.8` |
| `ip public` | Show public IP address | `netforge ip public` |
| `subnet` | CIDR subnet calculator | `netforge subnet 192.168.1.0/24` |
| `mac` | OUI / Manufacturer lookup | `netforge mac 00:1A:2B:3C:4D:5E` |
| `mtu` | Path MTU discovery | `netforge mtu google.com` |
| `socket list` | Active socket table | `netforge socket list --proto tcp` |

### 🧬 Deep Intelligence
| Command | Description | Example |
|---------|-------------|---------|
| `dnsfuzz brute` | Subdomain brute-force | `netforge dnsfuzz brute domain.com` |
| `bgp asn` | BGP / ASN resolver | `netforge bgp asn 8.8.8.8` |
| `tlsaudit` | TLS cipher audit + grading | `netforge tlsaudit google.com` |
| `quic` | QUIC / HTTP3 probe | `netforge quic cloudflare.com` |
| `snmp walk` | SNMP v1/v2c OID tree walk | `netforge snmp walk 192.168.1.1` |
| `snmp get` | SNMP single OID query | `netforge snmp get 192.168.1.1 --oid ...` |

### 📊 Monitoring
| Command | Description | Example |
|---------|-------------|---------|
| `monitor system` | OS & Memory info | `netforge monitor system` |
| `monitor network` | Interface & WiFi info | `netforge monitor network` |
| `monitor cpu` | Real-time CPU % | `netforge monitor cpu` |
| `monitor bandwidth` | Network speed | `netforge monitor bandwidth` |
| `monitor ports` | Listening port list | `netforge monitor ports` |

### ⚡ Performance & Security
| Command | Description | Example |
|---------|-------------|---------|
| `http get` | HTTP header inspection | `netforge http get https://...` |
| `http benchmark` | Concurrency load test | `netforge http benchmark url -n 100` |
| `ssl inspect` | SSL/TLS certificate info | `netforge ssl inspect google.com` |

### 🌊 Network Flow
| Command | Description | Example |
|---------|-------------|---------|
| `netflow listen` | NetFlow v5 UDP collector | `netforge netflow listen --port 2055` |

### 🦈 Packet Modules *(requires `-tags pcap`)*
| Command | Description | Example |
|---------|-------------|---------|
| `pcap capture` | Live packet capture | `netforge pcap capture -i eth0` |
| `topology scan` | ARP host discovery | `netforge topology scan --cidr 192.168.1.0/24` |
| `fingerprint` | OS TCP/IP fingerprinting | `netforge fingerprint 192.168.1.1` |
| `replay` | PCAP file replay | `netforge replay file.pcap --iface eth0` |

---

## 4. Troubleshooting
- **Ping/Socket Errors**: Ensure you are running as **Administrator** on Windows.
- **Mousetrap Message**: This tool is designed for the terminal. Avoid double-clicking the `.exe` in Explorer.
- **WiFi SSID "Disconnected"**: Ensure your WiFi is turned on and connected to an Access Point.
- **pcap/topology/fingerprint/replay not found**: These commands only appear if built with `go build -tags pcap`. They require Npcap (Windows) or libpcap (Linux/macOS).
- **SNMP timeout**: Ensure the target device has SNMP enabled and the community string matches.
- **NetFlow no data**: Configure your router to export flows to the machine's IP on the listening port.

---

## 5. Future Extensibility
To add a new feature:
1. Create a package in `internal/[feature]`.
2. Implement the logic.
3. Create a command file in `cmd/[feature].go` and register it in `init()`.
