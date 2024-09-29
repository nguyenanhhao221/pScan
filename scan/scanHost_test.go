package scan_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/nguyenanhhao221/pScan/scan"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}

	if ps.Open.String() != "closed" {
		t.Errorf("Expect %q, got %q\n", "closed", ps.Open.String())
	}

	ps.Open = true

	if ps.Open.String() != "open" {
		t.Errorf("Expect %q, got %q\n", "closed", ps.Open.String())
	}
}

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{"OpenPort", "open"},
		{"ClosePort", "closed"},
	}
	hl := &scan.HostList{}
	host := "localhost"
	if err := hl.Add(host); err != nil {
		t.Fatal(err)
	}
	ports := []int{}
	// Set up ports, 1 open 1 close
	for _, tc := range testCases {
		// he port 0 is a special value in network programming that tells the OS to choose any free port. Once the OS assigns a port, you can retrieve it using ln.Addr().
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}

		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)

		if tc.name == "ClosePort" {
			ln.Close()
		}
	}

	res := scan.Run(hl, ports)
	if len(res) != 1 {
		t.Fatalf("Expected 1 result, got %d instead\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("Expected %q, got %q instead\n", host, res[0].Host)
	}

	if res[0].NotFound {
		t.Errorf("Expected host %q to be found\n", host)
	}

	if len(res[0].PortStates) != 2 {
		t.Fatalf("Expected 2 port states, got %d instead\n", len(res[0].PortStates))
	}

	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expect %q, got %q\n", ports[i], res[0].PortStates[i].Port)
		}
		if res[0].PortStates[i].Open.String() != tc.expectState {
			t.Errorf("Expect port %d, to be %s\n", ports[i], tc.expectState)
		}
	}
}

// TestRunHostNotFound This test could be potentially broken depends on how the DNS server setup in the machine that run it
// For Example, some DNS server actually return an valid looking IP address instead of an error when it couldn't found a host.
// If that is the case, this test will fail
func TestRunHostNotFound(t *testing.T) {
	host := "foo.invalid.uiweyhriweu"
	hl := &scan.HostList{}

	if err := hl.Add(host); err != nil {
		t.Fatal(err)
	}

	res := scan.Run(hl, []int{})

	if len(res) != 1 {
		t.Fatalf("Expected 1 result, got %d instead\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("Expected %q, got %q instead\n", host, res[0].Host)
	}
	if !res[0].NotFound {
		t.Errorf("Expect host %q NOT to be found\n", host)
	}

	if len(res[0].PortStates) != 0 {
		t.Fatalf("Expected 0 port state, got %d instead\n", len(res[0].PortStates))
	}
}
