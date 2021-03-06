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

//Package tools is inner tools package of liquid.
package tools

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/cobra/cobra/cmd"
)

//OSSLicenses store license information of some OSS Licenses.
var OSSLicenses map[string]*License

//License represents License data
type License struct {
	Name   string
	Text   string
	Header string
}

func convertCLToLL(cl *cmd.License) *License {
	return &License{
		Name:   cl.Name,
		Header: cl.Header,
		Text:   cl.Text,
	}
}

func init() {
	OSSLicenses = make(map[string]*License)
	for key, cl := range cmd.Licenses {
		l := convertCLToLL(&cl)
		OSSLicenses[key] = l
		OSSLicenses[l.Name] = l
	}
}

//CreateCustomLicense create License struct from file
func CreateCustomLicense(headPath, textPath string) (*License, error) {
	h := headPath
	t := textPath
	e, err := IsExistFilePath(h)
	e2, err2 := IsExistFilePath(t)

	if e {
		if e2 {
			return &License{Name: "custom", Header: readFile(h), Text: readFile(t)}, nil
		}
		fmt.Println(err2)
		h = t
		hStr := readFile(h)
		return &License{Name: "custom", Header: hStr, Text: hStr}, nil
	}
	if e2 {
		t = h
		tStr := readFile(t)
		return &License{Name: "custom", Header: tStr, Text: tStr}, nil
	}

	return nil, fmt.Errorf("%s\r\n%s", err.Error(), err2.Error())

}

//GetOSSLicense get an OSSLicense struct from license name
func GetOSSLicense(licenseName string) *License {
	li, exist := OSSLicenses[licenseName]
	if !exist {
		err := fmt.Errorf("OSSLicenses not hit")
		fmt.Println(err)
		fmt.Println("liquid automatically choose mit")
		licenseName = "mit"
		li, _ = OSSLicenses[licenseName]
		//fmt.Println(OSSLicenses)
	}
	return li
}

func readFile(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

//WriteLicenseHeader write license header to w
func (l *License) WriteLicenseHeader(w io.Writer, author string) {
	ct := getNowCopyrightText(author)
	data := make(map[string]interface{})
	data["copyright"] = ct
	data["licenseHeader"] = l.Header

	template := `{{comment .copyright}}
{{comment .licenseHeader}}
`
	err := execTemplate(template, data, w)
	if err != nil {
		panic(err)
	}

}

func execTemplate(tmpl string, data interface{}, w io.Writer) error {
	t, err := template.New("").Funcs(template.FuncMap{"comment": CommentifyString}).Parse(tmpl)

	if err != nil {
		return err
	}

	return t.Execute(w, data)

}

//CommentifyString commentify string inspired by cobra's commentifyString
func CommentifyString(input string) string {
	nlcode := "\n"
	replacer := strings.NewReplacer("\r\n", nlcode, "\r", nlcode, "\n", nlcode)
	inputNLd := replacer.Replace(input)

	lines := strings.Split(inputNLd, "\n")
	var sb strings.Builder
	sb.Grow(len(input) + len(lines)*len("\n"))
	c := "// "
	for _, l := range lines {
		if strings.HasPrefix(l, c) {
			sb.WriteString(l)
			sb.WriteString(nlcode)
		} else {
			sb.WriteString(c)
			if l != "" {
				sb.WriteString(l)
			}
			sb.WriteString(nlcode)
		}
	}

	return strings.TrimSuffix(sb.String(), nlcode)
}

func getNowCopyrightText(author string) string {
	var sb strings.Builder
	sb.Grow(19 + len(author))
	sb.WriteString("Copyright (c) ")
	sb.WriteString(time.Now().Format("2006"))
	sb.WriteString(" ")
	sb.WriteString(author)
	return sb.String()
}

//GetDirLicense get license based on license file in dir
func GetDirLicense(dir string) *License {
	lc := findAndGetLicenseContent(dir)
	if lc == nil {
		return nil
	}

	var l *License
	l = nil
	lcStr := string(lc)
	lcStr = strings.TrimSpace(lcStr)
	for _, ol := range OSSLicenses {
		t := strings.TrimSpace(ol.Text)
		h := strings.TrimSpace(ol.Header)
		if strings.HasSuffix(lcStr, t) || strings.HasPrefix(lcStr, h) {
			l = ol
			break
		}
	}

	if l == nil {
		l = &License{
			Name:   "custom",
			Text:   lcStr,
			Header: lcStr,
		}
	}

	return l
}

func findAndGetLicenseContent(dir string) []byte {
	candidate := []string{filepath.Join(dir, "LICENSE"), filepath.Join(dir, "LICENSE.txt"), filepath.Join(dir, "LICENSE.md")}

	search := true

	i := 0
	var b []byte
	b = nil
	for search {
		if e, _ := IsExistFile(candidate[i]); e {
			search = false
			b, _ = ioutil.ReadFile(candidate[i])

		} else {
			i++
			if !(i < len(candidate)) {
				search = false
			}
		}
	}

	return b
}

//IsExistFilePath is validate whether val is exist filepath or not.
func IsExistFilePath(val string) (bool, error) {
	absPath, err := filepath.Abs(val)
	if err != nil {
		return false, err
	}
	if is, _ := govalidator.IsFilePath(absPath); !is {
		return is, fmt.Errorf("%s is not file path", val)
	}

	fi, err := os.Stat(val)
	if err != nil {
		return false, err
	}

	if fi.IsDir() {
		return false, fmt.Errorf("%s is not file", val)
	}

	return true, err
}

//IsExistFile validate whether val is exist or not
func IsExistFile(val string) (bool, error) {
	fi, err := os.Stat(val)

	if os.IsNotExist(err) {
		return false, nil
	}

	if fi.IsDir() {
		return false, fmt.Errorf("%s is directory", val)
	}

	return true, err
}

//IsFilePath validate whether val is file path or not
func IsFilePath(val string) (bool, error) {
	absPath, err := filepath.Abs(val)
	if err != nil {
		return false, err
	}

	if is, _ := govalidator.IsFilePath(absPath); !is {
		return is, fmt.Errorf("%s is not file path", val)
	}

	return true, nil
}
