package psx

import (
	"errors"
	"io"
	"os"
	"slices"

	"github.com/jbreckmckye/psx-license-tool/internal/cdformat"
)

const LICENSE_SECTORS = 16

func ReadLicense(f *os.File) ([]cdformat.XAForm1Sector, error) {
	sectors := make([]cdformat.XAForm1Sector, LICENSE_SECTORS)

	const CHUNK_SIZE = cdformat.ISO_SECTOR_SIZE

	for i := range LICENSE_SECTORS {
		buffer := make([]byte, CHUNK_SIZE)
		offset := int64(CHUNK_SIZE * i)

		bytesRead, err := f.ReadAt(buffer, offset)
		if err == io.EOF {
			return sectors, errors.New("reached end of file too early, check the path is actually a disc BIN image")
		}
		if bytesRead < CHUNK_SIZE {
			return sectors, errors.New("error reading disc sector, read too few bytes")
		}

		sector, err := cdformat.ParseSectorXAForm1(buffer)
		if err != nil {
			return sectors, err
		}

		sectors[i] = sector
	}

	return sectors, nil
}

func GetLicenseText(license []cdformat.XAForm1Sector) []byte {
	sector := license[4]
	return sector.Data[:70]
}

func GetLicenseTMD(license []cdformat.XAForm1Sector) []byte {
	sectors := license[5:]
	tmd := []byte{}

	for _, sector := range sectors {
    bytes := sector.Data[:]
		tmd = slices.Concat(tmd, bytes)
	}

	// Trim trailing nulls
  lastByte := len(tmd) - 1
	for lastByte >= 0 {
    byte := tmd[lastByte]
    if byte == 0xFF {
			lastByte--
		} else {
			break
		}
	}

	return tmd
}