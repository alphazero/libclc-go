// !!! FRIEND !!!
package libclc

import (
	"reflect"
	"unsafe"
)

// used as base reference and to permit functions on this type for
// macro-like one liner functions. We could use uintptr directly, but
// the (moving) GC interactions are not entirely clear.
//
// Note: this type is exported ("public") since we (a) require it in
// other (nested) packages, such as libclc/unit, and (b) Go does not
// permit selective access (via e.g. "protected" or "friend").
type Container []byte

// Verfies the buffer in terms of alignment, length, etc, and
// recasts as Container. For integrity checks, we rely on this function
// providing the 'container' value object, so refrain from directly
// casting []byte objects to libclc.Contaienr!
//
// panics.
func Using(b []byte) Container {
	// TODO asserts here once
	// TODO check if nil
	// TODO check if not aligned 0x40

	return Container(b)
}

/* ------------------------------------------------------------------------- */
/* convenience pointer ops --- */

// In a perfect world these will all be inlined by the compiler
// NOTE: all functions of type Unit assume verified pointer mapped via the
// libclc.Container type cast method.

func (p Container) Pointer() uintptr {
	return (*reflect.SliceHeader)(unsafe.Pointer(&p)).Data
}

// typically used to get container cmeta
func (p Container) BytePtr0() *byte {
	return (*byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&p)).Data))
}

// typically used to get container record rmeta
// note that function relies on 0x40 alignment of the Container for correct op.
func (p Container) BytePtr(xof uint8) *byte {
	return (*byte)(unsafe.Pointer(((*reflect.SliceHeader)(unsafe.Pointer(&p)).Data) | uintptr(xof)))
}

// typically used to get container state record R0
func (p Container) Uint64Ptr0() *uint64 {
	return (*uint64)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&p)).Data))
}

// typically used to get container record
// note that function relies on 0x40 alignment of the Container for correct op.
func (p Container) Uint64Ptr(xof uint8) *uint64 {
	return (*uint64)(unsafe.Pointer(((*reflect.SliceHeader)(unsafe.Pointer(&p)).Data) | uintptr(xof)))
}

/* ------------------------------------------------------------------------- */
/* meta-byte masks ----------- */

// REVU: keep conformant to c-11.RI.

// Use distinct consts for R/C-META (even though ++likely they will never
// diverge) just to allow for that possibility.

// REVU simply can not stand constants with camel case and Go requires init Cap
// letter for public objects, so compromise here is uniform treatment of constants
// in this codebase: M_ for mask, and then readable constant names.

const (
	/* r-meta word masks */
	M_ridx       byte = 0x07
	M_ridx_inv   byte = ^M_ridx
	M_rext       byte = 0x38
	M_rext_inv   byte = ^M_rext
	M_rdirty     byte = 0x40
	M_rdirty_inv byte = ^M_rdirty
	M_rlock      byte = 0x40
	M_rlock_inv  byte = ^M_rlock

	/* c-meta word masks */
	M_clen       byte = 0x07
	M_clen_inv   byte = ^M_clen
	M_cext       byte = 0x38
	M_cext_inv   byte = ^M_cext
	M_cdirty     byte = 0x40
	M_cdirty_inv byte = ^M_cdirty
	M_clock      byte = 0x40
	M_clock_inv  byte = ^M_clock
)

/* ------------------------------------------------------------------------- */
/* libclc status codes ------- */

// REVU: keep conformant to c-11.RI.

type Stat int

func (s Stat) IsError() bool { return s < 0 }

const (
	/* non-error stats */
	ok        = Stat(0)
	Full      = Stat(1)
	Empty     = Stat(2)
	Removed   = Stat(3)
	NotFound  = Stat(4)
	Duplicate = Stat(5)

	/* error stats: user | system error */
	ErrState     = Stat(-1)
	ErrAlignment = Stat(-2)
	ErrPointer   = Stat(-3)
	ErrArg       = Stat(-4)
	ErrSelector  = Stat(-5)
	ErrIndex     = Stat(-6)
	ErrRecord    = Stat(-7)
	ErrNotImpl   = Stat(-255)
)

// TODO: debug String() methd for Stat type.

/* ------------------------------------------------------------------------- */
/* systolic shift ops -------- */

type clc_rshift struct {
	m0, m1, m2, shft uint64
}

// REVU: keep conformant to c-11.RI.

var rmask_up = [8]clc_rshift{
	clc_rshift{0xffffffffffff00ff, 0x0000000000000000, 0x000000000000ff00, 0}, // nop
	clc_rshift{0xffffffffff0000ff, 0x000000000000ff00, 0x0000000000ff0000, 8},
	clc_rshift{0xffffffff000000ff, 0x0000000000ffff00, 0x00000000ff000000, 16},
	clc_rshift{0xffffff00000000ff, 0x00000000ffffff00, 0x000000ff00000000, 24},
	clc_rshift{0xffff0000000000ff, 0x000000ffffffff00, 0x0000ff0000000000, 32},
	clc_rshift{0xff000000000000ff, 0x0000ffffffffff00, 0x00ff000000000000, 40},
	clc_rshift{0x00000000000000ff, 0x00ffffffffffff00, 0xff00000000000000, 48},
	clc_rshift{0, 0, 0, 0},
}
var rmask_dn = [8]clc_rshift{
	clc_rshift{0x00000000000000ff, 0xffffffffffff0000, 0x000000000000ff00, 48},
	clc_rshift{0x000000000000ffff, 0xffffffffff000000, 0x0000000000ff0000, 40},
	clc_rshift{0x0000000000ffffff, 0xffffffff00000000, 0x00000000ff000000, 32},
	clc_rshift{0x00000000ffffffff, 0xffffff0000000000, 0x000000ff00000000, 24},
	clc_rshift{0x000000ffffffffff, 0xffff000000000000, 0x0000ff0000000000, 16},
	clc_rshift{0x0000ffffffffffff, 0xff00000000000000, 0x00ff000000000000, 8},
	clc_rshift{0x00ffffffffffffff, 0x0000000000000000, 0xff00000000000000, 0}, // nop
	clc_rshift{0, 0, 0, 0},
}
var rmask_r6_to = [8]clc_rshift{
	clc_rshift{0x00000000000000ff, 0x00ffffffffffff00, 0xff00000000000000, 48},
	clc_rshift{0x000000000000ffff, 0x00ffffffffff0000, 0xff00000000000000, 40},
	clc_rshift{0x0000000000ffffff, 0x00ffffffff000000, 0xff00000000000000, 32},
	clc_rshift{0x00000000ffffffff, 0x00ffffff00000000, 0xff00000000000000, 24},
	clc_rshift{0x000000ffffffffff, 0x00ffff0000000000, 0xff00000000000000, 16},
	clc_rshift{0x0000ffffffffffff, 0x00ff000000000000, 0xff00000000000000, 8},
	clc_rshift{0x00ffffffffffffff, 0x0000000000000000, 0xff00000000000000, 0}, // nop
	clc_rshift{0, 0, 0, 0},
}

// REVU: TODO: add boundary condition up/dn with hardcoded masks
// TODO: add basic up/dn/r_to systolic ops.
