# dnsping

A simple, fast DNS ping tool.

Query a DNS server repeatedly and measure:
- Response time (RTT)
- Packet loss
- DNS response size
- EDNS0 / DNSSEC / Recursion / Authoritative features
- EDNS0 UDP buffer size

---

## Usage

```bash
dnsping [options] <dns_server_ip>
```

**Options:**
- `-domain` : Domain to query (default: `example.com`)
- `-type`   : Query type (`A`, `AAAA`, `TXT`, `CNAME`) (default: `A`)
- `-count`  : Number of queries to send (default: `5`)
- `-timeout`: Timeout per query (default: `2s`)
- `-port`   : DNS server port (default: `53`)
- `-features`: Show EDNS0/DNSSEC/Recursion/Authoritative support (default: `false`)
- `-nocolor`: Disable colored output (default: `false`)
- `-version`: Show version and exit

---

## Example

```bash
dnsping -domain google.com -type A -count 5 -features 8.8.8.8
```

Output:

```
Pinging DNS server 8.8.8.8 for domain google.com with type A:

Reply from 8.8.8.8: time=24.5ms size=73 bytes
Reply from 8.8.8.8: time=25.0ms size=73 bytes
Reply from 8.8.8.8: time=24.0ms size=73 bytes
Reply from 8.8.8.8: time=24.3ms size=73 bytes
Reply from 8.8.8.8: time=24.8ms size=73 bytes

Features detected: [EDNS0, EDNS0 UDP Buffer Size: 1232, Recursion Available]

--- 8.8.8.8 dnsping statistics ---
5 packets transmitted, 5 received, 0% packet loss
rtt min/avg/max = 24.0ms/24.53ms/25.0ms
```

---

## Build

```bash
go build -o dnsping dnsping.go
```
or using Makefile:

```bash
make build
```