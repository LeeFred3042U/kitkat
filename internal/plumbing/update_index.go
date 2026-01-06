// MIT License

// Copyright (c) [2025] [Zeeshan Ahmad Alavi]

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package plumbing

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"io"
	"os"

	"github.com/LeeFred3042U/kitkat/internal/core"
)

const IndexPath = core.IndexPath

func UpdateIndex(entries []IndexEntry) error {
	var buf bytes.Buffer

	// ---- Header ----
	buf.Write([]byte("DIRC"))                                  // signature
	binary.Write(&buf, binary.BigEndian, uint32(2))            // version
	binary.Write(&buf, binary.BigEndian, uint32(len(entries))) // entry count

	// ---- Entries ----
	for _, e := range entries {
		if err := writeEntry(&buf, e); err != nil {
			return err
		}
	}

	// ---- Checksum ----
	sum := sha1.Sum(buf.Bytes())
	buf.Write(sum[:])

	return os.WriteFile(IndexPath, buf.Bytes(), 0644)
}

func writeEntry(w io.Writer, e IndexEntry) error {
	write := func(v interface{}) {
		_ = binary.Write(w, binary.BigEndian, v)
	}

	// ---- Stat fields ----
	write(e.CTimeSec)
	write(e.CTimeNSec)
	write(e.MTimeSec)
	write(e.MTimeNSec)
	write(e.Dev)
	write(e.Ino)
	write(e.Mode)
	write(e.UID)
	write(e.GID)
	write(e.Size)

	// ---- Object ID ----
	if _, err := w.Write(e.Hash[:]); err != nil {
		return err
	}

	// ---- Flags ----
	flags := uint16(len(e.Path)) & 0x0fff
	flags |= uint16(e.Stage&0x3) << 12
	write(flags)

	// ---- Path ----
	if _, err := w.Write([]byte(e.Path)); err != nil {
		return err
	}
	if _, err := w.Write([]byte{0}); err != nil {
		return err
	}

	padTo8(w)
	return nil
}

func padTo8(w io.Writer) {
	if b, ok := w.(*bytes.Buffer); ok {
		for b.Len()%8 != 0 {
			b.WriteByte(0)
		}
	}
}
