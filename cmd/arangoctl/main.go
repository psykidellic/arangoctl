package main

import (
	"github.com/psykidellic/arangoctl/cmd/arangoctl/subcmd"
)

var (
	Version = "dev"
)

func main() {
	subcmd.Execute(Version)
}
