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
	"bytes"
	"encoding/base64"
	"mime"

	"golang.org/x/exp/utf8string"
)

type Part struct {
	Header Header
	Body   []byte
}

func NewPart(mediaType string, body []byte) *Part {
	charset := "us-ascii"
	enc := "7bit"

	if !utf8string.NewString(string(body)).IsASCII() {
		charset = "utf-8"
		enc = "base64"
		body = []byte(base64.StdEncoding.EncodeToString(body))
	}

	typ := mime.FormatMediaType(mediaType, map[string]string{"charset": charset})

	h := NewHeader()
	h.Set("Content-Transfer-Encoding", enc)
	h.Set("Content-Type", typ)

	return &Part{Header: *h, Body: body}
}

func (p *Part) Render() ([]byte, error) {
	buf := new(bytes.Buffer)

	h, err := p.Header.Render()
	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(h); err != nil {
		return nil, err
	}

	if _, err := buf.WriteString("\r\n"); err != nil {
		return nil, err
	}

	if _, err := buf.Write(p.Body); err != nil {
		return nil, err
	}

	if _, err := buf.WriteString("\r\n"); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
