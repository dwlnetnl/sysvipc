package sysvipc

import (
	"structs"
	"unsafe"
)

// Additional control commands for msgctl, semctl, shmctl.
const (
	IPC_INFO = 3 // See ipcs.
)

// ipc_perm corresponds to struct ipc_perm.
// Size: 48 bytes, alignment: 8.
type ipc_perm struct {
	_    structs.HostLayout
	_key int32   // Key
	uid  uint32  // Owner's user ID
	gid  uint32  // Owner's group ID
	cuid uint32  // Creator's user ID
	cgid uint32  // Creator's group ID
	mode uint32  // Read/write permission
	_seq uint16  // Sequence number
	_    uint16  //
	_    [4]byte //
	_    uint64  // RESERVED
	_    uint64  // RESERVED
}

const (
	_ = uint(unsafe.Sizeof(ipc_perm{}) - 48)
	_ = uint(48 - unsafe.Sizeof(ipc_perm{}))
	_ = uint(unsafe.Alignof(ipc_perm{}) - 8)
	_ = uint(8 - unsafe.Alignof(ipc_perm{}))
)

func (p Perm) Key() int { return int(p.s._key) }
func (p Perm) Seq() int { return int(p.s._seq) }
func (p *Perm) SetMode(mode int) {
	p.s.mode = (p.s.mode &^ 0o777) | (uint32(mode) & 0o777)
}
