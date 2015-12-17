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
	"libclc"
	"libclc/unit"
	"testing"
)

func newInitClc() libclc.Container {
	var b = make([]byte, 4096)
	var c = libclc.Using(b)

	unit.Init(c)
	return c
}
func TestLen(t *testing.T) {

	c := newInitClc()

	// init len
	have := unit.Len(c)
	expect := uint8(0)
	if have != expect {
		t.Fatalf("unit.Len() - expect:%d have:%d", expect, have)
	}

	// set len
	have = unit.SetLen(c, 3)
	expect = 3
	if have != expect {
		t.Fatalf("unit.SetLen() - expect:%d have:%d", expect, have)
	}

	have = unit.Len(c)
	if have != expect {
		t.Fatalf("unit.Len() (2) - expect:%d have:%d", expect, have)
	}

	// unit.Reset should not affect length. Check it.
	unit.Reset(c)
	have = unit.Len(c)
	if have != expect {
		t.Fatalf("unit.Len() (3) - expect:%d have:%d", expect, have)
	}

	// unit.Init should reset length. Check it.
	expect = 0
	unit.Init(c)
	have = unit.Len(c)
	if have != expect {
		t.Fatalf("unit.Len() (4) - expect:%d have:%d", expect, have)
	}
}
