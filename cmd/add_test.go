package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	testFilePath := "./testdata/add/test.go"

	command := "liquid add " + testFilePath + " -l bsd -a author"
	t.Log("command:", command)
	args := strings.Split(command, " ")
	bout := new(bytes.Buffer)

	lcmd := newRootCmd()
	lcmd.SetOutput(bout)

	lcmd.SetArgs(args[1:])
	err := lcmd.Execute()

	g, _ := ioutil.ReadFile(testFilePath)
	fmt.Fprintln(bout, string(g))
	os.Remove(testFilePath)
	t.Log(bout.String())
	if err != nil {
		t.Error(err)
	}

}
