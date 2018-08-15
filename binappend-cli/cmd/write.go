// Copyright © 2018 Patrick Nuckolls <nuckollsp at gmail>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourfin/binappend"
)

var compressed bool

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write [append-ee] [file to append] [file to append] ...",
	Short: "Append files to another file",
	Long: `Writes all [file to append] out to the end of [append-ee]
along with a metadata block describing where they are at the end.
See github.com/yourfin/binappend for details on the format`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Must provide at least two arguments\ntry `binappender-cli write --help` for help")
			os.Exit(1)
		}
		appender, err := binappend.MakeAppender(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Make appender err: ", err)
			os.Exit(1)
		}
		for _, filename := range args[1:] {
			err = appender.AppendFile(filename, compressed)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Append \"%s\" err: %s\n", filename, err)
				os.Exit(1)
			}
		}
		_ = appender.Close()
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)
	rootCmd.PersistentFlags().BoolVarP(&compressed, "compress", "c", false, "Compress the files")
}
