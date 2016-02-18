package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
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
	tmpDir, err := ioutil.TempDir("", "api-gateway-config-supervisor-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	syncFolder = &tmpDir
	var sync_interval = time.Duration(time.Second * 1)
	syncInterval = &sync_interval

	// start the utility in background
	go main()

	// wait for some time to init
	time.Sleep(500 * time.Millisecond)

	//modifyFS: create a new file
	f1, err := createFile(t, tmpDir, "new-file-content")
	defer os.Remove(f1.Name()) // clean up

	// use a reload command to change the content of the new file added
	// so that we later test its content, expecting to find this change
	reload_cmd := "sed -i.old s/new/reload_cmd/ " + f1.Name()
	reloadCmd = &reload_cmd

	// wait for some time to track the changes
	time.Sleep(1000 * time.Millisecond)

	// check that reload cmd has been applied
	c, err := ioutil.ReadFile(f1.Name())
	if !strings.HasPrefix(string(c), "reload_cmd-file-content") {
		t.Fatal("reload cmd did not run correctly. File content was:" + string(c))
	}

	// wait for some time to check the changes
	time.Sleep(500 * time.Millisecond)
}
