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

package unit_test

import (
	"fmt"
	"libclc"
	"libclc/unit"
	"os"
	//	"reflect"
	"testing"
	//	"unsafe"
)

func _temp() { fmt.Println() }

func newInitClc() libclc.Container {
	var b = make([]byte, 4096)
	stat, c := libclc.Using(b)
	if stat.IsError() {
		panic("TEST-BUG - libclc.Using")
	}

	unit.Init(c)
	return c
}

func TestUnit(t *testing.T) {
	var b = newInitClc()
	fmt.Printf("0x%0x\n", &b[0])
	stat, c := libclc.Using(b)
	if stat.IsError() {
		t.Fatalf("TEST-BUG - libclc.Using - stat:%s", stat)
	}
	var u = c.Unit(1)
	fmt.Printf("0x%0x\n", u)
	/*
		_b := make([]byte, 4096)
		c := libclc.Using(_b)
		c3 := unit.Container(c, 3)
		fmt.Printf("%p %p\n", c, c3)
		unit.Init(c3)
		unit.Dump(c3)
	*/
}

func TestInit(t *testing.T) {

	_b := make([]byte, 4096)
	for i := 0; i < 4096; i++ {
		_b[i] = 0xA
	}
	stat, c := libclc.Using(_b)
	if stat.IsError() {
		t.Fatalf("TEST-BUG - libclc.Using - stat:%s", stat)
	}

	unit.Init(c)

	// verify expectations
	// test c-meta
	// len: 0, default record order: 0-6
	var b = []byte(c)

	// len +0
	expect := byte(0)
	have := b[0]
	if have != expect {
		t.Fatalf("unit.Init() - len - expect:%d have:%d", expect, have)
	}
	// init r-meta values are in (0..6) inclusive. In order.
	for xof := 1; xof < 8; xof++ {
		expect = byte(xof - 1)
		have = b[xof]
		if have != expect {
			unit.DumpTo(c, os.Stderr)
			t.Fatalf("unit.Init() - r-meta[%d] - expect:%d have:%d", xof, expect, have)
		}
	}

	// All records must be zero-value
	for xof := uint8(0); xof < 7; xof++ {
		expect := uint64(0)
		have := *c.RecordPtr(xof)
		if have != expect {
			unit.DumpTo(c, os.Stderr)
			t.Fatalf("unit.Init() - r[%d] - expect:%016x have:%016x", xof, expect, have)
		}
	}
}

func TestLen(t *testing.T) {

	c := newInitClc()

	// init len
	have := unit.Len(c)
	expect := uint8(0)
	if have != expect {
		unit.DumpTo(c, os.Stderr)
		t.Fatalf("unit.Len() - expect:%d have:%d", expect, have)
	}

	// set len
	have = unit.SetLen(c, 3)
	expect = 3
	if have != expect {
		unit.DumpTo(c, os.Stderr)
		t.Fatalf("unit.SetLen() - expect:%d have:%d", expect, have)
	}

	have = unit.Len(c)
	if have != expect {
		unit.DumpTo(c, os.Stderr)
		t.Fatalf("unit.Len() (2) - expect:%d have:%d", expect, have)
	}

	// unit.Reset should not affect length. Check it.
	unit.Reset(c)
	have = unit.Len(c)
	if have != expect {
		unit.DumpTo(c, os.Stderr)
		t.Fatalf("unit.Len() (3) - expect:%d have:%d", expect, have)
	}

	// unit.Init should reset length. Check it.
	expect = 0
	unit.Init(c)
	have = unit.Len(c)
	if have != expect {
		unit.DumpTo(c, os.Stderr)
		t.Fatalf("unit.Len() (4) - expect:%d have:%d", expect, have)
	}
}

/// benchmarks ////////////////////////////////////////////////////////////////

func BenchmarkInit(b *testing.B) {
	c := newInitClc()
	for n := 0; n < b.N; n++ {
		unit.Init(c)
	}
}

func BenchmarkLen(b *testing.B) {
	c := newInitClc()
	for n := 0; n < b.N; n++ {
		unit.Len(c)
	}
}

func BenchmarkSetLen(b *testing.B) {
	c := newInitClc()
	for n := 0; n < b.N; n++ {
		unit.SetLen(c, 0)
	}
}
