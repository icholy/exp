package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/tools/go/packages"
)

func main() {
	flag.Parse()
	cfg := &packages.Config{
		Mode: packages.LoadTypes | packages.NeedName | packages.NeedTypes | packages.NeedImports | packages.NeedDeps,
	}
	pkgs, err := packages.Load(cfg, flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) == 0 {
		log.Fatal("no matching packages found")
	}
	for _, pkg := range pkgs {
		fmt.Println(pkg.PkgPath)
		for _, err := range pkg.Errors {
			log.Printf("ERROR %v\n", err)
		}
	}
}
