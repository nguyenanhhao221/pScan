/*
Copyright © 2024 Hao Nguyen <hao@haonguyen.tech>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/nguyenanhhao221/pScan/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:          "delete <host1>...<hostn>",
	Short:        "delete host(s) from the host's list",
	Aliases:      []string{"d"},
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")
		return delAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(deleteCmd)
}

func delAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, host := range args {
		if err := hl.Remove(host); err != nil {
			return err
		}
		fmt.Fprintln(out, "Deleted host:", host)
	}
	return hl.Save(hostsFile)
}
