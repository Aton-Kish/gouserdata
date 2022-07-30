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
	"errors"
	"fmt"
	"io"
	"mime"
	"regexp"
)

const (
	defaultBoundary    = "+Go+User+Data+Boundary=="
	defaultMIMEVersion = "1.0"
)

var (
	boundaryRe = regexp.MustCompile(`^[0-9a-zA-Z'()+_,-./:=? ]{0,69}[0-9a-zA-Z'()+_,-./:=?]$`)
)

type Multipart struct {
	Header   Header
	Parts    []Part
	boundary string
}

func NewMultipart() *Multipart {
	h := NewHeader()
	h.Set("Mime-Version", defaultMIMEVersion)

	p := make([]Part, 0)

	m := &Multipart{Header: *h, Parts: p}
	m.SetBoundary(defaultBoundary)

	return m
}

func (m *Multipart) Boundary() string {
	return m.boundary
}

func (m *Multipart) SetBoundary(boundary string) error {
	if !boundaryRe.MatchString(boundary) {
		return errors.New("invalid boundary")
	}

	m.boundary = boundary

	typ := mime.FormatMediaType("multipart/mixed", map[string]string{"boundary": boundary})
	m.Header.Set("Content-Type", typ)

	return nil
}

func (m *Multipart) AddPart(mediaType MediaType, body []byte) {
	part := NewPart()
	part.SetBody(mediaType, body)
	m.Parts = append(m.Parts, *part)
}

func (m *Multipart) Render(w io.Writer) error {
	if err := m.Header.Render(w); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, "\r\n"); err != nil {
		return err
	}

	for _, part := range m.Parts {
		if _, err := fmt.Fprintf(w, "--%s\r\n", m.boundary); err != nil {
			return err
		}

		if err := part.Render(w); err != nil {
			return err
		}

		if _, err := fmt.Fprint(w, "\r\n"); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "--%s--\r\n", m.boundary); err != nil {
		return err
	}

	return nil
}
