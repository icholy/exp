// Command apidiff determines whether two versions of a package are compatible
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/types"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/icholy/exp/apidiff"
	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/go/packages"
)

var (
	incompatibleOnly = flag.Bool("incompatible", false, "display only incompatible changes")
	allowInternal    = flag.Bool("allow-internal", false, "allow apidiff to compare internal packages")
)

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage:\n")
		fmt.Fprintf(w, "apidiff OLD NEW\n")
		fmt.Fprintf(w, "   compares OLD and NEW package APIs\n")
		fmt.Fprintf(w, "   where OLD and NEW are import paths\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(2)
	}
	var oldpkg, newpkg *types.Package
	var g errgroup.Group
	g.Go(func() error {
		oldpkg = mustLoadPackage(flag.Arg(0)).Types
		return nil
	})
	g.Go(func() error {
		newpkg = mustLoadPackage(flag.Arg(1)).Types
		return nil
	})
	_ = g.Wait()
	if !*allowInternal {
		if isInternalPackage(oldpkg.Path()) && isInternalPackage(newpkg.Path()) {
			fmt.Fprintf(os.Stderr, "Ignoring internal package %s\n", oldpkg.Path())
			os.Exit(0)
		}
	}
	report := apidiff.Changes(oldpkg, newpkg)
	var err error
	if *incompatibleOnly {
		err = report.TextIncompatible(os.Stdout, false)
	} else {
		err = report.Text(os.Stdout)
	}
	if err != nil {
		log.Fatalf("writing report: %v", err)
	}
}

func mustLoadPackage(importPath string) *packages.Package {
	pkg, err := loadPackage(importPath, "")
	if err == nil {
		return pkg
	}
	pkg, err = downloadPackage(importPath)
	if err == nil {
		return pkg
	}
	log.Fatalf("loading %s: %v", importPath, err)
	panic("unreachable")
}

func loadPackage(importPath, dir string) (*packages.Package, error) {
	cfg := &packages.Config{
		Dir:  dir,
		Mode: packages.LoadTypes | packages.NeedName | packages.NeedTypes | packages.NeedImports | packages.NeedDeps,
	}
	pkgs, err := packages.Load(cfg, importPath)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("found no packages for import %s", importPath)
	}
	if len(pkgs[0].Errors) > 0 {
		return nil, pkgs[0].Errors[0]
	}
	return pkgs[0], nil
}

func isInternalPackage(pkgPath string) bool {
	switch {
	case strings.HasSuffix(pkgPath, "/internal"):
		return true
	case strings.Contains(pkgPath, "/internal/"):
		return true
	case pkgPath == "internal", strings.HasPrefix(pkgPath, "internal/"):
		return true
	}
	return false
}

func downloadPackage(importPath string) (*packages.Package, error) {
	tmp, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)
	cmd := exec.Command("go", "mod", "init", "tmp")
	cmd.Dir = tmp
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	var src bytes.Buffer
	src.WriteString("package tmp\n")
	src.WriteString("import _ \"" + importPath + "\"\n")
	src.WriteString("func main() {}\n")
	path := filepath.Join(tmp, "main.go")
	if err := os.WriteFile(path, src.Bytes(), os.ModePerm); err != nil {
		return nil, err
	}
	cmd = exec.Command("go", "get")
	cmd.Dir = tmp
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return loadPackage(importPath, tmp)
}
