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
	"io"

	"github.com/spf13/cobra"
	"github.com/yourfin/binappend"
)

var printNames bool

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read [appended_file] [name_of_data]",
	Short: "Read appended data by name from a file",
	Long: `Read appened data by name from appended_file, appened in
the same format that "binappend write" uses.
See github.com/yourfin/binappend for details on the format`,
	Run: func(cmd *cobra.Command, args []string) {
		if printNames {
			if len(args) != 1 {
				fmt.Fprintln(
					os.Stderr,
					"The --dump-table option takes exactly one argument.\n",
					"Given: ",
					len(args),
				)
			}
			extractor, err := binappend.MakeExtractor(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "make extractor err: ", err)
			}
			for _, name := range extractor.AvalibleData() {
				fmt.Println(name)
			}
		} else {
			if len(args) != 2 {
				fmt.Fprintln(os.Stderr, "read takes exactly two arguments\nAdd the --help option for help")
			}
			extractor, err := binappend.MakeExtractor(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "make extractor err: ", err)
				os.Exit(1)
			}
			reader, err := extractor.GetReader(args[1])
			if err != nil {
				fmt.Fprintln(os.Stderr, "make reader err: ", err)
				os.Exit(1)
			}
			_, err = io.Copy(os.Stdout, reader)
			if err != nil {
				fmt.Fprintln(os.Stderr, "reader err: ", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	readCmd.PersistentFlags().BoolVarP(&printNames, "dump-table", "d", false, "print all names that can be read from the file")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
