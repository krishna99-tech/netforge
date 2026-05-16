# NetForge v2.0 — Advanced Networking Features

This document covers all new hardcore networking modules added to NetForge. Each feature is implemented as a Cobra sub-command following the existing modular architecture.

---

## 🆕 New Features Overview

| Module | Command Prefix | Description |
|---|---|---|
| Packet Capture | `netforge pcap` | Live raw packet capture and analysis |
| Network Topology | `netforge topology` | ARP-based LAN host discovery & mapping |
| TCP Fingerprinting | `netforge fingerprint` | OS detection via TCP/IP stack analysis |
| DNS Fuzzer | `netforge dnsfuzz` | Subdomain enumeration & DNS brute-force |
| BGP & ASN Lookup | `netforge bgp` | BGP route, ASN, and prefix resolution |
| MTU Discovery | `netforge mtu` | Path MTU discovery with binary search |
| Network Namespace | `netforge netns` | Linux network namespace inspection |
| Socket Diagnostics | `netforge socket` | Active socket table dump with process mapping |
| VLAN Scanner | `netforge vlan` | 802.1Q VLAN probing on local interfaces |
| Packet Replay | `netforge replay` | Replay PCAP files against a live target |
| TCP Session Hijack Sim | `netforge tcpsim` | Educational TCP sequence number analysis |
| QUIC/HTTP3 Probe | `netforge quic` | QUIC handshake analysis and HTTP/3 support check |
| TLS Cipher Audit | `netforge tlsaudit` | Full TLS cipher suite enumeration and grading |
| SNMP Walker | `netforge snmp` | SNMP v1/v2c/v3 OID tree walk and query |
| NetFlow Collector | `netforge netflow` | Lightweight NetFlow v5/v9/IPFIX receiver |

---

## 📦 New Dependencies

Add these to your `go.mod`:

```go
require (
    github.com/google/gopacket       v1.1.19  // Packet capture (libpcap)
    github.com/gosnmp/gosnmp         v1.37.0  // SNMP v1/v2c/v3
    golang.org/x/net                 v0.25.0  // Low-level net primitives
    github.com/miekg/dns             v1.1.61  // Advanced DNS operations
    github.com/jackpal/gateway       v1.0.15  // Default gateway detection
)
```

