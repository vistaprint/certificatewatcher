package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {

	stopChan := make(chan struct{}, 1)
	go handleSigterm(stopChan)

	cfg := NewConfig()

	//Parse the flags
	err := cfg.ParseFlags(os.Args[1:])

	if err != nil {
		log.Fatalf("flag parsing error: %v", err)
	}

	ctrl, err := NewController(cfg)

	if err != nil {
		log.Fatalf("could not create controller %v", err)
	}

	if cfg.Once {
		err := ctrl.RunOnce()
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}
	ctrl.Run(stopChan)

}

func handleSigterm(stopChan chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	log.Info("Received SIGTERM. Terminating...")
	close(stopChan)
}
