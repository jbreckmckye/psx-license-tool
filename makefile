all: clean build pack

clean:
	rm -rf build

build:
	GOOS=linux   GOARCH=amd64 go build -o build/linux/ ./cmd/psxlicensedump
	GOOS=linux   GOARCH=amd64 go build -o build/linux/ ./cmd/psxlicensepatch
	GOOS=windows GOARCH=amd64 go build -o build/windows/ ./cmd/psxlicensedump
	GOOS=windows GOARCH=amd64 go build -o build/windows/ ./cmd/psxlicensepatch
	GOOS=darwin  GOARCH=arm64 go build -o build/macos/ ./cmd/psxlicensedump
	GOOS=darwin  GOARCH=arm64 go build -o build/macos/ ./cmd/psxlicensepatch

pack:
	(cd build/windows && zip -9 -y -r psxlicensetool.zip .)
	(cd build/linux   && zip -9 -y -r psxlicensetool.zip .)
	(cd build/macos   && zip -9 -y -r psxlicensetool.zip .)
