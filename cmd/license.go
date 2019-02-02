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

//Package cmd is inner package of liquid.
package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra/cobra/cmd"
)

//OSSLicenses store license information of some OSS Licenses.
var OSSLicenses map[string]cmd.License

func init() {
	OSSLicenses = cmd.Licenses
}

//CreateCustomLicense create License struct from file
func CreateCustomLicense(headPath, textPath string) (*cmd.License, error) {
	h := headPath
	t := textPath
	e, err := IsExistFilePath(h)
	e2, err2 := IsExistFilePath(t)

	if e {
		if e2 {
			return &cmd.License{Name: "custom", Header: readFile(h), Text: readFile(t)}, nil
		}
		fmt.Println(err2)
		h = t
		hStr := readFile(h)
		return &cmd.License{Name: "custom", Header: hStr, Text: hStr}, nil
	}
	if e2 {
		t = h
		tStr := readFile(t)
		return &cmd.License{Name: "custom", Header: tStr, Text: tStr}, nil
	}

	return nil, fmt.Errorf("%s\r\n%s", err.Error(), err2.Error())

}

func readFile(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(b)
}
