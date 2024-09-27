package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

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
