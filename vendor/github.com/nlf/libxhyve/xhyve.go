// +build darwin

package xhyve

// #cgo CFLAGS: -I${SRCDIR}/include -x c -std=c11 -fno-common -arch x86_64 -DXHYVE_CONFIG_ASSERT -DVERSION=v0.2.0 -Os -fstrict-aliasing -Wno-unknown-warning-option -Wno-reserved-id-macro -pedantic -fmessage-length=152 -fdiagnostics-show-note-include-stack -fmacro-backtrace-limit=0 -Wno-gnu-zero-variadic-macro-arguments
// #cgo LDFLAGS: -L${SRCDIR} -arch x86_64 -framework Hypervisor -framework vmnet
// #include <xhyve/xhyve.h>
// #include <xhyve/mevent.h>
// #include <string.h>
import "C"
import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

var termios syscall.Termios

// getTermios gets the current settings for the terminal.
func getTermios() syscall.Termios {
	var state syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(0), uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&state)), 0, 0, 0); err != 0 {
		fmt.Fprintln(os.Stderr, err)
	}
	return state
}

// setTermios restores terminal settings.
func setTermios(state syscall.Termios) {
	syscall.Syscall6(syscall.SYS_IOCTL, uintptr(0), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&state)), 0, 0, 0)
}

// go_callback_exit gets invoked from within xhyve.c whenever a trap
// suspending the VM is triggered. This is so we can clean up resources
// in Go land, restore terminal settings and allow the goroutine to be scheduled
// on multiple OS threads again by Go's scheduler.
//export go_callback_exit
func go_callback_exit(status C.int) {
	exitStatus := map[int]string{
		0:   "Reset",
		1:   "PowerOFF",
		2:   "Halt",
		3:   "TripleFault",
		100: "Internal error",
	}

	// Restores stty settings to the values that existed before running xhyve.
	setTermios(termios)

	fmt.Printf("VM has been suspended by %s event\n", exitStatus[int(status)])
	fmt.Printf("Releasing allocated memory from Go land... ")
	for _, arg := range argv {
		C.free(unsafe.Pointer(arg))
	}
	fmt.Println("done")

	// Turns exit flag On for mevent busy loop so that the next time kevent
	// receives an event, mevent handles it and exits the loop.
	fmt.Print("Signaling xhyve mevent dispatch loop to exit... ")
	C.exit_mevent_dispatch_loop = true

	// Forces kevent() to exit by using the self-pipe trick.
	C.mevent_exit()
	fmt.Println("done")

	// Allows Go's scheduler to move the goroutine to a different OS thread.
	runtime.UnlockOSThread()
}

// go_set_pty_name is called by xhyve whenever a master/slave pseudo-terminal is setup in
// COM1 or COM2.
//export go_set_pty_name
func go_set_pty_name(name *C.char) {
	if newPty == nil {
		return
	}
	newPty <- C.GoString(name)
}

func init() {
	// Saves stty settings
	termios = getTermios()

	// We need to stick the goroutine to its current OS thread so Go's scheduler
	// does not move it to a different thread while xhyve is running. By doing this
	// we make sure that once go_callback_exit is invoked from C land, it in turn
	// invokes C funtions from the same OS thread and thread context.
	runtime.LockOSThread()
}

// newPty is a channel to send through the devices path names
// for new devices created when a LPC device is added with the
// option: autopty. Example: -l com1,autopty
// If you add a LPC device on COM2 with autopty enabled, you might need to
// make sure the guest OS runs getty on the pseudo-terminal device created, so
// a login prompt is shown once you open such pseudo-terminal.
var newPty chan string

var argv []*C.char

// Run runs xhyve hypervisor.
func Run(params []string, newPtyCh chan string) error {
	newPty = newPtyCh

	argc := C.int(len(params))
	argv = make([]*C.char, argc)
	for i, arg := range params {
		argv[i] = C.CString(arg)
	}

	// Runs xhyve and blocks.
	if err := C.run_xhyve(argc, &argv[0]); err != 0 {
		fmt.Printf("ERROR => %s\n", C.GoString(C.strerror(err)))
		return fmt.Errorf("Error initializing hypervisor")
	}

	return nil
}
