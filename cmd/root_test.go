package cmd

import (
	"bytes"
	"strings"
	"testing"

	appkit "github.com/soundadam/appkit/go"
)

func TestVersionReportsTheBuildVersion(t *testing.T) {
	info := appkit.App{Name: "tea", Version: "9.8.7", Repo: "owner/tea"}
	for _, args := range [][]string{{"version"}, {"--version"}} {
		root := NewRoot(info)
		buf := new(bytes.Buffer)
		root.SetOut(buf)
		root.SetErr(buf)
		root.SetArgs(args)
		if err := root.Execute(); err != nil {
			t.Fatal(err)
		}
		if got := strings.TrimSpace(buf.String()); got != "tea 9.8.7" {
			t.Fatalf("%v: got %q", args, got)
		}
	}
}
