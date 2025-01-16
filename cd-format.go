package main

import (
	"errors"
)

const LICENSE_SECTORS = 16
const ISO_SECTOR_SIZE = 2352

/*
 * PlayStation discs use the CD-ROM XA (Extended Architecture) layout where each sector is made
 * of 2352 bytes, including headers and metadata. How much metadata depends on the "form", Form 1
 * includes subheader and error correction data, Form 2 does not.
 *
 * The PSX license data is stored on the first 16 sectors of the disc image in XA-Form1 format.
 * When reading / writing the data we don't care about the ECC as it's wrong anyway 
 */

type XAForm1Sector struct {
  Sync    [12]byte     // Sync pattern (usually 00 FF FF FF FF FF FF FF FF FF FF 00)
	Addr    [3]byte      // Sector address
	Mode    byte         // Mode (usually 2 for Mode 2 Form 1/2 sectors)
	SubHead [8]byte      // Sub-header (00 00 08 00 00 00 08 00 for Form 1 data sectors)
	Data    [2048]byte   // Data (form 1)
	EDC     [4]byte      // Error-detection code (CRC32 of data area)
	ECC     [276]byte    // Error-correction code (uses Reed-Solomon ECC algorithm)
}

func ParseSectorXAForm1 (sector []byte) (XAForm1Sector, error) {
	if len(sector) != ISO_SECTOR_SIZE {
		return XAForm1Sector{}, errors.New("could not parse the disc sector (wrong length)")
	}

	sync := sector[0:12]
	addr := sector[12:15]
	mode := sector[15]
	subh := sector[16:24]
	data := sector[24:2072]
	erdc := sector[2072:2076]
	ercc := sector[2076:2352]

	return XAForm1Sector{
		Sync: [12]byte(sync),
		Addr: [3]byte(addr),
		Mode: mode,
		SubHead: [8]byte(subh),
		Data: [2048]byte(data),
		EDC: [4]byte(erdc),
		ECC: [276]byte(ercc),
	}, nil
}
