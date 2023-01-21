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
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMultipart(t *testing.T) {
	tests := []struct {
		name     string
		expected Multipart
		err      error
	}{
		{
			name: "positive case",
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=\"+Go+User+Data+Boundary==\""},
						"Mime-Version": {"1.0"},
					},
				},
				parts:    []Part{},
				boundary: "+Go+User+Data+Boundary==",
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewMultipart()

			if tt.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err)
			}
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
		expected Multipart
		err      error
	}{
		{
			name: "positive case: quoted",
			args: args{
				boundary: "+Go+User+Data+Boundary==",
			},
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=\"+Go+User+Data+Boundary==\""},
						"Mime-Version": {"1.0"},
					},
				},
				parts:    []Part{},
				boundary: "+Go+User+Data+Boundary==",
			},
			err: nil,
		},
		{
			name: "positive case: non quoted",
			args: args{
				boundary: "+Go+User+Data+Boundary++",
			},
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=+Go+User+Data+Boundary++"},
						"Mime-Version": {"1.0"},
					},
				},
				parts:    []Part{},
				boundary: "+Go+User+Data+Boundary++",
			},
			err: nil,
		},
		{
			name: "positive case: not ending with white space",
			args: args{
				boundary: " Go User Data Boundary==",
			},
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=\" Go User Data Boundary==\""},
						"Mime-Version": {"1.0"},
					},
				},
				parts:    []Part{},
				boundary: " Go User Data Boundary==",
			},
			err: nil,
		},
		{
			name: "positive case: valid characters",
			args: args{
				boundary: "0-9a-zA-Z'()+_,-./:=?",
			},
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=\"0-9a-zA-Z'()+_,-./:=?\""},
						"Mime-Version": {"1.0"},
					},
				},
				parts:    []Part{},
				boundary: "0-9a-zA-Z'()+_,-./:=?",
			},
			err: nil,
		},
		{
			name: "negative case: empty boundary",
			args: args{
				boundary: "",
			},
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=\"+Go+User+Data+Boundary==\""},
						"Mime-Version": {"1.0"},
					},
				},
				parts:    []Part{},
				boundary: "+Go+User+Data+Boundary==",
			},
			err: &Error{Op: "new", Err: ErrInvalidBoundary},
		},
		{
			name: "negative case: ending with white space",
			args: args{
				boundary: "+Go+User+Data+Boundary ",
			},
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=\"+Go+User+Data+Boundary==\""},
						"Mime-Version": {"1.0"},
					},
				},
				parts:    []Part{},
				boundary: "+Go+User+Data+Boundary==",
			},
			err: &Error{Op: "new", Err: ErrInvalidBoundary},
		},
		{
			name: "negative case: over 70 characters",
			args: args{
				boundary: "+Go+User+Data+Boundary==+Go+User+Data+Boundary==+Go+User+Data+Boundary==",
			},
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=\"+Go+User+Data+Boundary==\""},
						"Mime-Version": {"1.0"},
					},
				},
				parts:    []Part{},
				boundary: "+Go+User+Data+Boundary==",
			},
			err: &Error{Op: "new", Err: ErrInvalidBoundary},
		},
		{
			name: "negative case: includes invalid character",
			args: args{
				boundary: "!Go+User+Data+Boundary==",
			},
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=\"+Go+User+Data+Boundary==\""},
						"Mime-Version": {"1.0"},
					},
				},
				parts:    []Part{},
				boundary: "+Go+User+Data+Boundary==",
			},
			err: &Error{Op: "new", Err: ErrInvalidBoundary},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewMultipartWithBoundary(tt.args.boundary)

			if tt.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err)
			}
		})
	}
}

func TestMultipart_Append(t *testing.T) {
	type args struct {
		part Part
	}

	tests := []struct {
		name      string
		multipart Multipart
		args      []args
		expected  Multipart
	}{
		{
			name: "positive case: ascii only",
			multipart: func() Multipart {
				m, _ := NewMultipart()
				return m
			}(),
			args: []args{
				{
					part: NewPart(MediaTypeCloudConfig, []byte("#cloud-config\n"+"timezone: Europe/London")),
				},
				{
					part: NewPart(MediaTypeXShellscript, []byte("#!/bin/bash\n"+"echo 'Hello World'")),
				},
			},
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=\"+Go+User+Data+Boundary==\""},
						"Mime-Version": {"1.0"},
					},
				},
				parts: []Part{
					&part{
						header: &header{
							textproto.MIMEHeader{
								"Content-Transfer-Encoding": {"7bit"},
								"Content-Type":              {"text/cloud-config; charset=us-ascii"},
							},
						},
						body: []byte("#cloud-config\n" + "timezone: Europe/London"),
					},
					&part{
						header: &header{
							textproto.MIMEHeader{
								"Content-Transfer-Encoding": {"7bit"},
								"Content-Type":              {"text/x-shellscript; charset=us-ascii"},
							},
						},
						body: []byte("#!/bin/bash\n" + "echo 'Hello World'"),
					},
				},
				boundary: "+Go+User+Data+Boundary==",
			},
		},
		{
			name: "positive case: include utf-8",
			multipart: func() Multipart {
				m, _ := NewMultipart()
				return m
			}(),
			args: []args{
				{
					part: NewPart(MediaTypeCloudConfig, []byte("#cloud-config\n"+"timezone: Asia/Tokyo")),
				},
				{
					part: NewPart(MediaTypeXShellscript, []byte("#!/bin/bash\n"+"echo 'こんにちは世界'")),
				},
			},
			expected: &multipart{
				header: &header{
					textproto.MIMEHeader{
						"Content-Type": {"multipart/mixed; boundary=\"+Go+User+Data+Boundary==\""},
						"Mime-Version": {"1.0"},
					},
				},
				parts: []Part{
					&part{
						header: &header{
							textproto.MIMEHeader{
								"Content-Transfer-Encoding": {"7bit"},
								"Content-Type":              {"text/cloud-config; charset=us-ascii"},
							},
						},
						body: []byte("#cloud-config\n" + "timezone: Asia/Tokyo"),
					},
					&part{
						header: &header{
							textproto.MIMEHeader{
								"Content-Transfer-Encoding": {"base64"},
								"Content-Type":              {"text/x-shellscript; charset=utf-8"},
							},
						},
						body: []byte(
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
				tt.multipart.Append(args.part)
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
				m, _ := NewMultipart()

				m.Append(NewPart(MediaTypeCloudConfig, []byte("#cloud-config\n"+"timezone: Europe/London")))
				m.Append(NewPart(MediaTypeXShellscript, []byte("#!/bin/bash\n"+"echo 'Hello World'")))

				return m
			}(),
			expected: "Content-Type: multipart/mixed; boundary=\"+Go+User+Data+Boundary==\"\r\n" +
				"Mime-Version: 1.0\r\n" +
				"\r\n" +
				"--+Go+User+Data+Boundary==\r\n" +
				"Content-Transfer-Encoding: 7bit\r\n" +
				"Content-Type: text/cloud-config; charset=us-ascii\r\n" +
				"\r\n" +
				"#cloud-config\n" +
				"timezone: Europe/London\r\n" +
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
				m, _ := NewMultipart()

				m.Append(NewPart(MediaTypeCloudConfig, []byte("#cloud-config\n"+"timezone: Asia/Tokyo")))
				m.Append(NewPart(MediaTypeXShellscript, []byte("#!/bin/bash\n"+"echo 'こんにちは世界'")))

				return m
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
