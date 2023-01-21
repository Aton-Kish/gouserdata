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
	"fmt"
	"io"
	"net/textproto"
	"sort"

	"golang.org/x/exp/maps"
)

type Header interface {
	Add(key, value string)
	Set(key, value string)
	Get(key string) string
	Values(key string) []string
	Del(key string)
	Renderer
}

type header struct {
	textproto.MIMEHeader
}

func NewHeader() Header {
	h := make(textproto.MIMEHeader)
	return &header{h}
}

func (h *header) Render(w io.Writer) error {
	keys := maps.Keys(h.MIMEHeader)
	sort.Strings(keys)

	for _, k := range keys {
		values := h.MIMEHeader[k]
		sort.Strings(values)

		for _, v := range values {
			if _, err := fmt.Fprintf(w, "%s: %s\r\n", k, v); err != nil {
				return err
			}
		}
	}

	return nil
}
