// Package sysvipc provides access to System V inter-process communication
// mechanisms.
package sysvipc

import "golang.org/x/sys/unix"

// Mode bits.
const (
	IPC_CREAT  = unix.IPC_CREAT  // Create entry if key does not exist
	IPC_EXCL   = unix.IPC_EXCL   // Fail if key exists
	IPC_NOWAIT = unix.IPC_NOWAIT // Error if request must wait
)

// Private key.
const IPC_PRIVATE = unix.IPC_PRIVATE

// Control commands.
const (
	IPC_RMID = unix.IPC_RMID // Remove identifier
	IPC_SET  = unix.IPC_SET  // Set options
	IPC_STAT = unix.IPC_STAT // Get options
)

// Common mode bits.
const (
	IPC_R = 0o000400 // Read permission
	IPC_W = 0o000200 // Write/alter permission
	IPC_M = 0o010000 // Modify control info permission
)

var _zero uintptr

// Perm corresponds to struct ipc_perm.
type Perm struct{ s ipc_perm }

func (p *Perm) Uid() int  { return int(p.s.Uid) }
func (p *Perm) Gid() int  { return int(p.s.Gid) }
func (p *Perm) Cuid() int { return int(p.s.Cuid) }
func (p *Perm) Cgid() int { return int(p.s.Cgid) }
func (p *Perm) Mode() int { return int(p.s.Mode) }

func (p *Perm) SetUid(uid int) { p.s.Uid = uint32(uid) }
func (p *Perm) SetGid(gid int) { p.s.Gid = uint32(gid) }

// Ftok creates IPC identifier from path name.
func Ftok(path string, id uint8) (int, error) {
	var stat unix.Stat_t
	if err := unix.Stat(path, &stat); err != nil {
		return -1, err
	}

	key := int(id&0xff)<<24 |
		int((stat.Dev&0xff)<<16) |
		int(stat.Ino&0xffff)

	return key, nil
}
