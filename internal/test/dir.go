package test

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	cwd_args = flag.String("cwd", "", "Root Folder Path")
)

func ChangeDirectory(t *testing.T) {
	var cwd string
	cwd, ok := os.LookupEnv("cwd")
	if !ok {
		flag.Parse()
		if cwd_args == nil {
			t.FailNow()
		} else {
			cwd = *cwd_args
		}
	}

	if cwd != "" {
		if err := os.Chdir(cwd); err != nil {
			fmt.Println("Chdir error:", err)
		}
	}

}