> **Linux/macOS**: Install `libpcap-dev` (Debian/Ubuntu) or `libpcap` (Homebrew) before building features that use `gopacket`.
> **Windows**: Install [Npcap](https://npcap.com) with WinPcap compatibility mode enabled.

---

## 1. 📡 Packet Capture (`netforge pcap`)

Live capture of raw packets on any network interface. Supports BPF filtering, live decoding, and PCAP file export. Requires root/Administrator.

### Installation Note

```bash
# Ubuntu / Debian
sudo apt-get install libpcap-dev

# macOS
brew install libpcap

# Windows — install Npcap from https://npcap.com
```

### Usage

```bash
# List available interfaces
netforge pcap list

# Capture 100 packets on eth0 with a BPF filter
sudo netforge pcap capture --iface eth0 --filter "tcp port 443" --count 100

# Save capture to a .pcap file
sudo netforge pcap capture --iface eth0 --out capture.pcap --duration 30s

# Decode and print packet details from a saved file
netforge pcap decode capture.pcap

# Print only HTTP host headers from live traffic
sudo netforge pcap capture --iface eth0 --filter "tcp port 80" --decode http
```

### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--iface`, `-i` | string | `eth0` | Interface to capture on |
| `--filter`, `-f` | string | `""` | BPF filter expression |
| `--count`, `-n` | int | `0` (unlimited) | Stop after N packets |
| `--duration` | duration | `0` (unlimited) | Stop after duration |
| `--out` | string | `""` | Write to .pcap file |
| `--decode` | string | `all` | Protocol decoder: `all`, `http`, `dns`, `tls` |
| `--promisc` | bool | `false` | Enable promiscuous mode |
| `--snaplen` | int | `65535` | Max bytes per packet |

### Implementation Skeleton

```go
// cmd/pcap.go
package cmd

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
)

func runPcapCapture(iface, filter string, count int) error {
    handle, err := pcap.OpenLive(iface, 65535, true, pcap.BlockForever)
    if err != nil {
        return err
    }
    defer handle.Close()

    if err := handle.SetBPFFilter(filter); err != nil {
        return err
    }

    src := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range src.Packets() {
        decodePacket(packet)
    }
    return nil
}

func decodePacket(packet gopacket.Packet) {
    if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
        tcp, _ := tcpLayer.(*layers.TCP)
        // process tcp.SrcPort, tcp.DstPort, tcp.Payload
    }
}
```

---

## 2. 🗺️ Network Topology Discovery (`netforge topology`)

Discovers all live hosts on a local subnet using ARP requests, then maps open ports and builds a topology graph. Outputs as table, JSON, or Graphviz DOT format.

### Usage

```bash
# Discover all hosts on default subnet
sudo netforge topology scan

# Scan a specific CIDR
sudo netforge topology scan --cidr 10.0.0.0/24

# Export as Graphviz DOT for visualization
sudo netforge topology scan --cidr 192.168.1.0/24 -o dot > topology.dot
dot -Tpng topology.dot -o topology.png

# Deep scan: ARP discovery + port scan per host
sudo netforge topology scan --cidr 192.168.1.0/24 --deep --ports 22,80,443,3389
```

### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--cidr` | string | auto-detected | Target subnet in CIDR notation |
| `--iface` | string | default | Interface to send ARP from |
| `--deep` | bool | `false` | Also port-scan each discovered host |
| `--ports` | string | `22,80,443` | Comma-separated ports for deep scan |
| `--timeout` | duration | `500ms` | ARP response timeout per host |

### Implementation Skeleton

```go
// cmd/topology.go
func arpScan(cidr string, iface *net.Interface) ([]Host, error) {
    hosts := []Host{}
    ips, _ := hostsInCIDR(cidr)

    handle, _ := pcap.OpenLive(iface.Name, 65535, true, pcap.BlockForever)
    defer handle.Close()

    for _, ip := range ips {
        sendARPRequest(handle, iface, ip)
    }
    // listen for ARP replies
    src := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range src.Packets() {
        if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
            arp := arpLayer.(*layers.ARP)
            if arp.Operation == layers.ARPReply {
                hosts = append(hosts, Host{IP: net.IP(arp.SourceProtAddress), MAC: net.HardwareAddr(arp.SourceHwAddress)})
            }
        }
    }
    return hosts, nil
}
```

---

## 3. 🔍 OS Fingerprinting (`netforge fingerprint`)

Performs passive and active TCP/IP stack fingerprinting to identify the remote host's operating system, similar to `nmap -O`. Analyzes TCP window size, TTL, IP flags, and option ordering.

### Usage

```bash
# Active TCP fingerprint (sends crafted probe packets)
sudo netforge fingerprint 192.168.1.1

# Fingerprint with confidence score output
sudo netforge fingerprint 192.168.1.1 --verbose

# Fingerprint multiple hosts
sudo netforge fingerprint --file hosts.txt
```

### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--port` | int | `80` | Probe port (choose an open port for best results) |
| `--probes` | int | `4` | Number of probe packets |
| `--verbose` | bool | `false` | Show raw TCP option analysis |
| `--file` | string | `""` | File with one IP per line |

### How It Works

NetForge sends crafted TCP SYN packets with specific option combinations and analyzes the SYN-ACK response:

- **Initial TTL**: `64` → Linux/macOS, `128` → Windows, `255` → Cisco IOS
- **TCP Window Size**: `65535` → macOS/BSD, `8192` → older Windows, `29200` → Linux
- **TCP Options order**: MSS, NOP, WScale, NOP, NOP, Timestamp, SACK → Linux fingerprint
- **DF bit behavior**: Windows sets DF; many embedded systems do not

---

## 4. 🌐 DNS Fuzzer & Subdomain Enumerator (`netforge dnsfuzz`)

High-concurrency DNS brute-force and subdomain enumeration. Supports wordlist-based discovery, wildcard detection, and zone transfer attempts.

### Usage

```bash
# Subdomain brute-force with built-in wordlist
netforge dnsfuzz brute example.com

# Use a custom wordlist
netforge dnsfuzz brute example.com --wordlist /path/to/subdomains.txt --concurrency 200

# Attempt DNS zone transfer
netforge dnsfuzz axfr example.com

# Detect wildcard DNS entries
netforge dnsfuzz wildcard example.com

# Enumerate all DNS record types for a domain
netforge dnsfuzz records example.com --types A,AAAA,MX,TXT,NS,SOA,SRV,CAA
```

### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--wordlist` | string | built-in 10k list | Path to wordlist file |
| `--concurrency`, `-c` | int | `100` | Parallel DNS query workers |
| `--resolver` | string | system default | Custom DNS resolver (e.g. `8.8.8.8:53`) |
| `--timeout` | duration | `3s` | Per-query timeout |
| `--types` | string | `A,AAAA` | Record types to query |
| `--delay` | duration | `0` | Delay between queries (rate limiting) |

### Implementation Skeleton

```go
// cmd/dnsfuzz.go
func bruteforceSubdomains(domain string, wordlist []string, concurrency int) []Result {
    sem := make(chan struct{}, concurrency)
    results := make(chan Result, len(wordlist))

    for _, word := range wordlist {
        sem <- struct{}{}
        go func(sub string) {
            defer func() { <-sem }()
            fqdn := sub + "." + domain
            addrs, err := net.LookupHost(fqdn)
            if err == nil {
                results <- Result{Host: fqdn, IPs: addrs}
            }
        }(word)
    }
    // collect results...
}
```

---

## 5. 🛣️ BGP & ASN Lookup (`netforge bgp`)

Queries BGP routing tables and resolves Autonomous System Numbers (ASNs) for any IP or prefix. Uses Cymru's DNS-based whois service and the Routeviews prefix-to-AS dataset.

### Usage

```bash
# Resolve ASN for an IP address
netforge bgp asn 8.8.8.8

# Look up full BGP prefix info for a CIDR
netforge bgp prefix 8.8.8.0/24

# Get all prefixes announced by an ASN
netforge bgp prefixes AS15169

# Check if an IP is in a specific ASN
netforge bgp check 1.1.1.1 --asn AS13335
```

### Sample Output

```
IP:          8.8.8.8
ASN:         AS15169
AS Name:     GOOGLE, US
BGP Prefix:  8.8.8.0/24
Country:     US
Registry:    ARIN
Allocated:   1992-12-01
```

### Implementation

NetForge queries `origin.asn.cymru.com` via DNS TXT records:

```go
// cmd/bgp.go
func asnLookup(ip string) (*ASNInfo, error) {
    reversed := reverseIP(ip)
    query := reversed + ".origin.asn.cymru.com"

    c := new(dns.Client)
    m := new(dns.Msg)
    m.SetQuestion(dns.Fqdn(query), dns.TypeTXT)

    r, _, err := c.Exchange(m, "8.8.8.8:53")
    if err != nil {
        return nil, err
    }
    // parse TXT record: "15169 | 8.8.8.0/24 | US | arin | 1992-12-01"
    return parseCymruTXT(r.Answer), nil
}
```

---

## 6. 📏 Path MTU Discovery (`netforge mtu`)

Discovers the Maximum Transmission Unit along the path to a target using a binary search approach with ICMP/UDP probes. Detects PMTUD black holes.

### Usage

```bash
# Discover path MTU to a host
sudo netforge mtu discover google.com

# Test a specific MTU size
sudo netforge mtu test google.com --size 1400

# Detect PMTUD black holes
sudo netforge mtu blackhole google.com
```

### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--min` | int | `576` | Minimum probe size (bytes) |
| `--max` | int | `9000` | Maximum probe size (bytes, jumbo frames) |
| `--protocol` | string | `icmp` | Probe protocol: `icmp` or `udp` |
| `--timeout` | duration | `2s` | Probe timeout |

### Algorithm

```
1. Start with low=576, high=1500 (or --max for jumbo frame discovery)
2. Send probe of size (low+high)/2 with DF bit set
3. If ICMP "Fragmentation Needed" received → high = probe_size - 1
4. If reply received → low = probe_size
5. Repeat until low == high → that is the path MTU
```

---

## 7. 🔌 Socket Diagnostics (`netforge socket`)

Dumps the active socket table from the OS kernel, mapping connections to PID, process name, and binary path. Similar to `ss -tulnp` or `netstat -b`.

### Usage

```bash
# Show all listening TCP/UDP sockets
netforge socket list

# Show sockets for a specific process
netforge socket list --pid 1234
netforge socket list --process nginx

# Show only established TCP connections
netforge socket list --state ESTABLISHED --proto tcp

# Watch mode: refresh every second
netforge socket watch --interval 1s
```

### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--proto` | string | `all` | Filter: `tcp`, `udp`, `tcp6`, `udp6` |
| `--state` | string | `all` | Filter by state: `LISTEN`, `ESTABLISHED`, etc. |
| `--pid` | int | `0` | Filter by process ID |
| `--process` | string | `""` | Filter by process name (substring match) |
| `--interval` | duration | `1s` | Watch mode refresh interval |

### Platform Support

| Platform | Source | Notes |
|---|---|---|
| Linux | `/proc/net/tcp`, `/proc/net/udp` | Full PID mapping via `/proc/<pid>/fd` |
| macOS | `sysctl` + `kinfo_file` | Requires no special privileges |
| Windows | `iphlpapi.dll` via `syscall` | `GetExtendedTcpTable` |

---

## 8. 🔒 TLS Cipher Suite Auditor (`netforge tlsaudit`)

Performs a full TLS handshake negotiation test against every cipher suite and protocol version, then grades the configuration against Mozilla's TLS recommendations.

### Usage

```bash
# Full TLS audit with grading
netforge tlsaudit google.com

# Test a specific port (default 443)
netforge tlsaudit mail.example.com --port 587

# Check for specific vulnerabilities
netforge tlsaudit example.com --check heartbleed,poodle,beast,robot

# Export detailed JSON report
netforge tlsaudit example.com -o json > tls_report.json
```

### Sample Output

```
Host:      google.com:443
Grade:     A+
TLS 1.3:   ✓ Supported    Cipher: TLS_AES_256_GCM_SHA384
TLS 1.2:   ✓ Supported    Cipher: ECDHE-RSA-AES256-GCM-SHA384
TLS 1.1:   ✗ Disabled
TLS 1.0:   ✗ Disabled
SSLv3:     ✗ Disabled

Certificate:
  Subject:   *.google.com
  Issuer:    GTS CA 1C3
  Valid:     2024-03-01 → 2024-05-24 (8 days remaining ⚠)
  HSTS:      max-age=31536000; includeSubDomains; preload
  OCSP:      Stapled ✓

Vulnerability Checks:
  Heartbleed:  ✗ Not vulnerable
  POODLE:      ✗ Not vulnerable
  BEAST:       ✗ Not vulnerable
  ROBOT:       ✗ Not vulnerable
```

### Grading Criteria

| Grade | Criteria |
|---|---|
| A+ | TLS 1.3 only or TLS 1.2+ with PFS, HSTS, OCSP stapling |
| A | TLS 1.2+ with strong ciphers, no known vulnerabilities |
| B | TLS 1.2 supported but with some weak ciphers |
| C | TLS 1.1 or RC4 ciphers present |
| F | SSLv3 or critical vulnerability detected |

---

## 9. 📊 SNMP Walker (`netforge snmp`)

Queries SNMP agents using v1, v2c, or v3. Supports full OID tree walks, MIB resolution, and device interface statistics collection.

### Usage

```bash
# Walk the full OID tree (SNMP v2c)
netforge snmp walk 192.168.1.1 --community public

# Get a specific OID
netforge snmp get 192.168.1.1 --oid 1.3.6.1.2.1.1.1.0

# Use SNMP v3 with authentication
netforge snmp walk 192.168.1.1 \
  --version 3 \
  --user netforge \
  --auth-pass MyAuthPass \
  --priv-pass MyPrivPass \
  --auth-proto SHA \
  --priv-proto AES

# Collect interface statistics (ifTable)
netforge snmp interfaces 192.168.1.1 --community public

# Monitor interface bandwidth via SNMP polling
netforge snmp monitor 192.168.1.1 --community public --interval 5s
```

### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--version` | int | `2` | SNMP version: `1`, `2`, `3` |
| `--community` | string | `public` | SNMP v1/v2c community string |
| `--port` | int | `161` | SNMP agent port |
| `--oid` | string | `1.3.6.1.2.1` | Starting OID for walk |
| `--timeout` | duration | `5s` | Request timeout |
| `--retries` | int | `3` | Request retries |
| `--user` | string | `""` | SNMPv3 username |
| `--auth-proto` | string | `SHA` | SNMPv3 auth protocol: `MD5`, `SHA`, `SHA256` |
| `--priv-proto` | string | `AES` | SNMPv3 privacy protocol: `DES`, `AES`, `AES256` |

---

## 10. 🌊 NetFlow / IPFIX Collector (`netforge netflow`)

Starts a lightweight UDP listener that receives and decodes NetFlow v5, NetFlow v9, and IPFIX flow records. Aggregates top talkers and protocols in real time.

### Usage

```bash
# Start NetFlow v9/IPFIX collector on UDP 2055
netforge netflow listen --port 2055

# Listen and export decoded flows to JSON file
netforge netflow listen --port 9995 --out flows.json

# Show top-10 talkers updated every 5 seconds
netforge netflow listen --port 2055 --top 10 --interval 5s

# Filter flows by source or destination subnet
netforge netflow listen --port 2055 --filter-src 10.0.0.0/8
```

### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--port`, `-p` | int | `2055` | UDP port to listen on |
| `--version` | int | `0` (auto) | Force NetFlow version: `5`, `9`, or `10` (IPFIX) |
| `--out` | string | `""` | Write decoded flows to JSON file |
| `--top` | int | `10` | Show top N talkers |
| `--interval` | duration | `10s` | Stats refresh interval |
| `--filter-src` | string | `""` | Filter by source CIDR |
| `--filter-dst` | string | `""` | Filter by destination CIDR |

### Router Configuration

To send flows to NetForge from a Cisco IOS router:

```
ip flow-export version 9
ip flow-export destination <netforge_host_ip> 2055
ip flow-export source GigabitEthernet0/0
interface GigabitEthernet0/0
 ip route-cache flow
```

---

## 11. 🎯 QUIC / HTTP3 Probe (`netforge quic`)

Analyzes QUIC handshake parameters and checks HTTP/3 support, alt-svc advertisement, and QUIC version negotiation.

### Usage

```bash
# Check if a host supports QUIC/HTTP3
netforge quic probe google.com

# Full QUIC handshake analysis
netforge quic analyze cloudflare.com --verbose

# Check alt-svc header for HTTP/3 advertisement
netforge quic altsvc example.com
```

### Sample Output

```
Host:           cloudflare.com:443
HTTP/3:         ✓ Supported
QUIC Versions:  QUIC_VERSION_1, QUIC_VERSION_DRAFT_29
alt-svc:        h3=":443"; ma=86400
0-RTT:          ✓ Supported
Connection ID:  a1b2c3d4...
Handshake RTT:  12.4ms
```

---

## 12. 🔄 Packet Replay (`netforge replay`)

Replays a PCAP file against a live target, with configurable timing and rate scaling. Useful for reproducing network events and testing firewalls/IDS.

### Usage

```bash
# Replay a capture at original speed
sudo netforge replay capture.pcap --target 192.168.1.100

# Replay at 10x speed
sudo netforge replay capture.pcap --target 192.168.1.100 --rate 10x

# Replay only TCP connections on port 80
sudo netforge replay capture.pcap --target 192.168.1.100 --filter "tcp port 80"

# Rewrite source/dest IP addresses in the replay
sudo netforge replay capture.pcap \
  --target 192.168.1.100 \
  --rewrite-src 10.0.0.5 \
  --rewrite-dst 192.168.1.100
```

### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--target` | string | required | Destination IP to replay traffic to |
| `--iface` | string | default | Outbound interface |
| `--rate` | string | `1x` | Replay speed multiplier (e.g. `0.5x`, `10x`) |
| `--filter` | string | `""` | BPF filter to select packets |
| `--rewrite-src` | string | `""` | Rewrite source IP in packets |
| `--rewrite-dst` | string | `""` | Rewrite destination IP in packets |
| `--loop` | int | `1` | Number of times to loop the replay |

---

## 🏗 Project Structure (Updated)

```
netforge/
├── cmd/
│   ├── root.go
│   ├── pcap.go            # NEW: Packet capture
│   ├── topology.go        # NEW: ARP-based LAN topology
│   ├── fingerprint.go     # NEW: OS fingerprinting
│   ├── dnsfuzz.go         # NEW: DNS subdomain fuzzer
│   ├── bgp.go             # NEW: BGP/ASN lookup
│   ├── mtu.go             # NEW: Path MTU discovery
│   ├── socket.go          # NEW: Socket diagnostics
│   ├── tlsaudit.go        # NEW: TLS cipher audit
│   ├── snmp.go            # NEW: SNMP walker
│   ├── netflow.go         # NEW: NetFlow collector
│   ├── quic.go            # NEW: QUIC/HTTP3 probe
│   ├── replay.go          # NEW: PCAP replay
│   ├── dns.go
│   ├── ping.go
│   ├── scan.go
│   ├── http.go
│   ├── ssl.go
│   ├── monitor.go
│   └── ...
├── internal/
│   ├── pcap/              # NEW: Packet capture engine
│   ├── arp/               # NEW: ARP scanning utilities
│   ├── fingerprint/       # NEW: TCP/IP fingerprint DB
│   ├── snmp/              # NEW: SNMP session management
│   ├── netflow/           # NEW: NetFlow decoder
│   └── ...
├── data/
│   ├── fingerprints.json  # NEW: OS fingerprint signatures
│   └── subdomains.txt     # NEW: Built-in subdomain wordlist
├── go.mod
├── go.sum
└── install.ps1
```

---

## 🔒 Privilege Requirements

| Feature | Linux | macOS | Windows |
|---|---|---|---|
| Packet Capture (pcap) | `root` or `CAP_NET_RAW` | `root` | Administrator + Npcap |
| ARP Topology Scan | `root` or `CAP_NET_RAW` | `root` | Administrator + Npcap |
| OS Fingerprinting | `root` or `CAP_NET_RAW` | `root` | Administrator + Npcap |
| ICMP Ping | `root` or `CAP_NET_RAW` | `root` | Administrator |
| MTU Discovery | `root` or `CAP_NET_RAW` | `root` | Administrator |
| Packet Replay | `root` or `CAP_NET_RAW` | `root` | Administrator + Npcap |
| DNS Fuzzer | None | None | None |
| BGP/ASN Lookup | None | None | None |
| Socket Diagnostics | None (`/proc`) | None | None |
| TLS Audit | None | None | None |
| SNMP Walker | None | None | None |
| NetFlow Collector | None | None | None |

### Granting `CAP_NET_RAW` on Linux (without full root)

```bash
sudo setcap cap_net_raw,cap_net_admin=eip ./bin/netforge
```

---

## 🧪 Testing New Features

```bash
# Run all new feature tests
go test ./cmd/... -run TestPcap -v
go test ./cmd/... -run TestTopology -v
go test ./cmd/... -run TestDNSFuzz -v

# Integration test (requires network access)
go test ./cmd/... -tags integration -v

# Build with all features
go build -tags pcap -o bin/netforge .

# Build without libpcap dependency (disables pcap/topology/fingerprint/replay)
go build -tags nopcap -o bin/netforge-lite .
```

---

## 📦 Cross-Platform Build Matrix (Updated)

```powershell
# Windows builds (no libpcap features)
GOOS=windows GOARCH=amd64 go build -tags nopcap -o bin/netforge-windows-amd64.exe .

# Linux builds (full features with CGO for libpcap)
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/netforge-linux-amd64 .

# macOS builds
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o bin/netforge-darwin-arm64 .
```

> **Note**: Features using `gopacket`/libpcap require `CGO_ENABLED=1`. The `nopcap` build tag produces a CGO-free binary with those commands disabled, suitable for static linking or Windows builds without Npcap.

---

## 📄 License
MIT License — see [LICENSE](./LICENSE).