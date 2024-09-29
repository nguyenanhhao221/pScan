package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nguyenanhhao221/pScan/scan"
)

func setUpFile(t *testing.T, initList bool, hosts []string) string {
	t.Helper()

	var tempDir = t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "pScan")
	if err != nil {
		t.Fatal(err)
	}

	if initList {
		hl := &scan.HostList{}
		for _, host := range hosts {
			if err := hl.Add(host); err != nil {
				t.Fatal(err)
			}
		}

		if err := hl.Save(tempFile.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return tempFile.Name()
}

func TestActions(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	testCases := []struct {
		name       string
		exp        string
		args       []string
		initList   bool
		actionFunc func(io.Writer, string, []string) error
	}{
		{
			name:       "ListAction",
			initList:   true,
			actionFunc: listAction,
			exp:        "host1\nhost2\nhost3\n",
		},
		{
			name:       "AddAction",
			initList:   false,
			actionFunc: addAction,
			exp:        "Added host: host1\nAdded host: host2\n",
			args:       []string{"host1", "host2"},
		},
		{
			name:       "DeleteAction",
			initList:   true,
			actionFunc: delAction,
			exp:        "Deleted host: host1\nDeleted host: host2\n",
			args:       []string{"host1", "host2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hostsFile := setUpFile(t, tc.initList, hosts)
			var b bytes.Buffer

			err := tc.actionFunc(&b, hostsFile, tc.args)
			if err != nil {
				t.Error(err)
			}

			out := b.String()
			if out != tc.exp {
				t.Errorf("Expect %s\n got %s", tc.exp, out)
			}
		})
	}
}

func TestScanAction(t *testing.T) {
	hosts := []string{"localhost", "invalidhost"}
	ports := []int{}
	tmpFile := setUpFile(t, true, hosts)

	// Set up ports, 1 open 1 close
	for i := 0; i < 2; i++ {
		ln, err := net.Listen("tcp", net.JoinHostPort("localhost", "0"))
		if err != nil {
			t.Fatal(err)
		}

		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		portNum, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}
		ports = append(ports, portNum)

		if i == 1 {
			ln.Close()
		}
	}

	var out bytes.Buffer
	err := scanAction(&out, tmpFile, ports)
	if err != nil {
		t.Errorf("Expect not error got %q", err)
	}

	var expectPrintOut string
	expectPrintOut += fmt.Sprintln("localhost:")
	expectPrintOut += fmt.Sprintf("\t%d: open\n", ports[0])
	expectPrintOut += fmt.Sprintf("\t%d: closed\n", ports[1])
	expectPrintOut += fmt.Sprintln()
	expectPrintOut += fmt.Sprintln("invalidhost: Host not found")
	expectPrintOut += fmt.Sprintln()

	got := out.String()
	if diff := cmp.Diff(expectPrintOut, got); diff != "" {
		t.Errorf("%s mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

func TestIntegration(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	hostsFile := setUpFile(t, false, hosts)
	var out bytes.Buffer

	if err := addAction(&out, hostsFile, hosts); err != nil {
		t.Fatalf("Expect no error, got: %v\n", err)
	}

	if err := listAction(&out, hostsFile, hosts); err != nil {
		t.Fatalf("Expect no error, got: %v\n", err)
	}

	hostToDel := []string{"host1"}

	if err := delAction(&out, hostsFile, hostToDel); err != nil {
		t.Fatalf("Expect no error, got: %v\n", err)
	}

	if err := listAction(&out, hostsFile, hosts); err != nil {
		t.Fatalf("Expect no error, got: %v\n", err)
	}

	if err := scanAction(&out, hostsFile, nil); err != nil {
		t.Fatalf("Expect no error, got: %v\n", err)
	}

	var expectOut string

	for _, h := range hosts {
		expectOut += fmt.Sprintf("Added host: %s\n", h)
	}

	expectOut += strings.Join(hosts, "\n")
	expectOut += fmt.Sprintln()
	hostsEnd := []string{"host2", "host3"}
	for _, h := range hostToDel {
		expectOut += fmt.Sprintf("Deleted host: %s\n", h)
	}
	expectOut += strings.Join(hostsEnd, "\n")
	expectOut += fmt.Sprintln()
	for _, h := range hostsEnd {
		expectOut += fmt.Sprintf("%s: Host not found\n", h)
		expectOut += fmt.Sprintln()
	}
	got := out.String()
	if diff := cmp.Diff(expectOut, got); diff != "" {
		t.Errorf("%s mismatch (-want +got):\n%s", t.Name(), diff)
	}
}
