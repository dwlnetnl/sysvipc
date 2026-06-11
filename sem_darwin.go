package sysvipc

import (
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Commands for semctl.
const (
	GETNCNT = 3 // Return the value of semncnt {READ}
	GETPID  = 4 // Return the value of sempid {READ}
	GETVAL  = 5 // Return the value of semval {READ}
	GETALL  = 6 // Return semvals into arg.array {READ}
	GETZCNT = 7 // Return the value of semzcnt {READ}
	SETVAL  = 8 // Set the value of semval to arg.val {ALTER}
	SETALL  = 9 // Set semvals from arg.array {ALTER}
)

// semid_ds corresponds to struct semid_ds.
type semid_ds struct {
	perm  ipc_perm // semaphore permissions
	base  int32    // 32 bit base ptr for semaphore set
	nsems uint16   // number of sems in set
	_     int16    //
	otime [8]byte  // last operation time
	_     int32    //
	_     int32    //
	ctime [8]byte  // last change time
	_     int32
	_     [4]int32
}

const (
	_ = uint(unsafe.Sizeof(semid_ds{}) - 76)
	_ = uint(76 - unsafe.Sizeof(semid_ds{}))
	_ = uint(unsafe.Alignof(semid_ds{}) - 4)
	_ = uint(4 - unsafe.Alignof(semid_ds{}))
)

func (s *SemDesc) Otime() time.Time { return bytesToTime(s.s.otime[:]) }
func (s *SemDesc) Ctime() time.Time { return bytesToTime(s.s.ctime[:]) }

//go:cgo_import_dynamic libSystem_semget semget "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libSystem_semop semop "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libSystem_semctl semctl "/usr/lib/libSystem.B.dylib"

var (
	libSystem_semget_trampoline_addr uintptr
	libSystem_semop_trampoline_addr  uintptr
	libSystem_semctl_trampoline_addr uintptr
)

func semget(key, nsems, flag int) (int, error) {
	r0, _, e1 := runtime_syscall(libSystem_semget_trampoline_addr, uintptr(key), uintptr(nsems), uintptr(flag))
	if int(r0) == -1 {
		return -1, e1
	}
	return int(r0), nil
}

func semop(id int, ops []Sembuf) (err error) {
	r0, _, e1 := runtime_syscall(libSystem_semop_trampoline_addr, uintptr(id), uintptr(unsafe.Pointer(&ops[0])), uintptr(len(ops)))
	if int(r0) == -1 {
		err = e1
	}
	return err
}

// Semctl corresponds to semctl.
//
// arg can be of type:
//
//	*SemControl    IPC_STAT, IPC_SET
//	[]uint16       GETALL, SETALL
//	int32          SETVAL
func Semctl(id, num, cmd int, arg any) (int, error) {
	var r0 uintptr
	var e1 unix.Errno

	if cmd != SETVAL {
		var p0 unsafe.Pointer
		switch cmd {
		case IPC_STAT, IPC_SET:
			switch v := arg.(type) {
			case *SemDesc:
				p0 = unsafe.Pointer(v)
			case SemDesc:
				p0 = unsafe.Pointer(&v)
			default:
				panic("arg is not a *SemControl value")
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
		r0, _, e1 = runtime_syscall6(libSystem_semctl_trampoline_addr, uintptr(id), uintptr(num), uintptr(cmd), uintptr(p0), 0, 0)

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
		r0, _, e1 = runtime_syscall6(libSystem_semctl_trampoline_addr, uintptr(id), uintptr(num), SETVAL, uintptr(val), 0, 0)
	}

	if int(r0) == -1 {
		return -1, e1
	}

	return int(r0), nil
}
