package sysvipc

import "golang.org/x/sys/unix"

// ipc_perm corresponds to struct ipc_perm.
type ipc_perm = unix.SysvIpcPerm

func (p *Perm) SetMode(mode int) {
	p.s.Mode = (p.s.Mode &^ 0o777) | (uint16(mode) & 0o777)
}
