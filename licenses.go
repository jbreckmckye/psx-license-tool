package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetPrefix("psx-license-tool: ")
	log.SetFlags(0)

	path, err := readArg(1)
	if err != nil {
		log.Fatal("Usage error: " + err.Error())
	}

	license, err := readLicense(path)
	if err != nil {
		log.Fatal(err)
	}

	printLicense(license)
}

func readLicense(path string) ([]XAForm1Sector, error) {
	sectors := make([]XAForm1Sector, LICENSE_SECTORS)

	f, err := os.Open(path)
	if err != nil {
		return sectors, err
	}
	defer f.Close()

	for i := range LICENSE_SECTORS {
		buffer := make([]byte, ISO_SECTOR_SIZE)
		offset := int64(ISO_SECTOR_SIZE * i)

		bytesRead, err := f.ReadAt(buffer, offset)
		if err == io.EOF {
			return sectors, errors.New("reached end of file too early, check the path is actually a disc BIN image")
		}
		if bytesRead < ISO_SECTOR_SIZE {
			return sectors, errors.New("error reading disc sector, read too few bytes")
		}

		sector, err := ParseSectorXAForm1(buffer)
		if err != nil {
			return sectors, err
		}

		sectors[i] = sector
	}

	return sectors, nil
}

func printLicense(license []XAForm1Sector) {
  sector := license[4]
	bytes := sector.Data[:70]
	fmt.Printf("License string:\n%q\n", bytes)
}

func readArg(n int) (string, error) {
	args := os.Args
	if len(args) < (n + 1) {
		return "", errors.New("not enough arguments passed")
	}

	return args[n], nil
}
