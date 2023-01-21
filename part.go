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
	"encoding/base64"
	"fmt"
	"io"
	"mime"

	"golang.org/x/exp/utf8string"
)

type Part interface {
	MediaType() MediaType
	SetBody(mediaType MediaType, body []byte)
	Renderer
}

type part struct {
	header    Header
	body      []byte
	mediaType MediaType
}

func NewPart() Part {
	h := NewHeader()
	return &part{header: h}
}

func (p *part) MediaType() MediaType {
	return p.mediaType
}

func (p *part) SetBody(mediaType MediaType, body []byte) {
	charset := "us-ascii"
	enc := "7bit"

	if !utf8string.NewString(string(body)).IsASCII() {
		charset = "utf-8"
		enc = "base64"
		body = []byte(base64.StdEncoding.EncodeToString(body))
	}

	typ := mime.FormatMediaType(string(mediaType), map[string]string{"charset": charset})

	p.header.Set("Content-Transfer-Encoding", enc)
	p.header.Set("Content-Type", typ)

	p.body = body
	p.mediaType = mediaType
}

func (p *part) Render(w io.Writer) error {
	if err := p.header.Render(w); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, "\r\n"); err != nil {
		return err
	}

	if _, err := w.Write(p.body); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, "\r\n"); err != nil {
		return err
	}

	return nil
}
