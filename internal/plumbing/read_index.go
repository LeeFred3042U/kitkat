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
	"fmt"
	"io"
	"os"
)

func ReadIndex(path string) ([]IndexEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Verify checksum
	if len(data) < 20 {
		return nil, fmt.Errorf("index file too small")
	}

	content := data[:len(data)-20]
	expected := data[len(data)-20:]

	sum := sha1.Sum(content)
	if !bytes.Equal(sum[:], expected) {
		return nil, fmt.Errorf("index checksum mismatch")
	}

	r := bytes.NewReader(content)

	// ---- HEADER ----
	var sig [4]byte
	if _, err := io.ReadFull(r, sig[:]); err != nil {
		return nil, err
	}
	if string(sig[:]) != "DIRC" {
		return nil, fmt.Errorf("invalid index signature")
	}

	var version uint32
	if err := binary.Read(r, binary.BigEndian, &version); err != nil {
		return nil, err
	}
	if version != 2 {
		return nil, fmt.Errorf("unsupported index version: %d", version)
	}

	var entryCount uint32
	if err := binary.Read(r, binary.BigEndian, &entryCount); err != nil {
		return nil, err
	}

	entries := make([]IndexEntry, 0, entryCount)

	// ---- ENTRIES ----
	for i := uint32(0); i < entryCount; i++ {
		var e IndexEntry

		read := func(v interface{}) error {
			return binary.Read(r, binary.BigEndian, v)
		}

		read(&e.CTimeSec)
		read(&e.CTimeNSec)
		read(&e.MTimeSec)
		read(&e.MTimeNSec)
		read(&e.Dev)
		read(&e.Ino)
		read(&e.Mode)
		read(&e.UID)
		read(&e.GID)
		read(&e.Size)

		if _, err := io.ReadFull(r, e.Hash[:]); err != nil {
			return nil, err
		}

		var flags uint16
		read(&flags)

		// stage = bits 12–13
		e.Stage = uint8((flags >> 12) & 0x3)

		// read path
		var pathBuf []byte
		for {
			b, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			if b == 0 {
				break
			}
			pathBuf = append(pathBuf, b)
		}
		e.Path = string(pathBuf)

		// skip padding to 8-byte alignment
		for (r.Size()-int64(r.Len()))%8 != 0 {
			r.ReadByte()
		}

		entries = append(entries, e)
	}

	for r.Len() > 0 {
		if r.Len() < 8 {
			break
		}

		var extSig [4]byte
		io.ReadFull(r, extSig[:])

		var extSize uint32
		binary.Read(r, binary.BigEndian, &extSize)

		if extSize > uint32(r.Len()) {
			return nil, fmt.Errorf("corrupt index extension")
		}

		r.Seek(int64(extSize), io.SeekCurrent)
	}

	return entries, nil
}
