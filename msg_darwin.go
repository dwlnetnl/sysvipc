package sysvipc

import (
	"encoding/binary"
	"structs"
	"syscall"
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

func (m *MsgControl) Cbytes() uint64 {
	return binary.NativeEndian.Uint64(m.s.cbytes[:])
}

func (m *MsgControl) Qnum() uint64 {
	return binary.NativeEndian.Uint64(m.s.qnum[:])
}

func (m *MsgControl) Qbytes() uint64 {
	return binary.NativeEndian.Uint64(m.s.qbytes[:])
}

func (m *MsgControl) Stime() time.Time {
	t := int64(binary.NativeEndian.Uint64(m.s.stime[:]))
	return time.Unix(t, 0)
}

func (m *MsgControl) Rtime() time.Time {
	t := int64(binary.NativeEndian.Uint64(m.s.rtime[:]))
	return time.Unix(t, 0)
}

func (m *MsgControl) Ctime() time.Time {
	t := int64(binary.NativeEndian.Uint64(m.s.ctime[:]))
	return time.Unix(t, 0)
}

func (m *MsgControl) SetQbytes(qbytes uint64) {
	binary.NativeEndian.PutUint64(m.s.qbytes[:], qbytes)
}

//go:linkname runtime_syscall syscall.syscall
//go:linkname runtime_syscall6 syscall.syscall6

func runtime_syscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
func runtime_syscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

//go:cgo_import_dynamic libSystem_msgget msgget "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libSystem_msgctl msgctl "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libSystem_msgsnd msgsnd "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libSystem_msgrcv msgrcv "/usr/lib/libSystem.B.dylib"

var (
	libSystem_msgget_trampoline_addr uintptr
	libSystem_msgctl_trampoline_addr uintptr
	libSystem_msgsnd_trampoline_addr uintptr
	libSystem_msgrcv_trampoline_addr uintptr
)

func msgget(key, msgflg int) (int, error) {
	r0, _, e1 := runtime_syscall(libSystem_msgget_trampoline_addr, uintptr(key), uintptr(msgflg), 0)
	if int(r0) == -1 {
		return -1, e1
	}
	return int(r0), nil
}

func msgsnd(qid int, msg unsafe.Pointer, msgsz, msgflg int) (err error) {
	r0, _, e1 := runtime_syscall6(libSystem_msgsnd_trampoline_addr, uintptr(qid), uintptr(msg), uintptr(msgsz), uintptr(msgflg), 0, 0)
	if int(r0) == -1 {
		err = e1
	}
	return err
}

func msgrcv(qid int, msg unsafe.Pointer, msgsz, msgtyp, msgflg int) (int, error) {
	r0, _, e1 := runtime_syscall6(libSystem_msgrcv_trampoline_addr, uintptr(qid), uintptr(msg), uintptr(msgsz), uintptr(msgtyp), uintptr(msgflg), 0)
	if int(r0) == -1 {
		return -1, e1
	}

	return int(r0), nil
}

func msgctl(qid, cmd int, buf any) (int, error) {
	var p0 unsafe.Pointer
	switch cmd {
	case unix.IPC_STAT, unix.IPC_SET:
		switch v := buf.(type) {
		case *MsgControl:
			p0 = unsafe.Pointer(v)
		case MsgControl:
			p0 = unsafe.Pointer(&v)
		default:
			panic("buf is not a *MsgControl value")
		}
	default:
		p0 = unsafe.Pointer(&_zero)
	}

	r0, _, e1 := runtime_syscall(libSystem_msgctl_trampoline_addr, uintptr(qid), uintptr(cmd), uintptr(p0))
	if int(r0) == -1 {
		return -1, e1
	}

	return int(r0), nil
}
