package sysvipc

import "unsafe"

// Possible values for the fifth parameter to Msgrcv,
// in addition to the IPC_NOWAIT flag, which is permitted.
const (
	MSG_NOERROR = 0o010000 // no error if message is too big
)

// MsgControl corresponds to struct msqid_ds.
type MsgControl struct{ s msqid_ds }

func (m *MsgControl) Perm() *Perm { return &Perm{m.s.perm} }
func (m *MsgControl) Lspid() int  { return int(m.s.lspid) }
func (m *MsgControl) Lrpid() int  { return int(m.s.lrpid) }

// Msgget corresponds to msgget.
func Msgget(key, flag int) (int, error) {
	return msgget(key, flag)
}

// Msgsnd corresponds to msgsnd.
func Msgsnd(id int, msg []byte, flag int) error {
	var p0 unsafe.Pointer
	if len(msg) > 0 {
		p0 = unsafe.Pointer(&msg[0])
	} else {
		p0 = unsafe.Pointer(&_zero)
	}
	return msgsnd(id, p0, len(msg), flag)
}

// Msgrcv corresponds to msgrcv.
func Msgrcv(id int, msg []byte, typ, flag int) (int, error) {
	var p0 unsafe.Pointer
	if len(msg) > 0 {
		p0 = unsafe.Pointer(&msg[0])
	} else {
		p0 = unsafe.Pointer(&_zero)
	}
	return msgrcv(id, p0, len(msg), typ, flag)
}
