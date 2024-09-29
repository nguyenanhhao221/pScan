package scan

import (
	"fmt"
	"net"
	"time"
)

type state bool

func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

// PortState represent the scan for a single port
type PortState struct {
	Port int
	Open state
}

// Results represents the scan results for a single host
type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

// Run perform a port scan on a hosts list
func Run(hl *HostList, ports []int) []Results {
	res := make([]Results, 0, len(hl.Hosts))
	for _, host := range hl.Hosts {
		r := Results{
			Host: host,
		}
		// Perform DNS lookup to see if the host exists
		// NOTE: this function is different on machine depends on the DNS and Internet Service prodiver. In my case, I use Vietnam Viettel Internet and default DNS set up on MacOS.
		// When given a host, this LookupHost go to the machine DNS settings, it as for an IP address from the DNS server, due to the way Viettel DNS server behave, when an invalid host is not found,
		// It does return an error to us, instead, it return an IP Address, which make our function thought that it actually found the host.
		// We can change this by updating our network DNS to use other DNS server such as Google or Cloudflare
		if _, err := net.LookupHost(host); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}
		// Scan the provided ports if the host is found
		for _, port := range ports {
			r.PortStates = append(r.PortStates, scanPort(host, port))
		}
		res = append(res, r)
	}
	return res
}

// scanPort perform TCP scan on a single port and host
func scanPort(host string, port int) PortState {
	p := PortState{Port: port}
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return p
	}

	scanConn.Close()
	p.Open = true
	return p
}
