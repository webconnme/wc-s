package main


import (
	"wc/ioctl"
)
import (
	"os"
	"log"
	"reflect"
	"unsafe"
)

const (
	IOCTL_WRITE = 0x40044900
	IOCTL_READ = 0x80044901
)

func main() {
	file, err := os.OpenFile("/dev/webconn", os.O_RDWR | os.O_SYNC, 0777)
	if err != nil {
		log.Fatal("open", err)
	}
	log.Println("open");
	defer file.Close()


	write_buffer := []byte("test")
	n, err := file.Write(write_buffer)
	if err != nil {
		log.Fatal("write", err)
	}
	log.Printf("write: %v\n", string(write_buffer[:n]))

	read_buffer := make([]byte, 128)
	n, err = file.Read(read_buffer)
	if err != nil {
		log.Fatal("read", err)
	}
	log.Printf("read: %v\n", string(read_buffer[:n]))



	ioctl_buffer := 0
	header := (*reflect.SliceHeader)(unsafe.Pointer(&ioctl_buffer))

	ioctl.IOCTL(uintptr(file.Fd()), ioctl.IOW('I', 0, 4), 0)
	ioctl.IOCTL(uintptr(file.Fd()), ioctl.IOR('I', 1, 4), header.Data)

	log.Printf("ioctl_read: %v\n", ioctl_buffer)

/*
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), IOCTL_WRITE, 0)
	if ep != 0 {
		//syscall.Errno(ep)
	}

	_, _, ep = syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), IOCTL_READ, 0)
	if ep != 0 {
		//syscall.Errno(ep)
	}
*/

	log.Println("close");
}
