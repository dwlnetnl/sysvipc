package sysvipc

import "golang.org/x/sys/unix"

// Possible flag values which may be OR'ed into the third
// argument to Shmat.
const (
	SHM_RDONLY = unix.SHM_RDONLY // read-only access
	SHM_RND    = unix.SHM_RND    // round attach address to SHMLBA boundary
)

// ShmControl corresponds to struct shmid_ds.
type ShmControl struct{ s shmid_ds }

func (s *ShmControl) Perm() *Perm    { return &Perm{s.s.Perm} }
func (s *ShmControl) Segsz() uint64  { return uint64(s.s.Segsz) }
func (s *ShmControl) Lpid() int      { return int(s.s.Lpid) }
func (s *ShmControl) Cpid() int      { return int(s.s.Cpid) }
func (s *ShmControl) Nattch() uint64 { return uint64(s.s.Segsz) }

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

// Shmctl corresponds to shmctl.
func Shmctl(id, cmd int, buf any) (int, error) {
	return shmctl(id, cmd, buf)
}
