package docdbuild

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// signalHandler : Handles all signals sent to DoCD
func signalHandler(signalChan chan os.Signal, exitChan chan int) {
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		for {
			s := <-signalChan
			switch s {
			case syscall.SIGHUP:
				exitChan <- 0

			case syscall.SIGINT:
				exitChan <- 0

			case syscall.SIGTERM:
				exitChan <- 0

			case syscall.SIGQUIT:
				exitChan <- 0

			default:
				fmt.Println("Unknown signal.")
				exitChan <- 1
			}
		}
	}()
}
