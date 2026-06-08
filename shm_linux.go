package sysvipc

import (
	"structs"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Additional flag values for Shmat.
const (
	SHM_REMAP = 0o040000  // take-over region on attach
	SHM_EXEC  = 0o0100000 // execution access
)

// Additional commands for Shmctl.
const (
	SHM_LOCK     = 11
	SHM_UNLOCK   = 12
	SHM_STAT     = 13
	SHM_INFO     = 14
	SHM_STAT_ANY = 15
)

// shmid_ds corresponds to struct shmid_ds.
type shmid_ds = unix.SysvShmDesc

func (s *ShmControl) Atime() time.Time { return time.Unix(s.s.Atime, 0) }
func (s *ShmControl) Dtime() time.Time { return time.Unix(s.s.Dtime, 0) }
func (s *ShmControl) Ctime() time.Time { return time.Unix(s.s.Ctime, 0) }

// ShmSystemInfo corresponds to struct shminfo64.
type ShmSystemInfo struct {
	_   structs.HostLayout
	Max uint // maximum segment size
	Min uint // minimum segment size; always 1
	Mni uint // maximum number of segments
	Seg uint // maximum number of segments that a process can attach (unused within kernel)
	All uint // maximum number of pages of shared memory, system-wide
	_   [4]uint
}

// ShmInfo corresponds to struct shm_info.
type ShmInfo struct {
	_             structs.HostLayout
	UsedIDs       int32 // number of currently existing segments
	ShmTot        uint  // total number of shared memory pages
	ShmRss        uint  // number of resident shared memory pages
	ShmSwp        uint  // number of swapped shared memory pages
	SwapAttempts  uint  // unused since Linux 2.4
	SwapSuccesses uint  // unused since Linux 2.4
}

// Shmctl corresponds to shmctl.
//
// arg can be of type:
//
//	*ShmControl    IPC_STAT, IPC_SET, SHM_STAT, SHM_STAT_ANY
//	*ShmInfo       IPC_INFO, SHM_INFO
func Shmctl(id, cmd int, arg any) (int, error) {
	var desc *unix.SysvShmDesc

	switch cmd {
	case IPC_STAT, IPC_SET, SHM_STAT, SHM_STAT_ANY:
		switch v := arg.(type) {
		case *ShmControl:
			desc = &v.s
		case ShmControl:
			desc = &v.s
		default:
			panic("arg is not a *ShmControl value")
		}
	case IPC_INFO:
		switch v := arg.(type) {
		case *ShmInfo:
			desc = (*unix.SysvShmDesc)(unsafe.Pointer(v))
		case ShmInfo:
			desc = (*unix.SysvShmDesc)(unsafe.Pointer(&v))
		default:
			panic("arg is not a *ShmInfo value")
		}
	case SHM_INFO:
		switch v := arg.(type) {
		case *ShmSystemInfo:
			desc = (*unix.SysvShmDesc)(unsafe.Pointer(v))
		case ShmSystemInfo:
			desc = (*unix.SysvShmDesc)(unsafe.Pointer(&v))
		default:
			panic("arg is not a *ShmSystemInfo value")
		}
	}

	return unix.SysvShmCtl(id, cmd, desc)
}
