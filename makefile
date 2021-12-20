build-linux64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o autoCheckIn ./src
build-win64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o autoCheckIn.exe ./src
build-macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o autoCheckIn ./src