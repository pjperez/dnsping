package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

const version = "1.0"

var disableColor bool

func main() {
	// CLI flags
	count := flag.Int("count", 5, "Number of DNS pings to send")
	timeout := flag.Duration("timeout", 2*time.Second, "Timeout per DNS query")
	qtype := flag.String("type", "A", "DNS query type (A, AAAA, TXT, CNAME)")
	domain := flag.String("domain", "example.com", "Domain name to query")
	port := flag.Int("port", 53, "DNS server port")
	showFeatures := flag.Bool("features", false, "Show EDNS/DNSSEC/recursion/authoritative feature support")
	nocolor := flag.Bool("nocolor", false, "Disable colored output")
	showVersion := flag.Bool("version", false, "Show version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Println("dnsping version", version)
		return
	}

	if flag.NArg() < 1 {
		fmt.Println("Usage: dnsping [options] <server>")
		flag.PrintDefaults()
		return
	}

	server := flag.Arg(0)

	if net.ParseIP(server) == nil {
		fmt.Println("Invalid DNS server IP address")
		return
	}

	disableColor = *nocolor

	var rtts []time.Duration
	var sent, received int
	var detectedFeatures []string
	var featurePrinted bool

	fmt.Printf("Pinging DNS server %s for domain %s with type %s:\n\n", server, *domain, *qtype)

	for i := 0; i < *count; i++ {
		ok, duration, size, rcode, features := sendDNSQuery(server, *domain, *qtype, *timeout, *port, *showFeatures)
		sent++

		if ok {
			received++
			rtts = append(rtts, duration)
			fmt.Printf("%s: time=%v size=%d bytes\n", green("Reply from "+server), duration, size)
		} else {
			fmt.Printf("%s (rcode: %s)\n", red("Timeout from "+server), rcode)
		}

		if *showFeatures && !featurePrinted && len(features) > 0 {
			detectedFeatures = features
			featurePrinted = true
		}

		time.Sleep(1 * time.Second)
	}

	if *showFeatures && len(detectedFeatures) > 0 {
		fmt.Printf("\nFeatures detected: %v\n", detectedFeatures)
	}

	fmt.Printf("\n--- %s dnsping statistics ---\n", server)
	loss := 100 * (sent - received) / sent
	fmt.Printf("%d packets transmitted, %d received, %d%% packet loss\n", sent, received, loss)

	if len(rtts) > 0 {
		min, avg, max := rttStats(rtts)
		fmt.Printf("rtt min/avg/max = %v/%v/%v\n", min, avg, max)
	}
}

// sendDNSQuery sends a DNS query and returns useful results
func sendDNSQuery(server string, domain string, qtype string, timeout time.Duration, port int, showFeatures bool) (bool, time.Duration, int, string, []string) {
	c := new(dns.Client)
	c.Timeout = timeout

	m := new(dns.Msg)

	switch qtype {
	case "A":
		m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	case "AAAA":
		m.SetQuestion(dns.Fqdn(domain), dns.TypeAAAA)
	case "TXT":
		m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)
	case "CNAME":
		m.SetQuestion(dns.Fqdn(domain), dns.TypeCNAME)
	default:
		m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	}

	// Always attach EDNS0 request
	m.SetEdns0(4096, true)

	start := time.Now()
	r, _, err := c.Exchange(m, net.JoinHostPort(server, fmt.Sprintf("%d", port)))
	duration := time.Since(start)

	if err != nil {
		return false, duration, 0, "", nil
	}

	features := []string{}
	if showFeatures {
		if opt := r.IsEdns0(); opt != nil {
			features = append(features, "EDNS0")
			bufferSize := opt.UDPSize()
			features = append(features, fmt.Sprintf("EDNS0 UDP Buffer Size: %d", bufferSize))
		}
		if r.AuthenticatedData {
			features = append(features, "DNSSEC OK")
		}
		if r.RecursionAvailable {
			features = append(features, "Recursion Available")
		}
		if r.Authoritative {
			features = append(features, "Authoritative Answer")
		}
	}

	return r.Rcode == dns.RcodeSuccess, duration, r.Len(), dns.RcodeToString[r.Rcode], features
}

// rttStats computes min, avg, max RTT
func rttStats(rtts []time.Duration) (min, avg, max time.Duration) {
	min = rtts[0]
	max = rtts[0]
	var total time.Duration

	for _, rtt := range rtts {
		if rtt < min {
			min = rtt
		}
		if rtt > max {
			max = rtt
		}
		total += rtt
	}

	avg = total / time.Duration(len(rtts))
	return
}

// Color helpers
func green(s string) string {
	if disableColor {
		return s
	}
	return "[32m" + s + "[0m"
}

func red(s string) string {
	if disableColor {
		return s
	}
	return "[31m" + s + "[0m"
}
