package sysvipc

// Possible flag values for Semop.
const (
	SEM_UNDO = 0x1000 // undo the operation on exit
)

// SemDesc corresponds to struct semid_ds.
type SemDesc struct{ s semid_ds }

// Sembuf corresponds to struct sembuf.
type Sembuf struct {
	Num uint16
	Op  int16
	Flg int16
}

// Semget corresponds to semget.
func Semget(key, nsems, flag int) (int, error) {
	return semget(key, nsems, flag)
}

// Semop corresponds to semop.
func Semop(id int, ops []Sembuf) error {
	return semop(id, ops)
}
