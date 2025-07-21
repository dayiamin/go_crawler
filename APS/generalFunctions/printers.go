package generalFunctions

import (
	"fmt"

	"os"
	"os/exec"
	"runtime"
	"time"
)

// Clear terminal
func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Print timer
func PrintActiveTimer(start time.Time, done chan bool) {
	for {
		select {
		case <-done:
			
			fmt.Println("Timer stopped")
			return
		default:
			elapsed := time.Since(start)
			fmt.Printf("\rCrawl time:%s ", elapsed.Truncate(time.Second))
			time.Sleep(1 * time.Second)
		}
	}
}
