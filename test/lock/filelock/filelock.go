package filelock

import (
	"os"
	"syscall"
	"fmt"
	"time"
)

//文件锁
type FileLocker struct {
	dir string
	f   *os.File
}

func New(dir string) *FileLocker {
	return &FileLocker{
		dir: dir,
	}
}

//加锁
func (l *FileLocker) Lock() error {
	f, err := os.Open(l.dir)
	if err != nil {
		return err
	}
	l.f = f
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("cannot flock directory %s - %s", l.dir, err)
	}
	return nil
}

//释放锁
func (l *FileLocker) Unlock() error {
	defer l.f.Close()
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}

type FilelockScrambler struct {
	FileLocker *FileLocker
	Locked     bool
	StopCh     chan struct{}
}

func (s *FilelockScrambler) Scramble() {
	for {
		select {
		case <- s.StopCh:
			break
		default:
			err := s.FileLocker.Lock()
			if err != nil {
				s.Locked = false
			} else {
				s.Locked = true
			}
			time.Sleep(time.Second * 5)
		}

	}
}
