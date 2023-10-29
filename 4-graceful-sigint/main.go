//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	proc := MockProcess{}
	go func() {
		proc.Run()
	}()
	sig := make(chan os.Signal, 1)
	done := make(chan interface{})
	signal.Notify(sig, syscall.SIGINT)
	<-sig
	go func() {
		proc.Stop()
		done <- nil
	}()

	select {
	case <-sig:
		os.Exit(1)
	case <-done:
		os.Exit(0)
	}
}
