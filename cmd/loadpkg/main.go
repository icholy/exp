package main

import (
	"flag"

	"golang.org/x/tools/go/packages"
)

func main() {
	flag.Parse()
	cfg := &packages.Config{
		Mode: packages.LoadTypes | packages.NeedName | packages.NeedTypes | packages.NeedImports | packages.NeedDeps,
	}
	pkgs, err := packages.Load(cfg, flag.Arg(0))
	if err != nil {
		panic(err)
	}
	packages.PrintErrors(pkgs)
}
