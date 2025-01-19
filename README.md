# PSX License Tool

This is a pair of utilities for dumping and patching license data on PlayStation 1 disc images.

Very new, still in beta, but local testing seems to work

## Installing

Please look to the `releases` for a binary matching your platform:

- Linux
- Windows x64
- MacOS ARM

Non-Linux binaries were cross-compiled and come with no guarantees (but the Windows ones seem to work on WINE)

## psxlicensedump

Given a BIN file from a `.CUE/.BIN` pair, writes out a LICENSE.TXT (license string) and LICENSE.TMD
file (PSX disc logo)

```shell
Usage: psxlicensedump [--output OUTPUT] BIN

Positional arguments:
  BIN                    path to a PSX disc image BIN

Options:
  --output OUTPUT        name for .TXT, .TMD output files [default: LICENSE]
  --help, -h             display this help and exit
```

## psxlicensepatch

Modifies the license data on a disc BIN archive (e.g from CUE/BIN pair)

- **Region** sets the region to apply to the license text section. If `text` is not set it will apply the default license text
- **Text** lets you set text. Bear in mind the data is truncated at 70 bytes
- **TMD** lets you specify a file for the TMD logo data. You'll get an error or warning if this looks too big.

```shell
Usage: psxlicensepatch [--region REGION] [--text TEXT] [--tmd TMD] BIN

Positional arguments:
  BIN                    path to a PSX disc image BIN

Options:
  --region REGION        Sets region string and / or padding. May be JP, EUR or US
  --text TEXT            Sets disc license text, overwriting region
  --tmd TMD              Path to TMD file to insert into license. Used for PSX boot logo
  --help, -h             display this help and exit
```

## Building

```shell
git clone git@github.com:jbreckmckye/psx-license-tool.git
cd psx-license-tool
make
```

## Notes

This was my first Golang project... don't expect any great shakes

I took most of my documentation from https://psx-spx.consoledev.net; go check it out.

Contributions welcome. Jimmy Breck-McKye 2025
