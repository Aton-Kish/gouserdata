# Go User Data

The library for cloud-init User Data

## Getting Started

Use `go get` to install the library

```shell
get github.com/Aton-Kish/gouserdata
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
	d := userdata.NewMultipart()

	cfg, err := os.ReadFile("cloud-config.yaml")
	if err != nil {
			log.Fatal(err)
	}
	d.AddPart(userdata.MediaTypeCloudConfig, cfg)

	j2, err := os.ReadFile("script.j2")
	if err != nil {
		log.Fatal(err)
	}
	d.AddPart(userdata.MediaTypeJinja2, j2)

	hook, err := os.ReadFile("boothook.sh")
	if err != nil {
		log.Fatal(err)
	}
	d.AddPart(userdata.MediaTypeCloudBoothook, hook)

	buf := new(bytes.Buffer)
	d.Render(buf)

	fmt.Println(buf.String())
}
```

## License

This library is licensed under the MIT License, see [LICENSE](./LICENSE).
