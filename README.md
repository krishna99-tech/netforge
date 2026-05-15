# NetForge - Advanced Networking CLI Toolkit

**NetForge** is a professional-grade, high-performance network utility tool built with Go. It is designed for system administrators, developers, and network engineers who need a fast, modular, and cross-platform toolkit for diagnostics, monitoring, and security auditing.

## 🌟 Key Features

- **🌐 Deep Networking**: DNS resolution, ICMP Ping (Privileged), Concurrent Port Scanning, and WHOIS lookups.
- **⚡ Performance Testing**: High-concurrency HTTP benchmarking with detailed latency statistics.
- **📊 Real-time Monitoring**: Live tracking of CPU usage, Bandwidth, Listening Ports, and detailed System/Network info.
- **🔒 Security Auditing**: Full SSL/TLS certificate chain inspection with expiration tracking.
- **🔄 Network Proxy**: Built-in Reverse Proxy for traffic forwarding and testing.
- **🛠 Modular Design**: Easily extensible architecture with professional command handling via Cobra.
- **📱 Cross-Platform**: Optimized binaries for Windows (x64, 386, ARM64), Linux, and macOS.

---

## 🛠 Installation & Build

### 1. Prerequisites
- Go 1.26 or higher
- (Windows) Administrator privileges for the `ping` command

### 2. Automated Installation (Remote)
To install NetForge instantly from GitHub, run this command in your PowerShell:
```powershell
irm https://raw.githubusercontent.com/krishna99-tech/netforge/main/install.ps1 | iex
```

### 3. Build from Source
If you have Go installed and want to build it yourself:
```powershell
./install.ps1
```

### 3. Run the tool
```powershell
./bin/netforge.exe --help
```

---

## 📖 Usage Examples

### Network Monitoring
```powershell
netforge monitor network   # View WiFi SSID and Network Interfaces
netforge monitor cpu       # Real-time CPU load
netforge monitor bandwidth # Live upload/download speeds
```

### Performance & Security
```powershell
netforge http benchmark https://google.com -n 100 -c 10
netforge ssl inspect google.com
```

### Professional Output
All commands support structured JSON or Table output for professional reporting:
```powershell
netforge monitor system -o json
netforge dns google.com -o table
```

---

## 📄 Full Documentation
For a complete technical guide, architectural details, and a full command reference, please see [DOCUMENTATION.md](./DOCUMENTATION.md).

## 📝 License
This project is licensed under the MIT License.
