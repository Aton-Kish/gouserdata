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

func TestHeader_Render(t *testing.T) {
	tests := []struct {
		name     string
		header   Header
		expected string
		err      error
	}{
		{
			name: "positive case: normal",
			header: func() Header {
				h := NewHeader()

				h.Set("Key1", "Key1-Value1")

				h.Set("Key2", "Key2-Value1")
				h.Add("Key2", "Key2-Value2")

				h.Set("Key3", "Key3-Value1")
				h.Add("Key3", "Key3-Value2")
				h.Add("Key3", "Key3-Value3")

				return *h
			}(),
			expected: "Key1: Key1-Value1\r\n" +
				"Key2: Key2-Value1\r\n" +
				"Key2: Key2-Value2\r\n" +
				"Key3: Key3-Value1\r\n" +
				"Key3: Key3-Value2\r\n" +
				"Key3: Key3-Value3\r\n",
		},
		{
			name: "positive case: sort keys",
			header: func() Header {
				h := NewHeader()

				h.Set("Key3", "Key3-Value3")
				h.Add("Key3", "Key3-Value2")
				h.Add("Key3", "Key3-Value1")

				h.Set("Key2", "Key2-Value2")
				h.Add("Key2", "Key2-Value1")

				h.Set("Key1", "Key1-Value1")

				return *h
			}(),
			expected: "Key1: Key1-Value1\r\n" +
				"Key2: Key2-Value1\r\n" +
				"Key2: Key2-Value2\r\n" +
				"Key3: Key3-Value1\r\n" +
				"Key3: Key3-Value2\r\n" +
				"Key3: Key3-Value3\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := tt.header.Render(buf)

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
