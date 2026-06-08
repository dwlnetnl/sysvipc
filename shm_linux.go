package sysvipc

import (
	"structs"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// shmat() shmflg values
const (
	SHM_RDONLY = 0o010000  // read-only access
	SHM_RND    = 0o020000  // round attach address to SHMLBA boundary
	SHM_REMAP  = 0o040000  // take-over region on attach
	SHM_EXEC   = 0o0100000 // execution access
)

// super user shmctl commands
const (
	SHM_LOCK   = 11
	SHM_UNLOCK = 12
)

// ipcs ctl commands
const (
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

func shmctl(id, cmd int, buf any) (int, error) {
	var dest *unix.SysvShmDesc

	switch cmd {
	case IPC_STAT, IPC_SET, SHM_STAT, SHM_STAT_ANY:
		switch v := buf.(type) {
		case *ShmControl:
			dest = &v.s
		case ShmControl:
			dest = &v.s
		default:
			panic("buf is not a *ShmControl value")
		}
	case IPC_INFO:
		switch v := buf.(type) {
		case *ShmInfo:
			dest = (*unix.SysvShmDesc)(unsafe.Pointer(v))
		case ShmInfo:
			dest = (*unix.SysvShmDesc)(unsafe.Pointer(&v))
		default:
			panic("buf is not a *ShmInfo value")
		}
	case SHM_INFO:
		switch v := buf.(type) {
		case *ShmSystemInfo:
			dest = (*unix.SysvShmDesc)(unsafe.Pointer(v))
		case ShmSystemInfo:
			dest = (*unix.SysvShmDesc)(unsafe.Pointer(&v))
		default:
			panic("buf is not a *ShmSystemInfo value")
		}
	}

	return unix.SysvShmCtl(id, cmd, dest)
}
