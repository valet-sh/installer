# Development

## Build for Linux

```bash
go build -o build/valet-sh-installer-linux-amd64 -ldflags="-s -w -X github.com/valet-sh/valet-sh-installer/cmd.version=<dev>" -v
``` 
