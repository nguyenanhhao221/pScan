package scan_test

import (
	"errors"
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
