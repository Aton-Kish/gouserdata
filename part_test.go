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
