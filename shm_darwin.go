package sysvipc

import "golang.org/x/sys/unix"

// shmid_ds corresponds to struct shmid_ds.
type shmid_ds = unix.SysvShmDesc

// Shmctl corresponds to shmctl.
//
// arg can be of type:
//
//	*ShmControl    IPC_STAT, IPC_SET
func Shmctl(id, cmd int, buf any) (int, error) {
	var desc *unix.SysvShmDesc

	switch cmd {
	case IPC_STAT, IPC_SET:
		switch v := buf.(type) {
		case *ShmControl:
			desc = &v.s
		case ShmControl:
			desc = &v.s
		default:
			panic("arg is not a *ShmControl value")
		}
	}

	return unix.SysvShmCtl(id, cmd, desc)
}
