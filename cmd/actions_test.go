package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	got := out.String()
	if diff := cmp.Diff(expectOut, got); diff != "" {
		t.Errorf("%s mismatch (-want +got):\n%s", t.Name(), diff)
	}
}
