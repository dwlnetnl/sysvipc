package sysvipc

import (
	"structs"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Additional commands for semctl.
const (
	GETPID       = 11 // get sempid
	GETVAL       = 12 // get semval
	GETALL       = 13 // get all semval's
	GETNCNT      = 14 // get semncnt
	GETZCNT      = 15 // get semzcnt
	SETVAL       = 16 // set semval
	SETALL       = 17 // set all semval's
	SEM_STAT     = 18
	SEM_INFO     = 19
	SEM_STAT_ANY = 20
)

const (
	SEM_DEST   = 0o2000 // Semaphore will be destroyed on last detach
	SEM_LOCKED = 0o4000 // Semaphore set is locked
)

// semid_ds corresponds to struct semid_ds.
type semid_ds struct {
	// TODO: probably will not work on some architectures, e.g. 32-bit?
	_     structs.HostLayout
	perm  ipc_perm // semaphore permissions
	otime int64    // last semop time
	ctime int64    // creation time/time of last modification via semctl()
	nsems uint64   // # of semaphores in set
	_     uint64
	_     uint64
}

const (
	_ = uint(unsafe.Sizeof(semid_ds{}) - 88)
	_ = uint(88 - unsafe.Sizeof(semid_ds{}))
	_ = uint(unsafe.Alignof(semid_ds{}) - 8)
	_ = uint(8 - unsafe.Alignof(semid_ds{}))
)

func (s *SemControl) Otime() time.Time { return time.Unix(s.s.otime, 0) }
func (s *SemControl) Ctime() time.Time { return time.Unix(s.s.ctime, 0) }

type SemInfo struct {
	_   structs.HostLayout
	Map int32 // Number of entries in semaphore map; unused within kernel
	Mni int32 // Maximum number of semaphore sets
	Mns int32 // Maximum number of semaphores in all semaphore sets
	Mnu int32 // System-wide maximum number of undo structures; unused within kernel
	Msl int32 // Maximum number of semaphores in a set
	Opm int32 // Maximum number of operations for semop
	Ume int32 // Maximum number of undo entries per process; unused within kernel
	Usz int32 // Size of struct sem_undo
	Vmx int32 // Maximum semaphore value
	Aem int32 // Maximum value that can be recorded for semaphore adjustment (SEM_UNDO)
}

const (
	_ = uint(unsafe.Sizeof(SemInfo{}) - 40)
	_ = uint(40 - unsafe.Sizeof(SemInfo{}))
	_ = uint(unsafe.Alignof(SemInfo{}) - 4)
	_ = uint(4 - unsafe.Alignof(SemInfo{}))
)

func semget(key, nsems, flag int) (int, error) {
	r0, _, e1 := unix.Syscall(unix.SYS_SEMGET, uintptr(key), uintptr(nsems), uintptr(flag))
	if int(r0) == -1 {
		return -1, e1
	}
	return int(r0), nil
}

func semop(id int, ops []Sembuf) (err error) {
	var p0 unsafe.Pointer
	if len(ops) > 0 {
		p0 = unsafe.Pointer(&ops[0])
	}

	r0, _, e1 := unix.Syscall(unix.SYS_SEMOP, uintptr(id), uintptr(p0), uintptr(len(ops)))
	if int(r0) == -1 {
		err = e1
	}

	return err
}

// Semtimedop corresponds to semtimedop.
func Semtimedmop(id int, ops []Sembuf, timeout time.Duration) (err error) {
	var p0 unsafe.Pointer
	if len(ops) > 0 {
		p0 = unsafe.Pointer(&ops[0])
	}

	var p1 unsafe.Pointer
	switch {
	case timeout > 0:
		ts := unix.NsecToTimespec(timeout.Nanoseconds())
		p1 = unsafe.Pointer(&ts)
	case timeout == 0:
		p1 = unsafe.Pointer(&_zero)
	}

	r0, _, e1 := unix.Syscall6(unix.SYS_SEMTIMEDOP, uintptr(id), uintptr(p0), uintptr(len(ops)), uintptr(p1), 0, 0)
	if int(r0) == -1 {
		err = e1
	}

	return err
}

// Semctl corresponds to semctl.
//
// arg can be of type:
//
//	*SemControl    IPC_STAT, IPC_SET, SEM_STAT, SEM_STAT_ANY
//	*SemInfo       IPC_INFO, SEM_INFO
//	[]uint16       GETALL, SETALL
//	int32          SETVAL
func Semctl(id, num, cmd int, arg any) (int, error) {
	var r0 uintptr
	var e1 unix.Errno

	if cmd != SETVAL {
		var p0 unsafe.Pointer
		switch cmd {
		case IPC_STAT, IPC_SET, SEM_STAT, SEM_STAT_ANY:
			switch v := arg.(type) {
			case *SemControl:
				p0 = unsafe.Pointer(v)
			case SemControl:
				p0 = unsafe.Pointer(&v)
			default:
				panic("arg is not a *SemControl value")
			}
		case IPC_INFO, SEM_INFO:
			switch v := arg.(type) {
			case *SemInfo:
				p0 = unsafe.Pointer(v)
			case SemInfo:
				p0 = unsafe.Pointer(&v)
			default:
				panic("arg is not a *SemInfo value")
			}
		case GETALL, SETALL:
			switch v := arg.(type) {
			case []uint16:
				p0 = unsafe.Pointer(&v[0])
			case *[]uint16:
				p0 = unsafe.Pointer(&(*v)[0])
			default:
				panic("arg is not a []uint16 value")
			}
		default:
			p0 = unsafe.Pointer(&_zero)
		}
		r0, _, e1 = unix.Syscall6(unix.SYS_MSGCTL, uintptr(id), uintptr(num), uintptr(cmd), uintptr(p0), 0, 0)

	} else {
		var val int32
		switch v := arg.(type) {
		case int32:
			val = v
		case *int32:
			val = *v
		default:
			panic("arg is not a int32 value")
		}
		r0, _, e1 = unix.Syscall6(unix.SYS_MSGCTL, uintptr(id), uintptr(num), SETVAL, uintptr(val), 0, 0)
	}

	if int(r0) == -1 {
		return -1, e1
	}

	return int(r0), nil
}
