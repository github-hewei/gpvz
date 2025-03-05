package main

import (
	"gpvz/gpvz"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	tools := gpvz.NewGPvz()
	go tools.Run()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	tools.Quit()
}
