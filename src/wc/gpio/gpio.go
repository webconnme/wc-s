package gpio

import (
	"syscall"
	"os"
	"reflect"
	"unsafe"
)

var config = struct {
	base uint64
	size uint
}{
	base : 0xC000A000,
	size : 0x1000,
}

type Gpio struct {
	file *os.File
	buffer []byte
}

func mmap(fd uintptr, off int64, len int) ([]byte, error) {
	b, err := syscall.Mmap(int(fd), 
		off, 
		len, 
		syscall.PROT_READ | syscall.PROT_WRITE, 
		syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (g *Gpio) Alt(module uint, bit uint, value int) {
	addr := 0x40 * module + 8

	low_mask := byte(1 << ((bit % 4) * 2))
	high_mask := byte(1 << ((bit % 4) * 2 + 1))

	g.buffer[addr + (bit/4)] &= ^(high_mask | low_mask)

	if value & 0x01 == 0x01 {
		g.buffer[addr + (bit/4)] |= low_mask
	}

	if value & 0x02 == 0x02 {
		g.buffer[addr + (bit/4)] |= high_mask
	}
}

func (g *Gpio) Dir(module uint, bit uint, value int) {
	addr := 0x40 * module + 1

	if (value != 0) {
		g.buffer[addr + (bit/8)] |= byte(1 << (bit % 8))
	} else {
		g.buffer[addr + (bit/8)] &= byte(^(1 << (bit % 8)))
	}
}

func (g *Gpio) Val(module uint, bit uint, value int) {
	addr := 0x40 * module

	if (value != 0) {
		g.buffer[addr + (bit/8)] |= byte(1 << (bit % 8))
	} else {
		g.buffer[addr + (bit/8)] &= byte(^(1 << (bit % 8)))
	}
}

func Open() (*Gpio, error) {
	var err error
	g := new(Gpio)

	g.file, err = os.OpenFile("/dev/mem", os.O_RDWR | os.O_SYNC, 0777)
	if err != nil {
		return nil, err
	}

	g.buffer, err = mmap(uintptr(g.file.Fd()), 0xC000A000, 0x1000)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Gpio) Close() error {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&g.buffer))

	_, _, errno := syscall.Syscall(syscall.SYS_MUNMAP, header.Data, uintptr(config.size), 0)
	if errno != 0 {
		return syscall.Errno(errno)
	}

	g.file.Close()
	return nil
}