/* !!! FRIEND !!! */
//
//    Copyright © 2015 Joubin Muhammad Houshyar. All rights reserved.
//
//    This file is part of libclc-go.
//
//    Foobar is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    Foobar is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with libclc-go.  If not, see <http://www.gnu.org/licenses/>.

// [STUB]
// Unit cache line container is provides functionality for unit containers
// used as either stand-alone or element(s) of a segmented-container.
//
// NOTES:
// 1 - All functions in this package with arg-0 'libclc.Container'  assume
// the value has been obtained via functionality of libclc and *do not*
// check for correct length, alignment, nil, etc.
package unit

import (
	"fmt"
	"io"
	"libclc"
	"os"
)

// @REF-CHECK
// Resets the container iteration order.
// Does NOT modify c-meta.
func Reset(u libclc.Container) {
	for b := uint8(0); b < 7; b++ {
		*u.BytePtr(b + 1) = b
	}
	return
}

// @REF-CHECK
// Sets container to initial state: default order, 0-len,
// and zero'd data records.
func Init(u libclc.Container) {
	Reset(u)
	*u.BytePtr0() = 0
	for r := uint8(0); r < 7; r++ {
		*u.RecordPtr(r) = 0
	}
	return
}

/// STORE /////////////////////////////////////////////////////////////////////

// Returns unit content length in range (0, 7) inclusive.
// See general package notes regarding argument asserts.
func Len(u libclc.Container) uint8 {
	return *u.BytePtr0() & libclc.M_clen
}

// Sets the unit length, with arg in range (0, 7) inclusive.
// Returns unit content length in range (0, 7) inclusive.
// See general package notes regarding argument asserts.
func SetLen(u libclc.Container, length uint8) uint8 {
	*u.BytePtr0() &= libclc.M_clen_inv
	*u.BytePtr0() |= length
	return *u.BytePtr0()
}

// non-exclusive put, adds rec to the next available slot per
// instrinsic unit record order.
func Put(u libclc.Container, rec uint64) (s libclc.Stat, rmeta byte) {
	if Len(u) == 7 {
		return libclc.Full, byte(0)
	}
	rmeta = *u.BytePtr(6)
	*u.Uint64Ptr((1 + (rmeta & libclc.M_ridx)) << 3) = rec
	return libclc.ErrNotImpl, byte(0)
}

func Dump(u libclc.Container) error {
	return DumpTo(u, os.Stderr)
}

// REVU match libclc's emit pattern
func DumpTo(u libclc.Container, w io.Writer) error {
	for xof := uint8(0); xof < 8; xof++ {
		_, e := fmt.Fprintf(w, "+%02d | %016x\n", xof<<3, *u.Uint64Ptr(xof << 3))
		if e != nil {
			return e
		}
	}
	return nil
}
