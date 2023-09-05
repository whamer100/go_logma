package logma

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

var verbosity = false

const (
	colorReset = "\033[0m"
	colorInfo  = colorReset
	colorWarn  = "\033[33m"
	colorFatal = "\033[31m"
	styleDim   = "\033[90m"
	colorDebug = "\033[2m" + styleDim
	resetAll   = colorReset + "\033[0m"
)

func Info(format string, a ...any) {
	fmt.Printf("%s[INFO] %s%s\n", colorInfo, fmt.Sprintf(format, a...), resetAll)
}

func Warn(format string, a ...any) {
	fmt.Printf("%s[WARN] %s%s\n", colorWarn, fmt.Sprintf(format, a...), resetAll)
}

func Fatal(format string, a ...any) {
	fmt.Printf("%s[FATAL] %s%s\n", colorFatal, fmt.Sprintf(format, a...), resetAll)
}

func Debug(format string, a ...any) {
	if verbosity {
		fmt.Printf("%s[DEBUG] %s%s\n", colorDebug, fmt.Sprintf(format, a...), resetAll)
	}
}

func DebugRaw(format string, a ...any) {
	if verbosity {
		fmt.Printf("%s%s%s\n", colorDebug, fmt.Sprintf(format, a...), resetAll)
	}
}

func SetVerbosity(v bool) {
	verbosity = v
	if verbosity {
		Debug("Debug logging enabled.")
	}
}

func Init() {
	// TODO: other operating systems if needed.
	if strings.Contains(runtime.GOOS, "windows") {
		var (
			kernel32                           = syscall.MustLoadDLL("Kernel32.dll")
			procSetConsoleMode                 = kernel32.MustFindProc("SetConsoleMode")
			procGetConsoleMode                 = kernel32.MustFindProc("GetConsoleMode")
			procGetStdHandle                   = kernel32.MustFindProc("GetStdHandle")
			ENABLE_VIRTUAL_TERMINAL_PROCESSING = uintptr(0x0004)
			E_ERROR                            = uintptr(0)          // nonzero = success
			INVALID_HANDLE_VALUE               = uintptr(4294967295) // uint32(-1)
			STD_OUTPUT_HANDLE                  = uintptr(4294967285) // uint32(-11)
			hOut                               uintptr               // HANDLE
			mode                               uintptr               // DWORD
			result                             uintptr               // BOOL
		)
		// This is literally the fastest way I can do this.
		// obtain console handle
		hOut, _, _ = syscall.SyscallN(procGetStdHandle.Addr(), STD_OUTPUT_HANDLE)
		if hOut == INVALID_HANDLE_VALUE {
			log.Fatalln("Failed to enable ANSI support. (Could not obtain console handle.)")
		}
		// get console mode
		result, _, _ = syscall.SyscallN(procGetConsoleMode.Addr(), hOut, uintptr(unsafe.Pointer(&mode)))
		if result == E_ERROR {
			log.Fatalln("Failed to enable ANSI support. (Could not get console mode.)")
		}
		// set console mode with virtual terminal processing enabled
		mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
		result, _, _ = syscall.SyscallN(procSetConsoleMode.Addr(), hOut, mode)
		if result == E_ERROR {
			log.Fatalln("Failed to enable ANSI support. (Could not set console mode.)")
		}
	}
}
