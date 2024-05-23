package websockets

import "syscall"

func InitWebsocketServer() {
	increaseResourcesLimitations()
}

func increaseResourcesLimitations() {
	// increase resources limitations
	var rLimit syscall.Rlimit

	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
}
