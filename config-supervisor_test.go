package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	Sync "sync"
	"testing"
	"time"

	"github.com/adobe-apiplatform/api-gateway-config-supervisor/sync"
	"github.com/koyachi/go-term-ansicolor/ansicolor"
)

var once Sync.Once

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func setup(t *testing.T) (tempdir string) {
	tmpDir, err := ioutil.TempDir("", "api-gateway-config-supervisor-")
	if err != nil {
		t.Fatal(err)
	}

	// setup a syncFolder to watch by the main program
	syncFolder = &tmpDir
	// setup a sync interval for the test
	var sync_interval = time.Duration(time.Second * 1)
	syncInterval = &sync_interval
	// setup extra debug information
	debugOn := true
	debug = &debugOn

	once.Do(func() {
		log.Println(ansicolor.IntenseBlack("Starting the main() program"))
		go main()
	})

	return tmpDir
}

func createFile(t *testing.T, tmpDir string, file_content string) (f *os.File, err error) {
	content := []byte(file_content)
	tmpfile, err := ioutil.TempFile(tmpDir, "new-file-")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpfile, err
}

func TestThatReloadCommandExecutesOnFsChanges(t *testing.T) {
	tmpDir := setup(t)
	defer os.RemoveAll(tmpDir)

	status = sync.GetStatusInstance()
	if time.Since(status.LastSync).Seconds() < 2 {
		t.Fatal("sync should not happen immediately" + time.Since(status.LastSync).String())
	}
	if time.Since(status.LastReload).Seconds() > 2 {
		t.Fatal("LastReload should be current but was " + time.Since(status.LastReload).String())
	}
	if time.Since(status.LastFSChangeDetected).Seconds() > 2 {
		t.Fatal("LastFSChangeDetected should be current but was " + time.Since(status.LastFSChangeDetected).String())
	}

	// wait for some time to init
	time.Sleep(500 * time.Millisecond)

	//modifyFS: create a new directory and file
	dir, err := ioutil.TempDir(tmpDir, "new-dir-")
	if err != nil {
		t.Fatal(err)
	}
	f1, err := createFile(t, dir, "new-file-content")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f1.Name()) // clean up

	// use a reload command to change the content of the new file added
	// so that we later test its content, expecting to find this change
	reload_cmd := "sed -i.old s/new/reload_cmd/ " + f1.Name()
	reloadCmd = &reload_cmd

	// wait for some time to track the changes
	time.Sleep(1000 * time.Millisecond)

	// check that reload cmd has been executed
	if time.Since(status.LastSync).Seconds() > 1 {
		t.Fatal("sync should have executed earlier than 1.5 but was executed " + time.Since(status.LastSync).String())
	}
	if time.Since(status.LastReload).Seconds() > 1 {
		t.Fatal("reload should have executed earlier than 1.5s but was executed " + time.Since(status.LastReload).String())
	}
	if time.Since(status.LastFSChangeDetected).Seconds() > 1 {
		t.Fatal("FS changes should have been detected earlier than 1.5s but was detected " + time.Since(status.LastFSChangeDetected).String())
	}

	c, err := ioutil.ReadFile(f1.Name())
	if !strings.HasPrefix(string(c), "reload_cmd-file-content") {
		t.Fatal("reload cmd did not run correctly. File content was:" + string(c))
	}

	//reset the reload command
	reload_cmd = "echo reload-cmd not defined"
}

func TestThatSyncCommandExecutes(t *testing.T) {
	tmpDir := setup(t)
	defer os.RemoveAll(tmpDir)

	// in order to test that the sync command executed we create a file to later verify that it exists on the disk
	sync_cmd_test := "touch " + tmpDir + "/sync_cmd.txt"
	syncCmd = &sync_cmd_test

	// wait for some time to init
	time.Sleep(1100 * time.Millisecond)

	// check that the sync_cmd.txt file was created when the sync command executes
	if _, err := os.Stat(tmpDir + "/sync_cmd.txt"); err != nil {
		t.Fatal("Expected to find the file created by the sync command " + tmpDir + "/sync_cmd.txt")
	}

	sync_cmd_test = "echo sync-cmd not defined"
}

func TestThatReloadCommandExecutesWhenNewFolderIsAdded(t *testing.T) {
	tmpDir := setup(t)
	defer os.RemoveAll(tmpDir)

	// in order to test that the reload command executed, we create a file and later verify that it exists on the disk
	reload_cmd_test := "touch " + tmpDir + "/reload_cmd.txt"
	reloadCmd = &reload_cmd_test

	syncFolder = &tmpDir
	go watchForFSChanges()

	// wait for some time
	time.Sleep(500 * time.Millisecond)

	//modifyFS: create a new directory and file
	dir, err := ioutil.TempDir(tmpDir, "new-folder-")
	if err != nil {
		t.Fatal(err)
	}
	f2, err := createFile(t, dir, "new-file-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f2.Name()) // clean up

	// wait for some time to track the changes
	time.Sleep(1000 * time.Millisecond)

	status = sync.GetStatusInstance()
	// check that reload cmd has been executed
	if time.Since(status.LastSync).Seconds() > 1 {
		t.Fatal("sync should have executed earlier than 1.5 but was executed " + time.Since(status.LastSync).String())
	}
	if time.Since(status.LastReload).Seconds() > 1 {
		t.Fatal("reload should have executed earlier than 1.5s but was executed " + time.Since(status.LastReload).String())
	}
	if time.Since(status.LastFSChangeDetected).Seconds() > 1 {
		t.Fatal("FS changes should have been detected earlier than 1.5s but was detected " + time.Since(status.LastFSChangeDetected).String())
	}

	// check that the reload_cmd.txt file was created when the reload command executed
	if _, err := os.Stat(tmpDir + "/reload_cmd.txt"); err != nil {
		t.Fatal("Expected to find the file created by the reload command " + tmpDir + "/reload_cmd.txt")
	}

	//delete the tmp file and the TempDir
	reload_cmd_test = "touch " + tmpDir + "/reload_cmd_after_rm.txt"
	reloadCmd = &reload_cmd_test
	// os.Remove(tmpDir + "/reload_cmd.txt")
	os.RemoveAll(dir)

	time.Sleep(3000 * time.Millisecond)

	// check that the reload_cmd_after_rm.txt file was created when the reload command executed
	if _, err := os.Stat(tmpDir + "/reload_cmd_after_rm.txt"); err != nil {
		t.Fatal("Expected to find the file created by the reload command " + tmpDir + "/reload_cmd_after_rm.txt")
	}

	//reset the reload command
	reload_cmd_test = "echo reload-cmd not defined"
}

func TestThatSyncCommandDoesntRunInParallel(t *testing.T) {
	tmpDir := setup(t)
	defer os.RemoveAll(tmpDir)

	// in order to test that the sync command executed we create a file to later verify that it exists on the disk
	sync_cmd_test := "sleep 3"
	syncCmd = &sync_cmd_test

	lastsync := status.LastSync

	time.Sleep(1500 * time.Millisecond)

	if status.LastSyncDuration != -1 {
		t.Fatal("After 1.5s the sync command should not have finished.LastSyncDuration should have been -1")
	}

	if status.LastSync.Equal(lastsync) {
		t.Fatal("Sync shouldn't have executed in the meanwhile.")
	}

	time.Sleep(1700 * time.Millisecond)

	if time.Since(status.LastSync) < 3 {
		t.Fatal("time since last sync is wrong: " + time.Since(status.LastSync).String())
	}

	sync_cmd_test = "echo sync-cmd not defined"
}
