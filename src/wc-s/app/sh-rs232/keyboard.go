package main
// https://groups.google.com/d/msg/golang-nuts/8o9fxPaeFu8/uSFYfobL5EgJ
// http://go.pastie.org/813153

import (
	"syscall"
	"unsafe"
	"os"
	"fmt"
)

var oldTermios syscall.Termios
var oldFlags int32

func getTermios() (result syscall.Termios, err error) {
	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCGETS, uintptr(unsafe.Pointer(&result)))
	if errno != 0 {
		return result, os.NewSyscallError("SYS_IOCTL", errno)
	}
	if r1 != 0 {
		return result, fmt.Errorf("Error: expected first syscall result to be 0, got %d", r1)
	}
	return result, nil
}

func setTermios(t syscall.Termios) error {
	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCSETS, uintptr(unsafe.Pointer(&t)))
	if errno != 0 {
		return os.NewSyscallError("SYS_IOCTL", errno)
	}
	if r1 != 0 {
		return fmt.Errorf("Error: expected first syscall result to be 0, got %d", r1)
	}
	return nil
}

func getFileStatusFlags() (int32, error) {
	r1, _, errno := syscall.Syscall(syscall.SYS_FCNTL, uintptr(syscall.Stdin), syscall.F_GETFL, 0)
	if errno != 0 {
		return 0, os.NewSyscallError("SYS_FCNTL", errno)
	}
	r := int32(r1)
	if r < 0 {
		return 0, fmt.Errorf("Error: expected first syscall result to be >= 0, got %d", r)
	}
	return r, nil
}

func setFileStatusFlags(f int32) error {
	r1, _, errno := syscall.Syscall(syscall.SYS_FCNTL, uintptr(syscall.Stdin), syscall.F_SETFL, uintptr(f))
	if errno != 0 {
		return os.NewSyscallError("SYS_FCNTL", errno)
	}
	if r1 != 0 {
		return fmt.Errorf("Error: expected first syscall result to be 0, got %d", r1)
	}
	return nil
}

// http://stackoverflow.com/a/6599441/23582

func readSingleKeypress() ([]byte, error) {
	oldFlags, err := getFileStatusFlags()
	if err != nil {
		return nil, err
	}
	oldTermios, err := getTermios()
	if err != nil {
		return nil, err
	}

	defer setFileStatusFlags(oldFlags)
	defer setTermios(oldTermios)

	newFl, newTermios := oldFlags, oldTermios

	newTermios.Iflag &^= (syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP)
	newTermios.Iflag &^= (syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON)
	newTermios.Oflag &^= syscall.OPOST
	newTermios.Cflag &^= (syscall.CSIZE | syscall.PARENB)
	newTermios.Cflag |= syscall.CS8
	newTermios.Lflag &^= (syscall.ECHONL | syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN)
	newTermios.Cc[syscall.VMIN] = 1
	newTermios.Cc[syscall.VTIME] = 0
	if err = setTermios(newTermios); err != nil {
		return nil, err
	}

	newFl &^= syscall.O_NONBLOCK
	if err = setFileStatusFlags(newFl); err != nil {
		return nil, err
	}

	keys := []byte{0}
	n, err := syscall.Read(syscall.Stdin, keys)
	if err != nil {
		return nil, err
	}
	if n != 1 {
		return nil, fmt.Errorf("Expected to read 1 byte, got %d", n)
	}

	return keys, nil
}