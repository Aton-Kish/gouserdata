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

package userdata

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMultipart(t *testing.T) {
	tests := []struct {
		name     string
		expected *Multipart
	}{
		{
			name: "positive case",
			expected: &Multipart{
				Header: func() Header {
					h := NewHeader()

					h.Set("Content-Type", `multipart/mixed; boundary="+Go+User+Data+Boundary=="`)
					h.Set("Mime-Version", "1.0")

					return *h
				}(),
				Parts:    []Part{},
				boundary: "+Go+User+Data+Boundary==",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewMultipart()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestNewMultipartWithBoundary(t *testing.T) {
	type args struct {
		boundary string
	}

	tests := []struct {
		name     string
		args     args
		expected *Multipart
	}{
		{
			name: "positive case: quoted",
			args: args{
				boundary: "+Go+User+Data+Boundary==",
			},
			expected: &Multipart{
				Header: func() Header {
					h := NewHeader()

					h.Set("Content-Type", "multipart/mixed; boundary=\"+Go+User+Data+Boundary==\"")
					h.Set("Mime-Version", "1.0")

					return *h
				}(),
				Parts:    []Part{},
				boundary: "+Go+User+Data+Boundary==",
			},
		},
		{
			name: "positive case: non quoted",
			args: args{
				boundary: "+Go+User+Data+Boundary++",
			},
			expected: &Multipart{
				Header: func() Header {
					h := NewHeader()

					h.Set("Content-Type", "multipart/mixed; boundary=+Go+User+Data+Boundary++")
					h.Set("Mime-Version", "1.0")

					return *h
				}(),
				Parts:    []Part{},
				boundary: "+Go+User+Data+Boundary++",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewMultipartWithBoundary(tt.args.boundary)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestMultipart_AddPart(t *testing.T) {
	type args struct {
		mediaType MediaType
		body      []byte
	}

	tests := []struct {
		name      string
		multipart Multipart
		args      []args
		expected  Multipart
	}{
		{
			name:      "positive case: ascii only",
			multipart: *NewMultipart(),
			args: []args{
				{
					mediaType: CloudConfig,
					body:      []byte("#cloud-config\n" + "timezone: America/Virgin"),
				},
				{
					mediaType: XShellscript,
					body:      []byte("#!/bin/bash\n" + "echo 'Hello World'"),
				},
			},
			expected: Multipart{
				Header: func() Header {
					h := NewHeader()

					h.Set("Content-Type", "multipart/mixed; boundary=\"+Go+User+Data+Boundary==\"")
					h.Set("Mime-Version", "1.0")

					return *h
				}(),
				Parts: []Part{
					{
						Header: func() Header {
							h := NewHeader()

							h.Set("Content-Transfer-Encoding", "7bit")
							h.Set("Content-Type", "text/cloud-config; charset=us-ascii")

							return *h
						}(),
						Body: []byte("#cloud-config\n" + "timezone: America/Virgin"),
					},
					{
						Header: func() Header {
							h := NewHeader()

							h.Set("Content-Transfer-Encoding", "7bit")
							h.Set("Content-Type", "text/x-shellscript; charset=us-ascii")

							return *h
						}(),
						Body: []byte("#!/bin/bash\n" + "echo 'Hello World'"),
					},
				},
				boundary: "+Go+User+Data+Boundary==",
			},
		},
		{
			name:      "positive case: include utf-8",
			multipart: *NewMultipart(),
			args: []args{
				{
					mediaType: CloudConfig,
					body:      []byte("#cloud-config\n" + "timezone: Asia/Tokyo"),
				},
				{
					mediaType: XShellscript,
					body:      []byte("#!/bin/bash\n" + "echo 'こんにちは世界'"),
				},
			},
			expected: Multipart{
				Header: func() Header {
					h := NewHeader()

					h.Set("Content-Type", "multipart/mixed; boundary=\"+Go+User+Data+Boundary==\"")
					h.Set("Mime-Version", "1.0")

					return *h
				}(),
				Parts: []Part{
					{
						Header: func() Header {
							h := NewHeader()

							h.Set("Content-Transfer-Encoding", "7bit")
							h.Set("Content-Type", "text/cloud-config; charset=us-ascii")

							return *h
						}(),
						Body: []byte("#cloud-config\n" + "timezone: Asia/Tokyo"),
					},
					{
						Header: func() Header {
							h := NewHeader()

							h.Set("Content-Transfer-Encoding", "base64")
							h.Set("Content-Type", "text/x-shellscript; charset=utf-8")

							return *h
						}(),
						Body: []byte(
							// base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\n" + "echo 'こんにちは世界'")),
							"IyEvYmluL2Jhc2gKZWNobyAn44GT44KT44Gr44Gh44Gv5LiW55WMJw==",
						),
					},
				},
				boundary: "+Go+User+Data+Boundary==",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, args := range tt.args {
				tt.multipart.AddPart(args.mediaType, args.body)
			}
			assert.Equal(t, tt.expected, tt.multipart)
		})
	}
}

func TestMultipart_Render(t *testing.T) {
	tests := []struct {
		name      string
		multipart Multipart
		expected  string
		err       error
	}{
		{
			name: "positive case: ascii only",
			multipart: func() Multipart {
				d := NewMultipart()

				d.AddPart(CloudConfig, []byte("#cloud-config\n"+"timezone: America/Virgin"))
				d.AddPart(XShellscript, []byte("#!/bin/bash\n"+"echo 'Hello World'"))

				return *d
			}(),
			expected: "Content-Type: multipart/mixed; boundary=\"+Go+User+Data+Boundary==\"\r\n" +
				"Mime-Version: 1.0\r\n" +
				"\r\n" +
				"--+Go+User+Data+Boundary==\r\n" +
				"Content-Transfer-Encoding: 7bit\r\n" +
				"Content-Type: text/cloud-config; charset=us-ascii\r\n" +
				"\r\n" +
				"#cloud-config\n" +
				"timezone: America/Virgin\r\n" +
				"\r\n" +
				"--+Go+User+Data+Boundary==\r\n" +
				"Content-Transfer-Encoding: 7bit\r\n" +
				"Content-Type: text/x-shellscript; charset=us-ascii\r\n" +
				"\r\n" +
				"#!/bin/bash\n" +
				"echo 'Hello World'\r\n" +
				"\r\n" +
				"--+Go+User+Data+Boundary==--\r\n",
		},
		{
			name: "positive case: include utf-8",
			multipart: func() Multipart {
				d := NewMultipart()

				d.AddPart(CloudConfig, []byte("#cloud-config\n"+"timezone: Asia/Tokyo"))
				d.AddPart(XShellscript, []byte("#!/bin/bash\n"+"echo 'こんにちは世界'"))

				return *d
			}(),
			expected: "Content-Type: multipart/mixed; boundary=\"+Go+User+Data+Boundary==\"\r\n" +
				"Mime-Version: 1.0\r\n" +
				"\r\n" +
				"--+Go+User+Data+Boundary==\r\n" +
				"Content-Transfer-Encoding: 7bit\r\n" +
				"Content-Type: text/cloud-config; charset=us-ascii\r\n" +
				"\r\n" +
				"#cloud-config\n" +
				"timezone: Asia/Tokyo\r\n" +
				"\r\n" +
				"--+Go+User+Data+Boundary==\r\n" +
				"Content-Transfer-Encoding: base64\r\n" +
				"Content-Type: text/x-shellscript; charset=utf-8\r\n" +
				"\r\n" +
				// base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\n"+"echo 'こんにちは世界'")) +
				"IyEvYmluL2Jhc2gKZWNobyAn44GT44KT44Gr44Gh44Gv5LiW55WMJw==\r\n" +
				"\r\n" +
				"--+Go+User+Data+Boundary==--\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := tt.multipart.Render(buf)

			if tt.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, buf.String())
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err)
			}
		})
	}
}
