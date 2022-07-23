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
	"fmt"
	"mime"
)

const (
	defaultBoundary    = "+Go+User+Data+Boundary=="
	defaultMIMEVersion = "1.0"
)

type Multipart struct {
	Header   Header
	Parts    []Part
	boundary string
}

func NewMultipart() *Multipart {
	return NewMultipartWithBoundary(defaultBoundary)
}

func NewMultipartWithBoundary(boundary string) *Multipart {
	typ := mime.FormatMediaType("multipart/mixed", map[string]string{"boundary": boundary})

	h := NewHeader()
	h.Set("Content-Type", typ)
	h.Set("Mime-Version", defaultMIMEVersion)

	p := make([]Part, 0)

	return &Multipart{Header: *h, Parts: p, boundary: boundary}
}

func (m *Multipart) AddPart(mediaType MediaType, body []byte) {
	part := NewPart(mediaType, body)
	m.Parts = append(m.Parts, *part)
}

func (m *Multipart) Render() ([]byte, error) {
	buf := new(bytes.Buffer)

	h, err := m.Header.Render()
	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(h); err != nil {
		return nil, err
	}

	if _, err := buf.WriteString("\r\n"); err != nil {
		return nil, err
	}

	for _, part := range m.Parts {
		if _, err := buf.WriteString(fmt.Sprintf("--%s\r\n", m.boundary)); err != nil {
			return nil, err
		}

		p, err := part.Render()
		if err != nil {
			return nil, err
		}

		if _, err := buf.Write(p); err != nil {
			return nil, err
		}

		if _, err := buf.WriteString("\r\n"); err != nil {
			return nil, err
		}
	}

	if _, err := buf.WriteString(fmt.Sprintf("--%s--\r\n", m.boundary)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
