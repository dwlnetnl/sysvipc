package sysvipc

import "golang.org/x/sys/unix"

// Possible flag values for Shmat.
const (
	SHM_RDONLY = unix.SHM_RDONLY // read-only access
	SHM_RND    = unix.SHM_RND    // round attach address to SHMLBA boundary
)

// ShmDesc corresponds to struct shmid_ds.
type ShmDesc struct{ s shmid_ds }

func (s *ShmDesc) Perm() *Perm    { return &Perm{s.s.Perm} }
func (s *ShmDesc) Segsz() uint64  { return uint64(s.s.Segsz) }
func (s *ShmDesc) Lpid() int      { return int(s.s.Lpid) }
func (s *ShmDesc) Cpid() int      { return int(s.s.Cpid) }
func (s *ShmDesc) Nattch() uint64 { return uint64(s.s.Segsz) }

// Shmget corresponds to shmget.
func Shmget(key, size, flag int) (int, error) {
	return unix.SysvShmGet(key, size, flag)
}

// Shmat corresponds to shmat.
func Shmat(id int, addr uintptr, flag int) ([]byte, error) {
	return unix.SysvShmAttach(id, addr, flag)
}

// Shmdt corresponds to shmdt.
func Shmdt(data []byte) error {
	return unix.SysvShmDetach(data)
}
