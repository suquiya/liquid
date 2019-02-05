// Copyright © 2019 suquiya
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
	"github.com/spf13/cobra"
)

// newHeadCmd
func newHeadCmd() *cobra.Command {

	headCmd := &cobra.Command{
		Use:   "sethead",
		Short: "set license header to .go files in input directory or specified files.",
		Long:  `liquid head add header to .go files in input directory or  input specified files. If user specified files already have license header, liquid change header to specified license.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	headCmd.Flags().BoolP("directory", "d", true, "This flag shows whether input is directory or not (default true).")
	headCmd.Flags().BoolP("repeat", "r", false, "This flag decide whether add license to subdirectory recursively or not")
	headCmd.Flags().BoolP("file", "f", false, "If this flag is true, input paths are assumed files.")

	return headCmd
}
