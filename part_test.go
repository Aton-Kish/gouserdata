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

package multipart

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPart(t *testing.T) {
	type args struct {
		mediaType string
		body      []byte
	}

	tests := []struct {
		name     string
		args     args
		expected *Part
	}{
		{
			name: "positive case: ascii",
			args: args{
				mediaType: "text/x-shellscript",
				body:      []byte("#!/bin/bash\n" + "echo 'Hello World'"),
			},
			expected: &Part{
				Header: func() Header {
					h := NewHeader()

					h.Set("Content-Transfer-Encoding", "7bit")
					h.Set("Content-Type", "text/x-shellscript; charset=us-ascii")

					return *h
				}(),
				Body: []byte("#!/bin/bash\n" + "echo 'Hello World'"),
			},
		},
		{
			name: "positive case: utf-8",
			args: args{
				mediaType: "text/x-shellscript",
				body:      []byte("#!/bin/bash\n" + "echo 'こんにちは世界'"),
			},
			expected: &Part{
				Header: func() Header {
					h := NewHeader()

					h.Set("Content-Transfer-Encoding", "base64")
					h.Set("Content-Type", "text/x-shellscript; charset=utf-8")

					return *h
				}(),
				Body: []byte(base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\n" + "echo 'こんにちは世界'"))),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewPart(tt.args.mediaType, tt.args.body)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPart_Render(t *testing.T) {
	tests := []struct {
		name     string
		part     Part
		expected []byte
		err      error
	}{
		{
			name: "positive case: ascii",
			part: *NewPart("text/x-shellscript", []byte("#!/bin/bash\n"+"echo 'Hello World'")),
			expected: []byte(
				"Content-Transfer-Encoding: 7bit\r\n" +
					"Content-Type: text/x-shellscript; charset=us-ascii\r\n" +
					"\r\n" +
					"#!/bin/bash\n" + "echo 'Hello World'" +
					"\r\n",
			),
		},
		{
			name: "positive case: utf-8",
			part: *NewPart("text/x-shellscript", []byte("#!/bin/bash\n"+"echo 'こんにちは世界'")),
			expected: []byte(
				"Content-Transfer-Encoding: base64\r\n" +
					"Content-Type: text/x-shellscript; charset=utf-8\r\n" +
					"\r\n" +
					base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\n"+"echo 'こんにちは世界'")) +
					"\r\n",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.part.Render()

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
