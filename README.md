# Go User Data

The library for cloud-init User Data

## Getting Started

Use `go get` to install the library

```shell
go get github.com/Aton-Kish/gouserdata
```

Import in your application

```go
import (
	userdata "github.com/Aton-Kish/gouserdata"
)
```

## Usage

This example shows how to make mime multipart user data.

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	userdata "github.com/Aton-Kish/gouserdata"
)

func main() {
	m, err := userdata.NewMultipart()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := os.ReadFile("cloud-config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	m.Append(userdata.NewPart(userdata.MediaTypeCloudConfig, cfg))

	j2, err := os.ReadFile("script.j2")
	if err != nil {
		log.Fatal(err)
	}
	m.Append(userdata.NewPart(userdata.MediaTypeJinja2, j2))

	hook, err := os.ReadFile("boothook.sh")
	if err != nil {
		log.Fatal(err)
	}
	m.Append(userdata.NewPart(userdata.MediaTypeCloudBoothook, hook))

	buf := new(bytes.Buffer)
	if err := m.Render(buf); err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())
}
```

## Development

### doc

```shell
: install godoc
go install golang.org/x/tools/cmd/godoc@latest

: run godoc server
godoc -http ":6060"

: uninstall godoc
rm $(go env GOPATH)/bin/godoc
```

### test

```shell
go test ./...
```

## License

This library is licensed under the MIT License, see [LICENSE](./LICENSE).
