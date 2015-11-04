package sync

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func executeCmd(cmd string, wg *sync.WaitGroup) {
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	cmdRunner := exec.Command(head, parts...)
	cmdReader, err := cmdRunner.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StderrPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("sync stderr | %s\n", scanner.Text())
		}
	}()

	stdOutReader, err := cmdRunner.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scannerOut := bufio.NewScanner(stdOutReader)
	go func() {
		for scannerOut.Scan() {
			fmt.Printf("sync stdout | %s\n", scannerOut.Text())
		}
	}()

	err = cmdRunner.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmdRunner.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		// TODO: exit after a number of retries only
		os.Exit(1)
	}

	log.Println("done")
	wg.Done() // Need to signal to waitgroup that this goroutine is done
}

func Execute(syncCmd string) {
	log.Println("Executing sync cmd:", syncCmd)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	executeCmd(syncCmd, wg)
	wg.Wait()
}
