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
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/cobra/cmd"
)

// newCmd represents the new command
func newAddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add [filename]",
		Short: "create newfile of source code",
		Long:  `This command create new file of source code using specified license`,
		Run: func(cmd *cobra.Command, args []string) {
			_, license, author := ProcessArg(cmd, args)
			packageName, _ := cmd.Flags().GetString("package")
			input := cmd.Flags().Args()
			for _, fileName := range input {
				createNew(fileName, license, author, packageName, cmd.OutOrStdout())
			}
		},
	}

	addCmd.Flags().StringP("package", "p", "", "package name for go file")

	return addCmd
}

func createNew(fn string, l *cmd.License, author, packageName string, messageWriter io.Writer) {
	isFilePath, err := IsFilePath(fn)
	if isFilePath {
		fp, err := filepath.Abs(fn)
		if err != nil {
			panic(err)
		}
		isExist, err := IsExistFile(fp)
		if isExist {
			fmt.Fprintf(messageWriter, "%s is exist. liquid did not process about %s.\r\n", fp, fp)
			return
		}

		if err != nil {
			fmt.Fprintf(messageWriter, err.Error())
			fmt.Fprintln(messageWriter)
			fmt.Fprintf(messageWriter, "liquid did not process about %s.\r\n", fp)
			return
		}

		dir := filepath.Dir(fp)
		pn := packageName
		if pn == "" {
			pn = filepath.Base(dir)
		}

	} else {
		fmt.Println(err)
	}
}
