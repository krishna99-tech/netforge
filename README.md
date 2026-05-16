# NetForge v1.0.1 - Advanced Networking CLI Toolkit

**NetForge** is a professional-grade, high-performance network utility tool built with Go. It is designed for system administrators, developers, and network engineers who need a fast, modular, and cross-platform toolkit for diagnostics, monitoring, and security auditing.

## 🌟 Key Features

- **🌐 Deep Networking**: DNS, ICMP Ping, Concurrent Port Scanner (with banner grabbing), WHOIS, and Traceroute.
- **⚡ Performance Testing**: High-concurrency HTTP benchmarking with detailed latency statistics.
- **📊 Real-time Monitoring**: Live tracking of CPU usage, Bandwidth, Listening Ports, and detailed System/Network info.
- **🔒 Security Auditing**: SSL/TLS certificate inspection, and a full **TLS Cipher Suite Auditor** with security grading.
- **🔄 Network Proxy**: Built-in Reverse Proxy for traffic forwarding and testing.
- **🧰 Advanced Utilities**: Subnet Calculator, MAC OUI Lookup, IP Geolocation, Path MTU Discovery, and Socket Diagnostics.
- **🧬 Deep Intelligence**: DNS Fuzzer (subdomain brute-force), BGP/ASN Resolver, QUIC/HTTP3 Probe, and SNMP Walker.
- **🌊 Network Flow**: Lightweight NetFlow v5 collector for receiving and decoding flow records.
- **🦈 Packet Modules** *(requires Npcap/libpcap, build with `-tags pcap`)*: Live Packet Capture, ARP Topology Discovery, OS Fingerprinting, and PCAP Replay.
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

### 🧰 Advanced Utilities
```powershell
netforge ip info 8.8.8.8          # Geolocation & ISP lookup
netforge ip public                # Get your current public IP
netforge subnet 192.168.1.0/24   # CIDR Subnet calculator
netforge mac 00:1A:2B:3C:4D:5E   # OUI / Manufacturer lookup
netforge traceroute google.com   # Trace network path (hops)
netforge scan localhost --banner # Port scan + service banner grabbing
netforge mtu google.com          # Discover Path MTU (binary search)
```

### 🧬 Deep Intelligence
```powershell
netforge dnsfuzz brute domain.com --concurrency 100  # Subdomain brute-force
netforge bgp asn 8.8.8.8                             # BGP / ASN lookup
netforge tlsaudit google.com                          # Full TLS cipher audit
netforge tlsaudit google.com -o json                  # Export as JSON
netforge quic cloudflare.com                          # QUIC/HTTP3 probe
netforge snmp walk 192.168.1.1 --community public     # SNMP OID tree walk
netforge snmp get 192.168.1.1 --oid 1.3.6.1.2.1.1.1.0
netforge socket list                                  # Active socket table
netforge socket list --proto tcp --state ESTABLISHED
```

### 🌊 NetFlow Collection
```powershell
netforge netflow listen --port 2055   # Start NetFlow v5 collector
```

### 🦈 Packet Modules (Requires Npcap + `-tags pcap` build)
```powershell
netforge pcap capture -i eth0 -f "tcp port 80" -n 100
netforge topology scan --cidr 192.168.1.0/24 --iface eth0
netforge fingerprint 192.168.1.1 --port 80
netforge replay capture.pcap --iface eth0 --rate 2x
```

---

## 📄 Full Documentation
For a complete technical guide, architectural details, and a full command reference, please see [DOCUMENTATION.md](./DOCUMENTATION.md).

## 📝 License
This project is licensed under the MIT License.
