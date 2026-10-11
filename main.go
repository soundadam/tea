package main

import (
	"os"

	appkit "github.com/soundadam/appkit/go"

	"github.com/soundadam/tea/cmd"
	"github.com/soundadam/tea/internal/privilege"
	teaversion "github.com/soundadam/tea/internal/version"
)

// version is injected by release builds with -ldflags "-X main.version=<ver>".
var version = "dev"

func main() {
	teaversion.Current = version
	if len(os.Args) > 1 && os.Args[1] == privilege.InternalCommand {
		os.Exit(privilege.RunHelper(os.Args[2:]))
	}
	os.Exit(cmd.Execute(appkit.App{Name: "tea", Version: version, Repo: "soundadam/tea"}))
}
