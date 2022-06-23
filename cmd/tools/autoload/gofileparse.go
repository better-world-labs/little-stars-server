package main

import (
	"aed-api-server/internal/pkg/utils"
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/tools/go/packages"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

type parseResult struct {
	Path        string
	PkgName     string
	InjectNames []string
}

const InjectTag = "inject-component"

var packageReg = regexp.MustCompile("^package ([a-zA-Z][a-zA-Z0-9_]*)")
var injectReg = regexp.MustCompile(fmt.Sprintf("^//go:%s(\\s+.*|$)", InjectTag))
var funcReg = regexp.MustCompile("^func\\s+([A-Z][a-zA-Z0-9_]*)\\s*\\(")

func goFileParse(goFilepath string) (*parseResult, error) {
	file, err := os.Open(goFilepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	var PkgName string
	var InjectNames = make([]string, 0)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if PkgName == "" {
			match := packageReg.FindSubmatch(line)
			if len(match) == 2 {
				PkgName = string(match[1])
			}
		} else {
			if injectReg.Match(line) {
				line, _, err := reader.ReadLine()
				if err == io.EOF {
					return nil, errors.New("file unexpected end")
				}
				math := funcReg.FindSubmatch(line)
				if len(math) == 2 {
					InjectNames = append(InjectNames, string(math[1]))
				}
			}
		}
	}

	if PkgName == "" || len(InjectNames) == 0 {
		return nil, nil
	}
	absPath, err := filepath.Abs(goFilepath)
	if err != nil {
		return nil, err
	}
	return &parseResult{
		Path:        path.Dir(absPath),
		PkgName:     PkgName,
		InjectNames: InjectNames,
	}, nil
}

func goModuleInfo(dir string) (moduleName string, moduleAbsPath string, err error) {
	defer utils.TimeStat("goModuleInfo:" + dir)()
	cfg := &packages.Config{
		Mode:  packages.NeedModule,
		Dir:   dir,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, "")

	if err != nil {
		return "", "", err
	}

	if len(pkgs) == 0 {
		return "", "", errors.New("not found go module")
	}

	p := pkgs[0]

	return p.Module.Path, p.Module.Dir, nil
}
