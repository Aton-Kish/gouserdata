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

package userdata_test

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	userdata "github.com/Aton-Kish/gouserdata"
)

func ExampleMultipart_Render() {
	d := userdata.NewMultipart()

	cfg := []byte(`#cloud-config
timezone: Europe/London`)
	d.AddPart(userdata.MediaTypeCloudConfig, cfg)

	scr := []byte(`#!/bin/bash
echo 'Hello World'`)
	d.AddPart(userdata.MediaTypeXShellscript, scr)

	buf := new(bytes.Buffer)
	d.Render(buf)

	output := buf.String()
	output = strings.ReplaceAll(output, "\r\n", "\n") // for testing
	fmt.Println(output)
	// Output:
	// Content-Type: multipart/mixed; boundary="+Go+User+Data+Boundary=="
	// Mime-Version: 1.0
	//
	// --+Go+User+Data+Boundary==
	// Content-Transfer-Encoding: 7bit
	// Content-Type: text/cloud-config; charset=us-ascii
	//
	// #cloud-config
	// timezone: Europe/London
	//
	// --+Go+User+Data+Boundary==
	// Content-Transfer-Encoding: 7bit
	// Content-Type: text/x-shellscript; charset=us-ascii
	//
	// #!/bin/bash
	// echo 'Hello World'
	//
	// --+Go+User+Data+Boundary==--
}

func ExampleMultipart_Render_includesUtf8() {
	d := userdata.NewMultipart()

	cfg := []byte(`#cloud-config
timezone: Asia/Tokyo`)
	d.AddPart(userdata.MediaTypeCloudConfig, cfg)

	scr := []byte(`#!/bin/bash
echo 'こんにちは世界'`)
	d.AddPart(userdata.MediaTypeXShellscript, scr)

	buf := new(bytes.Buffer)
	d.Render(buf)

	output := buf.String()
	output = strings.ReplaceAll(output, "\r\n", "\n") // for testing
	fmt.Println(output)
	// Output:
	// Content-Type: multipart/mixed; boundary="+Go+User+Data+Boundary=="
	// Mime-Version: 1.0
	//
	// --+Go+User+Data+Boundary==
	// Content-Transfer-Encoding: 7bit
	// Content-Type: text/cloud-config; charset=us-ascii
	//
	// #cloud-config
	// timezone: Asia/Tokyo
	//
	// --+Go+User+Data+Boundary==
	// Content-Transfer-Encoding: base64
	// Content-Type: text/x-shellscript; charset=utf-8
	//
	// IyEvYmluL2Jhc2gKZWNobyAn44GT44KT44Gr44Gh44Gv5LiW55WMJw==
	//
	// --+Go+User+Data+Boundary==--
}

func ExampleMultipart_SetBoundary() {
	d := userdata.NewMultipart()

	if err := d.SetBoundary("+Custom+User+Data+Boundary+"); err != nil {
		log.Fatal(err)
	}

	scr := []byte(`#!/bin/bash
echo 'Hello World'`)
	d.AddPart(userdata.MediaTypeXShellscript, scr)

	buf := new(bytes.Buffer)
	d.Render(buf)

	output := buf.String()
	output = strings.ReplaceAll(output, "\r\n", "\n") // for testing
	fmt.Println(output)
	// Output:
	// Content-Type: multipart/mixed; boundary=+Custom+User+Data+Boundary+
	// Mime-Version: 1.0
	//
	// --+Custom+User+Data+Boundary+
	// Content-Transfer-Encoding: 7bit
	// Content-Type: text/x-shellscript; charset=us-ascii
	//
	// #!/bin/bash
	// echo 'Hello World'
	//
	// --+Custom+User+Data+Boundary+--
}
