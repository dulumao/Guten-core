package graceful

import (
	"fmt"
	"github.com/gookit/color"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(onExit func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go handleSignals(sigChan, onExit)
}

func handleSignals(c chan os.Signal, cb func()) {
	fmt.Printf("Notice: System signal monitoring is enabled(watch: SIGINT,SIGTERM,SIGQUIT)\n")

	switch <-c {
	case syscall.SIGINT:
		color.Info.Println("\nShutdown by Ctrl+C")
	case syscall.SIGTERM: // by kill
		color.Info.Println("\nShutdown quickly")
	case syscall.SIGQUIT:
		color.Info.Println("\nShutdown gracefully")
		// TODO do graceful shutdown
	}

	// if callback exist
	if cb != nil {
		cb()
	}

	// wait some time
	time.Sleep(1e9 / 2)

	color.Info.Println("GoodBye :) ...")
	os.Exit(0)
}
