package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/jbreckmckye/psx-license-tool/internal/psx"
)

const LICENSE_SECTORS = 16

func main() {
	log.SetPrefix("psxlicensedump | ")
	log.SetFlags(0)

	var args struct {
		BIN string `arg:"positional,required" help:"path to a PSX disc image BIN"`
		Out string `arg:"--output" help:"name for .TXT, .TMD output files" default:"LICENSE"`
	}
	arg.MustParse(&args)

	file, err := os.Open(args.BIN)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	license, err := psx.ReadLicense(file)
	if err != nil {
		log.Fatal(err)
	}

	region := psx.GetLicenseText(license)
	tmd := psx.GetLicenseTMD(license)

	os.WriteFile(args.Out+".TXT", region, 0644)
	os.WriteFile(args.Out+".TMD", tmd, 0644)
}
