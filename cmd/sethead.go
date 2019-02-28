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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// newHeadCmd
func newAddCmd() *cobra.Command {

	headCmd := &cobra.Command{
		Use:   "sethead [directoryPath]",
		Short: "add license header to .go files in input directory or specified files.",
		Long:  `liquid head add header to .go files in input directory or  input specified files. If user specified files already have license header, liquid change header to specified license.`,
		Run: func(cmd *cobra.Command, args []string) {
			config, license, author, LIsNotSet := ProcessArg(cmd, args)
			var input []string
			if cmd.Flags().NArg() < 1 {
				input = make([]string, 1, 1)
				var err error
				input[0], err = os.Getwd()
				if err != nil {
					cmd.Println("input is empty and current directory cannnot be gotten.")
					cmd.Println(err)
				}
			} else {
				input = cmd.Flags().Args()
			}

			r, err := cmd.Flags().GetBool("recursively")

			if err != nil {
				panic(err)
			}

			if r {
				rinput := make([]string, 0, 0)
				for _, dir := range input {
					rinput = append(rinput, dir)
					err := filepath.Walk(dir, func(p string, fi os.FileInfo, err error) error {
						if fi.IsDir() {
							rinput = append(rinput, p)
						}
						return nil
					})
					if err != nil {
						cmd.Println(err)
					}
				}

				input = rinput
			}
			for _, dirPath := range input {
				err := SetHeaderLicense(dirPath, license, author, cmd.OutOrStdout(), LIsNotSet, config)
				if err != nil {
					cmd.Println(err)
				}
			}
		},
	}

	headCmd.Flags().BoolP("directory", "d", true, "This flag shows whether input is directory or not (default true).")
	headCmd.Flags().BoolP("recursively", "r", false, "This flag decide whether add license to subdirectory recursively or not. default is false")
	headCmd.Flags().BoolP("file", "f", false, "If this flag is true, input paths are assumed files.")

	return headCmd
}

//SetHeaders is add license header to files that do not have license header and change files' license header if the files already have license header.
func SetHeaders(dirPath string, l *License, author string, messageW io.Writer, LIsNotSet bool, config *Config) error {

	license := l
	if LIsNotSet {
		license = GetDirLicense(dirPath)
	}

	sfis, err := ioutil.ReadDir(dirPath)

	if err != nil {
		return err
	}

	for _, file := range sfis {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".go" {
			fp := filepath.Join(dirPath, file.Name())
		}
	}

	return err
}

//SetFileHeader set file header to specified license.
func SetFileHeader(fp string, fi os.FileInfo, l *License) {

}
