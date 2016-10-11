package sync

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func executeCmd(cmd string) {
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
		//os.Exit(1) // we continue b/c maybe at the next sync the issue is fixed
	}

	log.Println("done")
}

func Execute(syncCmd string) {
	log.Println("Executing sync cmd:", syncCmd)
	executeCmd(syncCmd)
}
