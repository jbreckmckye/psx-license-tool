# PSX License Tool

This is a pair of utilities for dumping and patching license data on PlayStation 1 disc images.

## Installing

Please look to the `releases` for a binary matching your platform.

Windows executables have been cross-compiled with Linux and come with no guarantees

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

## Building on your platform

(Unix)

```shell
git clone git@github.com:jbreckmckye/psx-license-tool.git
cd psx-license-tool
go build -o build ./cmd/psxlicensedump
go build -o build ./cmd/psxlicensepatch
```
