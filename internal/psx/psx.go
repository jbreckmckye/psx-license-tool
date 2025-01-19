package psx

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/jbreckmckye/psx-license-tool/internal/cdformat"
)

/**
 * The PSX license area format is documented on https://psx-spx.consoledev.net/cdromformat/#system-area-prior-to-volume-descriptors
 *
 * It is constructed of 16 CD-XA Form1 sectors. Data is therefore interleaved with ECC (error correction) data, but the ECC is
 * actually unused
 *
 * ECCs are ignored due to a bug with early versions of Sony's mastering tool. Early discs were burned with hardware like the
 * CDU-920 which didn't understand CD-XA and therefore didn't spot the problem... by the time Sony realised what had happened
 * there were too many games released with "broken" ECCs so they had to support it to be backwards compatible. Therefore the
 * PSX BIOS only reads the first 0x800 (2048) bytes of each sector.
 *
 * Sector overview
 *   Sector 0..3   - Zerofilled (Mode2/Form1, 4x800h bytes, plus ECC/EDC)
 *   Sector 4      - Licence String
 *   Sector 5..11  - Playstation Logo (3278h bytes) (remaining bytes FFh-filled)
 *   Sector 12..15 - Zerofilled (Mode2/Form2, 4x914h bytes, plus EDC)
 */

const LICENSE_SECTORS = 16

const TEXT_SECTOR = 4
const TMD_FIRST_SECTOR = 5

// EU/US: String is followed repeating 0x00 bytes
var EUR_STRING = [70]byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x4C, 0x69, 0x63, 0x65, 0x6E, 0x73, 0x65, 0x64, 0x20, 0x20, 0x62, 0x79, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x53, 0x6F, 0x6E, 0x79, 0x20, 0x43, 0x6F, 0x6D, 0x70, 0x75, 0x74, 0x65, 0x72, 0x20, 0x45, 0x6E, 0x74, 0x65, 0x72, 0x74, 0x61, 0x69, 0x6E, 0x6D, 0x65, 0x6E, 0x74, 0x20, 0x45, 0x75, 0x72, 0x6F, 0x20, 0x70, 0x65, 0x20, 0x20, 0x20}
var USA_STRING = [70]byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x4C, 0x69, 0x63, 0x65, 0x6E, 0x73, 0x65, 0x64, 0x20, 0x20, 0x62, 0x79, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x53, 0x6F, 0x6E, 0x79, 0x20, 0x43, 0x6F, 0x6D, 0x70, 0x75, 0x74, 0x65, 0x72, 0x20, 0x45, 0x6E, 0x74, 0x65, 0x72, 0x74, 0x61, 0x69, 0x6E, 0x6D, 0x65, 0x6E, 0x74, 0x20, 0x41, 0x6D, 0x65, 0x72, 0x20, 0x20, 0x69, 0x63, 0x61, 0x20}

// JP: String is followed by repeating 64 byte pattern of 62*0x30, 1*0x0A, 1*0x30 - this continues 31 times (1,984 bytes)
var JP_STRING = [65]byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x4C, 0x69, 0x63, 0x65, 0x6E, 0x73, 0x65, 0x64, 0x20, 0x20, 0x62, 0x79, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x53, 0x6F, 0x6E, 0x79, 0x20, 0x43, 0x6F, 0x6D, 0x70, 0x75, 0x74, 0x65, 0x72, 0x20, 0x45, 0x6E, 0x74, 0x65, 0x72, 0x74, 0x61, 0x69, 0x6E, 0x6D, 0x65, 0x6E, 0x74, 0x20, 0x49, 0x6E, 0x63, 0x2E, 0x0A}

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

func PatchLicense(f *os.File, license []cdformat.XAForm1Sector) error {
  for i, sector := range license {
    offset := int64(cdformat.ISO_SECTOR_SIZE * i)  
		bytes:= cdformat.SerialiseSectorXAForm1(sector)
    
		bytesWritten, err := f.WriteAt(bytes, offset)
		if err != nil {
			return err
		} 
		if bytesWritten < len(bytes) {
			return errors.New("wrote too few bytes")
		}
	}

	return nil
}

func GetLicenseText(license []cdformat.XAForm1Sector) [70]byte {
	sector := license[TEXT_SECTOR]
	return [70]byte(sector.Data[:70])
}

func PatchLicenseText(license []cdformat.XAForm1Sector, text []byte, japanese bool) {
	patchedSector := license[TEXT_SECTOR]

	padded := fmt.Sprintf("%-70s", text) // Text plus 70 spaces right-padding
	textBytes := []byte(padded)

	// Blank data section
	patchedSector.Data = [2048]byte{}
	var cursor = 0

	for cursor < 70 {
		patchedSector.Data[cursor] = textBytes[cursor]
		cursor++
	}

	for cursor < 2048 {
		if japanese {
			// Fill the rest of the sector with a repeating pattern of:
			// 62 times {0x30}; 1 single {0x0A}; 1 single {0x30}
			// This repeats 31 times... meaning it goes *one over* the sector size limit... we sort that further down...
			offset := cursor - 70
			patternPosition := offset % 64
			if patternPosition <= 61 {
				patchedSector.Data[cursor] = 0x30
			} else if patternPosition == 62 {
				patchedSector.Data[cursor] = 0x0A
			} else if patternPosition == 63 {
				patchedSector.Data[cursor] = 0x30
			}
		} else {
			patchedSector.Data[cursor] = 0x00
		}

		cursor++
	}

	if japanese {
		// Because the JP sector padding overruns into the EDC data... this format is a mess...
		patchedSector.EDC[0] = 0x30
	} // We don't bother un-setting this for non-JP discs as the EDC isn't used anyway

	license[TEXT_SECTOR] = patchedSector
}

func GetLicenseTMD(license []cdformat.XAForm1Sector) []byte {
	sectors := license[TMD_FIRST_SECTOR:]
	var tmd []byte

	for _, sector := range sectors {
		bytes := sector.Data[:]

		if !allZeroes(bytes) {
			tmd = slices.Concat(tmd, bytes)
		}
	}

	return tmd
}

// Returns a tuple: "writable" is whether the TMD *could* fit on the disc, "overBy" is how much over the assumed BIOS limit
func ValidateTMDSize(tmd []byte) (writable bool, overBy int) {
	assumedLimit := 7 * 2048 // Sectors 5..11, 2048 bytes per sector
	absoluteLimit := 11 * 2048
	size := len(tmd)

	writable = size < absoluteLimit
	if writable {
		overBy = size - assumedLimit
	} else {
		overBy = size - absoluteLimit
	}

	return writable, overBy
}

func PatchLicenseTMD(license []cdformat.XAForm1Sector, tmd []byte) {
	// Blank the data
	for i := TMD_FIRST_SECTOR; i < LICENSE_SECTORS; i++ {
		sector := license[i]

		if i < 12 {
			// TMD sectors are 0xFF filled
			sector.Data = [2048]byte{}
			for j := 0; j < 2048; j++ {
				sector.Data[j] = 0xFF
			}

		} else {
			// Final sectors are 0x00 filled
			sector.Data = [2048]byte{}
		}
	}

	// Copy in the TMD. Allow overrun into final sectors, although such TMDs probably won't work
	for pos, byte := range tmd {
		sectorN := TMD_FIRST_SECTOR + (pos / 2048)
		sectorPos := pos % 2048

		if sectorN > 15 {
			panic("Exceeded license data length")
		}

		license[sectorN].Data[sectorPos] = byte
	}
}

func allZeroes(input []byte) bool {
	for _, byte := range input {
		if byte != 0x00 {
			return false
		}
	}
	return true
}
