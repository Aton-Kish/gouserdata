// Copyright (c) 2022 Aton-Kish
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// The library for cloud-init User Data
//
// Getting Started
//
// Use `go get` to install the library
//
// 	go get github.com/Aton-Kish/gouserdata
//
// Import in your application
//
// 	import (
// 		userdata "github.com/Aton-Kish/gouserdata"
// 	)
//
// Usage
//
// This example shows how to make mime multipart user data.
//
// 	package main
//
// 	import (
// 		"bytes"
// 		"fmt"
// 		"log"
// 		"os"
//
// 		userdata "github.com/Aton-Kish/gouserdata"
// 	)
//
// 	func main() {
// 		d := userdata.NewMultipart()
//
// 		cfg, err := os.ReadFile("cloud-config.yaml")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		d.AddPart(userdata.MediaTypeCloudConfig, cfg)
//
// 		j2, err := os.ReadFile("script.j2")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		d.AddPart(userdata.MediaTypeJinja2, j2)
//
// 		hook, err := os.ReadFile("boothook.sh")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		d.AddPart(userdata.MediaTypeCloudBoothook, hook)
//
// 		buf := new(bytes.Buffer)
// 		d.Render(buf)
//
// 		fmt.Println(buf.String())
// 	}
//
// License
//
// This library is licensed under the MIT License.
package userdata
