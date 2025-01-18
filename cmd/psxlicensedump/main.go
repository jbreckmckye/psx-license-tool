package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/jbreckmckye/psx-license-tool/internal/psx"
)

func main() {
	log.SetPrefix("[psxlicensedump]")
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
	regionLength := 70

	if region == psx.EUR_STRING {
    fmt.Println("Detected European license")
	} else if region == psx.USA_STRING {
    fmt.Println("Detected American license")
	} else {
		japanMatch := [65]byte(region[:65])
    if japanMatch == psx.JP_STRING {
      fmt.Println("Detected Japanese license")
			regionLength = 65
		} else {
			fmt.Println("Unknown license type? Check file is a PSX disc image BIN. Attempting to continue...")
		}
	}

	tmd := psx.GetLicenseTMD(license)

	err = os.WriteFile(args.Out+".TXT", region[:regionLength], 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(args.Out+".TMD", tmd, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Dumped license data to %v.TXT, %v.TMD\n", args.Out, args.Out)
}
