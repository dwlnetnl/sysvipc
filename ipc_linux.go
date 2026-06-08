package sysvipc

import (
	"golang.org/x/sys/unix"
)

// Additional commands for Msgctl, Semctl, Shmctl.
const (
	IPC_INFO = 3
)

// ipc_perm corresponds to struct ipc_perm.
type ipc_perm = unix.SysvIpcPerm

func (p Perm) Key() int { return int(p.s.Key) }
func (p Perm) Seq() int { return int(p.s.Seq) }
func (p *Perm) SetMode(mode int) {
	p.s.Mode = (p.s.Mode &^ 0o777) | (uint32(mode) & 0o777)
}
