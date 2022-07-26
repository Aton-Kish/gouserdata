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

	fmt.Println(buf.String())
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

	fmt.Println(buf.String())
}

func ExampleMultipart_SetBoundary() {
	d := userdata.NewMultipart()

	d.SetBoundary("+Custom+User+Data+Boundary+")

	scr := []byte(`#!/bin/bash
echo 'Hello World'`)
	d.AddPart(userdata.MediaTypeXShellscript, scr)

	buf := new(bytes.Buffer)
	d.Render(buf)

	fmt.Println(strings.ReplaceAll(buf.String(), "\r\n", "\n"))
}
