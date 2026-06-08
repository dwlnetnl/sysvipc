package sysvipc

import (
	"encoding/binary"
	"structs"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// msqid_ds corresponds to struct msqid_ds.
type msqid_ds struct {
	_      structs.HostLayout
	perm   ipc_perm // msg queue permissions
	_      int32    // RESERVED: kernel use only
	_      int32    // RESERVED: kernel use only
	cbytes [8]byte  // # of bytes on the queue
	qnum   [8]byte  // number of msgs on the queue
	qbytes [8]byte  // max bytes on the queue
	lspid  int32    // pid of last msgsnd()
	lrpid  int32    // pid of last msgrcv()
	stime  [8]byte  // time of last msgsnd()
	_      int32    //
	rtime  [8]byte  // time of last msgrcv()
	_      int32    //
	ctime  [8]byte  // time of last msgctl()
	_      int32    //
	_      [4]int32 //
}

const (
	_ = uint(unsafe.Sizeof(msqid_ds{}) - 116)
	_ = uint(116 - unsafe.Sizeof(msqid_ds{}))
	_ = uint(unsafe.Alignof(msqid_ds{}) - 4)
	_ = uint(4 - unsafe.Alignof(msqid_ds{}))
)

func (m *MsgControl) Cbytes() uint64   { return bytesToUint64(m.s.cbytes[:]) }
func (m *MsgControl) Qnum() uint64     { return bytesToUint64(m.s.qnum[:]) }
func (m *MsgControl) Qbytes() uint64   { return bytesToUint64(m.s.qbytes[:]) }
func (m *MsgControl) Stime() time.Time { return bytesToTime(m.s.stime[:]) }
func (m *MsgControl) Rtime() time.Time { return bytesToTime(m.s.rtime[:]) }
func (m *MsgControl) Ctime() time.Time { return bytesToTime(m.s.ctime[:]) }

func (m *MsgControl) SetQbytes(qbytes uint64) {
	binary.NativeEndian.PutUint64(m.s.qbytes[:], qbytes)
}

//go:cgo_import_dynamic libSystem_msgget msgget "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libSystem_msgsnd msgsnd "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libSystem_msgrcv msgrcv "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libSystem_msgctl msgctl "/usr/lib/libSystem.B.dylib"

var (
	libSystem_msgget_trampoline_addr uintptr
	libSystem_msgsnd_trampoline_addr uintptr
	libSystem_msgrcv_trampoline_addr uintptr
	libSystem_msgctl_trampoline_addr uintptr
)

func msgget(key, flag int) (int, error) {
	r0, _, e1 := runtime_syscall(libSystem_msgget_trampoline_addr, uintptr(key), uintptr(flag), 0)
	if int(r0) == -1 {
		return -1, e1
	}
	return int(r0), nil
}

func msgsnd(id int, msg unsafe.Pointer, size, flag int) (err error) {
	r0, _, e1 := runtime_syscall6(libSystem_msgsnd_trampoline_addr, uintptr(id), uintptr(msg), uintptr(size), uintptr(flag), 0, 0)
	if int(r0) == -1 {
		err = e1
	}
	return err
}

func msgrcv(id int, msg unsafe.Pointer, size, typ, flag int) (int, error) {
	r0, _, e1 := runtime_syscall6(libSystem_msgrcv_trampoline_addr, uintptr(id), uintptr(msg), uintptr(size), uintptr(typ), uintptr(flag), 0)
	if int(r0) == -1 {
		return -1, e1
	}
	return int(r0), nil
}

// Msgctl corresponds to msgctl.
//
// arg can be of type:
//
//	*MsgControl    IPC_STAT, IPC_SET
func Msgctl(id, cmd int, arg any) (int, error) {
	var p0 unsafe.Pointer
	switch cmd {
	case unix.IPC_STAT, unix.IPC_SET:
		switch v := arg.(type) {
		case *MsgControl:
			p0 = unsafe.Pointer(v)
		case MsgControl:
			p0 = unsafe.Pointer(&v)
		default:
			panic("arg is not a *MsgControl value")
		}
	default:
		p0 = unsafe.Pointer(&_zero)
	}

	r0, _, e1 := runtime_syscall(libSystem_msgctl_trampoline_addr, uintptr(id), uintptr(cmd), uintptr(p0))
	if int(r0) == -1 {
		return -1, e1
	}

	return int(r0), nil
}
