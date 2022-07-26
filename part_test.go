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

func TestPart_MediaType(t *testing.T) {
	tests := []struct {
		name     string
		part     Part
		expected MediaType
	}{
		{
			name: "positive case: empty",
			part: func() Part {
				p := NewPart()

				return *p
			}(),
			expected: MediaType(""),
		},
		{
			name: "positive case: text/cloud-boothook",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeCloudBoothook, []byte{})

				return *p
			}(),
			expected: MediaTypeCloudBoothook,
		},
		{
			name: "positive case: text/cloud-config",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeCloudConfig, []byte{})

				return *p
			}(),
			expected: MediaTypeCloudConfig,
		},
		{
			name: "positive case: text/cloud-config-archive",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeCloudConfigArchive, []byte{})

				return *p
			}(),
			expected: MediaTypeCloudConfigArchive,
		},
		{
			name: "positive case: text/cloud-config-jsonp",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeCloudConfigJsonp, []byte{})

				return *p
			}(),
			expected: MediaTypeCloudConfigJsonp,
		},
		{
			name: "positive case: text/jinja2",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeJinja2, []byte{})

				return *p
			}(),
			expected: MediaTypeJinja2,
		},
		{
			name: "positive case: text/part-handler",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypePartHandler, []byte{})

				return *p
			}(),
			expected: MediaTypePartHandler,
		},
		{
			name: "positive case: text/x-include-once-url",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeXIncludeOnceUrl, []byte{})

				return *p
			}(),
			expected: MediaTypeXIncludeOnceUrl,
		},
		{
			name: "positive case: text/x-include-url",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeXIncludeUrl, []byte{})

				return *p
			}(),
			expected: MediaTypeXIncludeUrl,
		},
		{
			name: "positive case: text/x-shellscript",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeXShellscript, []byte{})

				return *p
			}(),
			expected: MediaTypeXShellscript,
		},
		{
			name: "positive case: text/x-shellscript-per-boot",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeXShellscriptPerBoot, []byte{})

				return *p
			}(),
			expected: MediaTypeXShellscriptPerBoot,
		},
		{
			name: "positive case: text/x-shellscript-per-instance",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeXShellscriptPerInstance, []byte{})

				return *p
			}(),
			expected: MediaTypeXShellscriptPerInstance,
		},
		{
			name: "positive case: text/x-shellscript-per-once",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeXShellscriptPerOnce, []byte{})

				return *p
			}(),
			expected: MediaTypeXShellscriptPerOnce,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.part.MediaType()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPart_SetBody(t *testing.T) {
	type args struct {
		mediaType MediaType
		body      []byte
	}

	tests := []struct {
		name     string
		part     Part
		args     args
		expected Part
	}{
		{
			name: "positive case: ascii",
			part: *NewPart(),
			args: args{
				mediaType: MediaTypeXShellscript,
				body:      []byte("#!/bin/bash\n" + "echo 'Hello World'"),
			},
			expected: Part{
				Header: Header{
					textproto.MIMEHeader{
						"Content-Transfer-Encoding": {"7bit"},
						"Content-Type":              {"text/x-shellscript; charset=us-ascii"},
					},
				},
				Body:      []byte("#!/bin/bash\n" + "echo 'Hello World'"),
				mediaType: MediaTypeXShellscript,
			},
		},
		{
			name: "positive case: utf-8",
			part: *NewPart(),
			args: args{
				mediaType: MediaTypeXShellscript,
				body:      []byte("#!/bin/bash\n" + "echo 'こんにちは世界'"),
			},
			expected: Part{
				Header: Header{
					textproto.MIMEHeader{
						"Content-Transfer-Encoding": {"base64"},
						"Content-Type":              {"text/x-shellscript; charset=utf-8"},
					},
				},
				Body: []byte(
					// base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\n" + "echo 'こんにちは世界'")),
					"IyEvYmluL2Jhc2gKZWNobyAn44GT44KT44Gr44Gh44Gv5LiW55WMJw==",
				),
				mediaType: MediaTypeXShellscript,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.part.SetBody(tt.args.mediaType, tt.args.body)
			assert.Equal(t, tt.expected, tt.part)
		})
	}
}

func TestPart_Render(t *testing.T) {
	tests := []struct {
		name     string
		part     Part
		expected string
		err      error
	}{
		{
			name: "positive case: ascii",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeXShellscript, []byte("#!/bin/bash\n"+"echo 'Hello World'"))

				return *p
			}(),
			expected: "Content-Transfer-Encoding: 7bit\r\n" +
				"Content-Type: text/x-shellscript; charset=us-ascii\r\n" +
				"\r\n" +
				"#!/bin/bash\n" +
				"echo 'Hello World'\r\n",
		},
		{
			name: "positive case: utf-8",
			part: func() Part {
				p := NewPart()

				p.SetBody(MediaTypeXShellscript, []byte("#!/bin/bash\n"+"echo 'こんにちは世界'"))

				return *p
			}(),
			expected: "Content-Transfer-Encoding: base64\r\n" +
				"Content-Type: text/x-shellscript; charset=utf-8\r\n" +
				"\r\n" +
				// base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\n"+"echo 'こんにちは世界'")) +
				"IyEvYmluL2Jhc2gKZWNobyAn44GT44KT44Gr44Gh44Gv5LiW55WMJw==\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := tt.part.Render(buf)

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
