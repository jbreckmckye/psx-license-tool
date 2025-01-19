package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/jbreckmckye/psx-license-tool/internal/psx"
)

type Region int

const (
	REGION_JAPAN Region = iota
	REGION_EUROPE
	REGION_US
	REGION_UNSET
)

func main() {
	log.SetPrefix("[psxlicensepatch] ")
	log.SetFlags(0)

	var args struct {
		BIN    string `arg:"positional,required" help:"path to a PSX disc image BIN"`
		Region string `arg:"--region" help:"Sets region string and / or padding. May be JP, EUR or US"`
		Text   string `arg:"--text" help:"Sets disc license text, overwriting region"`
		TMD    string `arg:"--tmd" help:"Path to TMD file to insert into license. Used for PSX boot logo"`
	}
	arg.MustParse(&args)

	region := parseRegion(args.Region)

	if len(args.Text) > 70 {
		log.Println("WARNING: Text is above 70 characters, will be truncated")
	}

	tmd := loadTMD(args.TMD)
	if tmd != nil {
		checkTMD(tmd)
	}

	_, err := os.Stat(args.BIN)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(args.BIN, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	license, err := psx.ReadLicense(file)
	if err != nil {
		log.Fatal(err)
	}

	if region != REGION_UNSET && args.Text == "" {
		// Region was passed but no custom text - overwrite with default text for region
		switch region {
		case REGION_JAPAN:
			{
				psx.PatchLicenseText(license, psx.JP_STRING[:], true)
				break
			}
		case REGION_EUROPE:
			{
				psx.PatchLicenseText(license, psx.EUR_STRING[:], false)
				break
			}
		case REGION_US:
			{
				psx.PatchLicenseText(license, psx.USA_STRING[:], false)
				break
			}
		}
	}

	if args.Text != "" {
		// Text was passed so overwrite region text sector
		psx.PatchLicenseText(license, []byte(args.Text), region == REGION_JAPAN)
	}

	if tmd != nil {
		psx.PatchLicenseTMD(license, tmd)
	}

	err = psx.PatchLicense(file, license)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("BIN was patched")
}

func parseRegion(input string) Region {
	switch input {
	case "":
		return REGION_UNSET
	case "JP":
		return REGION_JAPAN
	case "EUR":
		return REGION_EUROPE
	case "US":
		return REGION_US
	default:
		log.Fatal("Region must be one of JP, EUR or US")
		return REGION_UNSET // Not used
	}
}

func loadTMD(path string) []byte {
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		return data
	}
	return nil
}

func checkTMD(tmd []byte) {
	writable, overBy := psx.ValidateTMDSize(tmd)
	if !writable {
		msg := fmt.Sprintf("ERROR: The TMD file is too large to write into the disc (oversized by %v bytes)", overBy)
		log.Fatal(msg)
	}
	if overBy > 0 {
		msg := fmt.Sprintf("WARNING: The TMD file may be larger than the PSX BIOS will read (oversized by %v bytes)", overBy)
		log.Println(msg)
		log.Println("Continuing...")
	}
}
