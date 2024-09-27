// Package scan provides types and functions to perform TCP port
// scan on a list of hosts
package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
)

var (
	ErrExists    = errors.New("host already exist in the list")
	ErrNotExists = errors.New("host not in the list")
)

type HostList struct {
	Hosts []string
}

func (hl *HostList) search(host string) (bool, int) {
	slices.Sort(hl.Hosts)

	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}
	return false, -1
}

func (hl *HostList) Add(host string) error {
	found, _ := hl.search(host)
	if found {
		return ErrExists
	}
	hl.Hosts = append(hl.Hosts, host)
	return nil
}

func (hl *HostList) Remove(host string) error {
	found, i := hl.search(host)
	if found {
		hl.Hosts = slices.Delete(hl.Hosts, i, i+1)
		return nil
	}
	return ErrNotExists
}

func (hl *HostList) Save(hostFile string) error {
	var output string
	for _, host := range hl.Hosts {
		output += fmt.Sprintln(host)

	}
	return os.WriteFile(hostFile, []byte(output), 0644)
}

func (hl *HostList) Load(hostFile string) error {
	f, err := os.Open(hostFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}

	return nil
}
