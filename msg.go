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
func Msgget(key, msgflg int) (int, error) {
	return msgget(key, msgflg)
}

// Msgsnd corresponds to msgsnd.
func Msgsnd(qid int, msg []byte, msgflg int) error {
	var p0 unsafe.Pointer
	if len(msg) > 0 {
		p0 = unsafe.Pointer(&msg[0])
	} else {
		p0 = unsafe.Pointer(&_zero)
	}
	return msgsnd(qid, p0, len(msg), msgflg)
}

// Msgrcv corresponds to msgrcv.
func Msgrcv(qid int, msg []byte, msgtyp, msgflg int) (int, error) {
	var p0 unsafe.Pointer
	if len(msg) > 0 {
		p0 = unsafe.Pointer(&msg[0])
	} else {
		p0 = unsafe.Pointer(&_zero)
	}
	return msgrcv(qid, p0, len(msg), msgtyp, msgflg)
}

// Msgctl corresponds to msgctl.
func Msgctl(qid, cmd int, buf any) (int, error) {
	return msgctl(qid, cmd, buf)
}
