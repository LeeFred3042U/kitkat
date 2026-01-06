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

type IndexEntry struct {
	// ---- Logical (Git semantics) ----
	Path  string   // file path relative to repo root
	Hash  [20]byte // SHA-1 object ID (blob)
	Mode  uint32   // 0100644, 0100755, 0120000
	Stage uint8    // 0–3 (merge stages)

	// ---- Cached stat info (filesystem) ----
	CTimeSec  uint32
	CTimeNSec uint32
	MTimeSec  uint32
	MTimeNSec uint32
	Dev       uint32
	Ino       uint32
	UID       uint32
	GID       uint32
	Size      uint32
}
