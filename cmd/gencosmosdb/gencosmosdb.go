package main

import (
	"flag"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/jim-minter/go-cosmosdb/pkg/gencosmosdb"
)

var (
	pkg = flag.String("package", "cosmosdb", "package")

	packageRegexp          = regexp.MustCompile(`^package .*`)
	importRegexp           = regexp.MustCompile(`(?m)^\tpkg "[^"]+"$`)
	pluralRegexp           = regexp.MustCompile(`templates`)
	pluralExportedRegexp   = regexp.MustCompile(`Templates`)
	singularRegexp         = regexp.MustCompile(`template`)
	singularExportedRegexp = regexp.MustCompile(`Template`)
)

func writeFile(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	data = packageRegexp.ReplaceAll(data, []byte("// Code generated by github.com/jim-minter/go-cosmosdb, DO NOT EDIT.\n\npackage "+*pkg))

	_, err = f.Write(data)
	return err
}

func run() error {
	for _, name := range gencosmosdb.AssetNames() {
		if name == "template.go" {
			continue
		}

		err := writeFile("zz_generated_"+name, gencosmosdb.MustAsset(name))
		if err != nil {
			return err
		}
	}

	for _, arg := range flag.Args() {
		args := strings.Split(arg, ",")

		importpkg := args[0]
		singularExported := args[1]
		pluralExported := singularExported + "s"
		if len(args) == 3 {
			pluralExported = args[2]
		}
		singular := strings.ToLower(singularExported)
		plural := strings.ToLower(pluralExported)

		data := gencosmosdb.MustAsset("template.go")

		data = importRegexp.ReplaceAll(data, []byte("\tpkg \""+importpkg+"\""))

		// plural must be done before singular ("template" is a sub-string of "templates")
		data = pluralRegexp.ReplaceAll(data, []byte(plural))
		data = pluralExportedRegexp.ReplaceAll(data, []byte(pluralExported))
		data = singularRegexp.ReplaceAll(data, []byte(singular))
		data = singularExportedRegexp.ReplaceAll(data, []byte(singularExported))

		err := writeFile("zz_generated_"+singular+".go", data)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
