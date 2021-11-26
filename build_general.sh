mkdir -p "#temp"
mkdir -p "#temp/arm64"
env GOARCH=arm64 go build -o "#temp/arm64/meta"