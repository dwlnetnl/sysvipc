package sysvipc

import "golang.org/x/sys/unix"

// shmid_ds corresponds to struct shmid_ds.
type shmid_ds = unix.SysvShmDesc

func shmctl(id, cmd int, buf any) (int, error) {
	var dest *unix.SysvShmDesc

	switch cmd {
	case IPC_STAT, IPC_SET:
		switch v := buf.(type) {
		case *ShmControl:
			dest = &v.s
		case ShmControl:
			dest = &v.s
		default:
			panic("buf is not a *ShmControl value")
		}
	}

	return unix.SysvShmCtl(id, cmd, dest)
}
