package scan_test

import (
	"errors"
	"os"
	"slices"
	"testing"

	"github.com/nguyenanhhao221/go-cobra/scan"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name     string
		host     string
		hostList scan.HostList
		exp      []string
		expErr   error
	}{
		{name: "AddSuccess", host: "baz", exp: []string{"bar", "foo", "baz"}, expErr: nil, hostList: scan.HostList{Hosts: []string{"foo", "bar"}}},
		{name: "AddFailExist", host: "foo", expErr: scan.ErrExists, hostList: scan.HostList{Hosts: []string{"foo", "bar"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.hostList.Add(tc.host)
			if tc.expErr != nil {
				if err == nil {
					t.Error("Expect error, got 'nil'")
					return
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expect error: %q, got %q instead", tc.expErr, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Expect no error, got error: %q", err)
			}

			if !slices.Equal(tc.exp, tc.hostList.Hosts) {
				t.Errorf("Expect: %v, got %v\n", tc.exp, tc.hostList.Hosts)
			}
		})
	}
}
func TestRemove(t *testing.T) {
	testCases := []struct {
		name     string
		host     string
		hostList scan.HostList
		exp      []string
		expErr   error
	}{
		{name: "RemoveSuccess", host: "foo", exp: []string{"bar", "baz"}, expErr: nil, hostList: scan.HostList{Hosts: []string{"foo", "bar", "baz"}}},
		{name: "RemoveFailNotFound", host: "foo", expErr: scan.ErrNotExists, hostList: scan.HostList{Hosts: []string{"bar", "bar"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.hostList.Remove(tc.host)
			if tc.expErr != nil {
				if err == nil {
					t.Error("Expect error, got 'nil'")
					return
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expect error: %q, got %q instead", tc.expErr, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Expect no error, got error: %q", err)
			}

			if !slices.Equal(tc.exp, tc.hostList.Hosts) {
				t.Errorf("Expect: %v, got %v\n", tc.exp, tc.hostList.Hosts)
			}
		})
	}
}

func TestSaveLoad(t *testing.T) {
	hl1 := &scan.HostList{}
	hl2 := &scan.HostList{}

	hostName := "host1"
	if err := hl1.Add(hostName); err != nil {
		t.Errorf("Fail to add: %q", err)
	}

	tf, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatalf("Fail to set up temp file %s", err)
	}

	hostFile := tf.Name()

	if err := hl1.Save(hostFile); err != nil {
		t.Fatalf("Error saving host file %d", err)
	}

	if err := hl2.Load(hostFile); err != nil {
		t.Fatalf("Not loading host files %d", err)
	}

	if !slices.Equal(hl1.Hosts, hl2.Hosts) {
		t.Errorf("Host %v should match %v\n", hl1.Hosts, hl2.Hosts)
	}
}

func TestLoadNoFile(t *testing.T) {
	tf, err := os.CreateTemp(t.TempDir(), "")

	if err != nil {
		t.Fatalf("Fail to set up temp file %s", err)
		return
	}
	hostFile := tf.Name()
	if err := os.Remove(hostFile); err != nil {
		t.Fatalf("Error removing temp file %s\n", err)
	}

	hl := scan.HostList{}

	if err := hl.Load(hostFile); err != nil {
		t.Errorf("Expect no error, got %q instead \n", err)
	}
}
