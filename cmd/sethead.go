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
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suquiya/liquid/tools"
)

// newHeadCmd
func newHeadCmd() *cobra.Command {

	headCmd := &cobra.Command{
		Use:   "sethead [Paths of files or directories]",
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
				for _, inputPath := range input {
					ii, err := os.Stat(inputPath)
					if err != nil {
						cmd.Println(err)
					} else {

						rinput = append(rinput, inputPath)
						if ii.IsDir() {
							err := filepath.Walk(inputPath, func(p string, fi os.FileInfo, err error) error {
								if fi.IsDir() {
									rinput = append(rinput, p)
								}
								return nil
							})
							if err != nil {
								cmd.Println(err)
							}
						}
					}
				}

				input = rinput
			}
			for _, inputPath := range input {
				err := SetHeaderLicense(inputPath, license, author, cmd.OutOrStdout(), LIsNotSet, config)
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

//SetHeaderLicense is add license header to files that do not have license header and change files' license header if the files already have license header.
func SetHeaderLicense(inputPath string, l *tools.License, author string, messageW io.Writer, LIsNotSet bool, config *Config) error {

	license := l
	ii, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if LIsNotSet {
		if ii.IsDir() {
			license = tools.GetDirLicense(inputPath)
		} else {
			license = tools.GetDirLicense(filepath.Dir(inputPath))
		}
	}

	if ii.IsDir() {
		sfis, err := ioutil.ReadDir(inputPath)

		if err != nil {
			return err
		}

		for _, file := range sfis {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".go" {
				fp := filepath.Join(inputPath, file.Name())
				err := SetFileHeader(fp, file, license, author)
				if err != nil {
					fmt.Fprintln(messageW, err)
				} else {
					fmt.Fprintln(messageW, "added license header to ", fp, ".")
				}
			}
		}
	} else {
		err := SetFileHeader(inputPath, ii, license, author)
		if err != nil {
			return err
		}
		fmt.Fprintln(messageW, "added license header to ", inputPath, ".")
	}

	return err
}

//SetFileHeader set file header to specified license.
func SetFileHeader(fp string, fi os.FileInfo, l *tools.License, author string) error {
	f, err := os.Open(fp)
	if err != nil {
		return err
	}

	tmp, err := os.Create(fp + ".tmp")
	if err != nil {
		return err
	}
	tmpw := bufio.NewWriter(tmp)

	sc := bufio.NewScanner(f)

	loop := true

	text := ""
	s := true
	for loop {
		s = sc.Scan()
		if s {
			text = sc.Text()
			ttext := strings.TrimSpace(text)
			if ttext != "" && ttext != "//" {
				loop = false
			}
		} else {
			loop = s
		}
	}
	if strings.HasPrefix(text, "//") {
		//file begin with comment line
		ttext := strings.TrimPrefix(text, "//")
		ttext = strings.TrimSpace(ttext)
		if strings.HasPrefix(ttext, "Copyright") {
			//detect Copyright Comment block
			detecting := true
			for detecting {
				s := sc.Scan()
				if s {
					text = sc.Text()
					if !strings.HasPrefix(text, "//") {
						detecting = false
					}
				} else {
					text = ""
					detecting = false
				}
			}
			l.WriteLicenseHeader(tmpw, author)
			if text != "" {
				for sc.Scan() {
					tmpw.WriteString(sc.Text())
				}
			}
		} else {
			l.WriteLicenseHeader(tmpw, author)
			tmpw.WriteString(text)
			for sc.Scan() {
				tmpw.WriteString(sc.Text())
			}
		}
	} else if strings.HasPrefix(text, "/*") {
		var commentStack strings.Builder
		if strings.Contains(text, "*/") {
			ttext := strings.TrimPrefix(text, "/*")
			ttext = strings.TrimSpace(ttext)
			if strings.HasPrefix(ttext, "Copyright") {
				l.WriteLicenseHeader(tmpw, author)
				e := strings.Index(text, "*/")
				tmpw.WriteString(text[e+2:])
				for sc.Scan() {
					tmpw.WriteString(sc.Text())
				}
			}
		} else {
			commentStack.WriteString(text)
			//detecting comment block
			detecting := true
			for detecting {
				if sc.Scan() {
					commentStack.WriteString(sc.Text())
					if strings.Contains(sc.Text(), "*/") {
						detecting = false
					}
				} else {
					detecting = false
				}
			}
			comment := commentStack.String()
			tcomment := strings.TrimPrefix(comment, "/*")
			tcomment = strings.TrimSpace(tcomment)
			if strings.HasPrefix(tcomment, "Copyright") {
				l.WriteLicenseHeader(tmpw, author)
				for sc.Scan() {
					tmpw.WriteString(sc.Text())
				}
			} else {
				l.WriteLicenseHeader(tmpw, author)
				tmpw.WriteString(comment)
				for sc.Scan() {
					tmpw.WriteString(sc.Text())
				}
			}
		}

	} else {
		l.WriteLicenseHeader(tmpw, author)
		tmpw.WriteString(text)
		for sc.Scan() {
			tmpw.WriteString(sc.Text())
		}
	}
	tmpw.Flush()
	f.Close()
	tmp.Close()
	err = os.Rename(fp+".tmp", fp)
	if err != nil {
		return err
	}

	//fmt.Fprintln(messageW, "added license header to ", fp, ".")
	return err
}
