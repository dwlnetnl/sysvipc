package sysvipc

import (
	"encoding/binary"
	"syscall"
	"time"
	_ "unsafe"

	"golang.org/x/sys/unix"
)

// ipc_perm corresponds to struct ipc_perm.
type ipc_perm = unix.SysvIpcPerm

func (p *Perm) SetMode(mode int) {
	p.s.Mode = (p.s.Mode &^ 0o777) | (uint16(mode) & 0o777)
}

//go:linkname runtime_syscall syscall.syscall
//go:linkname runtime_syscall6 syscall.syscall6

func runtime_syscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
func runtime_syscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

func bytesToUint64(b []byte) uint64 {
	return binary.NativeEndian.Uint64(b)
}

func bytesToTime(b []byte) time.Time {
	t := int64(binary.NativeEndian.Uint64(b))
	return time.Unix(t, 0)
}
