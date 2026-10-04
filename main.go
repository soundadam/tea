package main

import (
	"os"

	"github.com/soundadam/tea/cmd"
	"github.com/soundadam/tea/internal/privilege"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == privilege.InternalCommand {
		os.Exit(privilege.RunHelper(os.Args[2:]))
	}
	os.Exit(cmd.Execute())
}
