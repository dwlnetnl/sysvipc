package sysvipc

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestMsg_RoundTrip(t *testing.T) {
	qid, err := Msgget(IPC_PRIVATE, IPC_CREAT|IPC_EXCL|0o666)
	if err != nil {
		t.Fatal("msgget:", err)
	}

	t.Cleanup(func() {
		if _, err := Msgctl(qid, IPC_RMID, nil); err != nil {
			t.Fatal("msgctl IPC_RMID:", err)
		}
	})

	buf := make([]byte, 0, 256)
	buf = binary.NativeEndian.AppendUint64(buf, 1)
	buf = append(buf, "hello world!"...)

	if err := Msgsnd(qid, buf, 0); err != nil {
		t.Fatal("msgsnd:", err)
	}

	var c MsgControl
	if _, err := Msgctl(qid, IPC_STAT, &c); err != nil {
		t.Fatal("msgctl IPC_STAT:", err)
	}
	if got, want := c.Qnum(), uint64(1); got != want {
		t.Errorf("qnum = %d, want: %d", got, want)
	}
	if got, want := c.Cbytes(), uint64(len(buf)); got != want {
		t.Errorf("cbytes = %d, want: %d", got, want)
	}

	var msg [256]byte
	n, err := Msgrcv(qid, msg[:], 1, IPC_NOWAIT)
	if n != len(buf) || err != nil {
		t.Errorf("msgrcv: %d, %v, want: %d, <nil>", n, len(buf), err)
	}
	if !bytes.Equal(msg[:n], buf) {
		t.Errorf("msg = %q, want: %q", msg[:n], buf)
	}
}
