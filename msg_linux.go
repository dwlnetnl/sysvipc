package sysvipc

import (
	"structs"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// ipcs ctl commands.
const (
	MSG_STAT     = 11
	MSG_INFO     = 12
	MSG_STAT_ANY = 13
)

// msqid_ds corresponds to struct msqid_ds (new layout).
type msqid_ds struct {
	_      structs.HostLayout
	perm   ipc_perm // msg queue permissions
	stime  int64    // time of last msgsnd()
	rtime  int64    // time of last msgrcv()
	ctime  int64    // time of last change
	cbytes uint64   // current bytes on queue
	qnum   uint64   // number of msgs on the queue
	qbytes uint64   // max bytes on the queue
	lspid  int32    // pid of last msgsnd()
	lrpid  int32    // pid of last msgrcv()
	_      uint64   // RESERVED
	_      uint64   // RESERVED
}

const (
	_ = uint(unsafe.Sizeof(msqid_ds{}) - 120)
	_ = uint(120 - unsafe.Sizeof(msqid_ds{}))
	_ = uint(unsafe.Alignof(msqid_ds{}) - 8)
	_ = uint(8 - unsafe.Alignof(msqid_ds{}))
)

func (m *MsgControl) Cbytes() uint64          { return m.s.cbytes }
func (m *MsgControl) Qnum() uint64            { return m.s.qnum }
func (m *MsgControl) Qbytes() uint64          { return m.s.qbytes }
func (m *MsgControl) Stime() time.Time        { return time.Unix(m.s.stime, 0) }
func (m *MsgControl) Rtime() time.Time        { return time.Unix(m.s.rtime, 0) }
func (m *MsgControl) Ctime() time.Time        { return time.Unix(m.s.ctime, 0) }
func (m *MsgControl) SetQbytes(qbytes uint64) { m.s.qbytes = qbytes }

// MsgInfo corresponds to struct msginfo.
type MsgInfo struct {
	_    structs.HostLayout
	Pool int32  // size in kibibytes of buffer pool (unused within kernel)
	Map  int32  // max number of entries in message map (unused within kernel)
	Max  int32  // max number of bytes per message
	Mnb  int32  // max number of bytes on queue
	Mni  int32  // max number of message queues
	Ssz  int32  // message segment size (unused within kernel)
	Tql  int32  // max number of messages on all queues (unused within kernel)
	Seg  uint16 // max number of segments (unused within kernel)
	_    [2]byte
}

const (
	_ = uint(unsafe.Sizeof(MsgInfo{}) - 32)
	_ = uint(32 - unsafe.Sizeof(MsgInfo{}))
	_ = uint(unsafe.Alignof(MsgInfo{}) - 4)
	_ = uint(4 - unsafe.Alignof(MsgInfo{}))
)

func msgget(key, msgflg int) (int, error) {
	r0, _, e1 := unix.Syscall(unix.SYS_MSGGET, uintptr(key), uintptr(msgflg), 0)
	if int(r0) == -1 {
		return -1, e1
	}
	return int(r0), nil
}

func msgsnd(qid int, msg unsafe.Pointer, msgsz, msgflg int) (err error) {
	r0, _, e1 := unix.Syscall6(unix.SYS_MSGSND, uintptr(qid), uintptr(msg), uintptr(msgsz), uintptr(msgflg), 0, 0)
	if int(r0) == -1 {
		err = e1
	}
	return err
}

func msgrcv(qid int, msg unsafe.Pointer, msgsz, msgtyp, msgflg int) (int, error) {
	r0, _, e1 := unix.Syscall6(unix.SYS_MSGRCV, uintptr(qid), uintptr(msg), uintptr(msgsz), uintptr(msgtyp), uintptr(msgflg), 0)
	if int(r0) == -1 {
		return -1, e1
	}

	return int(r0), nil
}

func msgctl(qid, cmd int, buf any) (int, error) {
	var p0 unsafe.Pointer
	switch cmd {
	case IPC_STAT, IPC_SET, MSG_STAT:
		switch v := buf.(type) {
		case *MsgControl:
			p0 = unsafe.Pointer(v)
		case MsgControl:
			p0 = unsafe.Pointer(&v)
		default:
			panic("buf is not a *MsgControl value")
		}
	case IPC_INFO, MSG_INFO:
		switch v := buf.(type) {
		case *MsgInfo:
			p0 = unsafe.Pointer(v)
		case MsgInfo:
			p0 = unsafe.Pointer(&v)
		default:
			panic("buf is not a *MsgInfo value")
		}
	default:
		p0 = unsafe.Pointer(&_zero)
	}

	r0, _, e1 := unix.Syscall(unix.SYS_MSGCTL, uintptr(qid), uintptr(cmd), uintptr(p0))
	if int(r0) == -1 {
		return -1, e1
	}

	return int(r0), nil
}
