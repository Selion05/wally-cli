package main

import (
	"fmt"
	"github.com/caarlos0/spin"
	"gopkg.in/cheggaaa/pb.v1"
	"os"
	"path/filepath"
	"time"
)

var appVersion = "2.0.0"

type state struct {
	step  int
	total int
	sent  int
}

func main() {
	var args = os.Args[1:]
	s := state{step: 0, total: 0, sent: 0}

	if len(args) != 1 {
		fmt.Println("Usage: wally-cli <firmware file>")
		os.Exit(2)
	}

	if args[0] == "--version" {
		appVersion := fmt.Sprintf("wally-cli v%s", appVersion)
		fmt.Println(appVersion)
		os.Exit(0)
	}

	path := args[0]
	extension := filepath.Ext(path)
	if extension != ".bin" && extension != ".hex" {
		fmt.Println("Please provide a valid firmware file: a .hex file (ErgoDox EZ) or a .bin file (Moonlander / Planck EZ)")
		os.Exit(2)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("The file path you specified does not exist")
		os.Exit(1)
	}

	spinner := spin.New("%s Press the reset button of your keyboard.")
	spinner.Start()
	spinnerStopped := false

	var progress *pb.ProgressBar
	progressStarted := false

	if extension == ".bin" {
		go dfuFlash(path, &s)
	}
	if extension == ".hex" {
		go teensyFlash(path, &s)
	}

	for s.step != 2 {
		time.Sleep(500 * time.Millisecond)
		if s.step > 0 {
			if spinnerStopped == false {
				spinner.Stop()
				spinnerStopped = true
			}
			if progressStarted == false {
				progressStarted = true
				progress = pb.StartNew(s.total)
			}
			progress.Set(s.sent)
		}
	}
	progress.Finish()
	fmt.Println("Your keyboard was successfully flashed and rebooted. Enjoy the new firmware!")
}
