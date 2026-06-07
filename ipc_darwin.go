package sysvipc

import (
	"structs"
	"unsafe"
)

// ipc_perm corresponds to struct ipc_perm.
// Size: 24 bytes, alignment: 4.
type ipc_perm struct {
	_    structs.HostLayout
	uid  uint32 // Owner's user ID
	gid  uint32 // Owner's group ID
	cuid uint32 // Creator's user ID
	cgid uint32 // Creator's group ID
	mode uint16 // Read/write permission
	_seq uint16 // RESERVED: internal use
	_key int32  // RESERVED: internal use
}

const (
	_ = uint(unsafe.Sizeof(ipc_perm{}) - 24)
	_ = uint(24 - unsafe.Sizeof(ipc_perm{}))
	_ = uint(unsafe.Alignof(ipc_perm{}) - 4)
	_ = uint(4 - unsafe.Alignof(ipc_perm{}))
)

func (p Perm) Key() int { return int(p.s._key) }
func (p Perm) Seq() int { return int(p.s._seq) }
func (p *Perm) SetMode(mode int) {
	p.s.mode = (p.s.mode &^ 0o777) | (uint16(mode) & 0o777)
}
