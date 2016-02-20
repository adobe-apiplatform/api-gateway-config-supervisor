package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/adobe-apiplatform/api-gateway-config-supervisor/sync"
	"github.com/adobe-apiplatform/api-gateway-config-supervisor/ws"
	_ "net/http/pprof"

	"github.com/carlescere/scheduler"

	"github.com/spf13/pflag"

	"github.com/koyachi/go-term-ansicolor/ansicolor"
)

var (
	// Flags
	cpuprofile   = pflag.StringP("cpuprofile", "", "", "Write cpu profile to file")
	version      = pflag.BoolP("version", "V", false, "Print the version number")
	syncInterval = pflag.DurationP("sync-interval", "", time.Second*5, "Time interval for the next sync")
	syncCmd      = pflag.StringP("sync-cmd", "", "echo sync-cmd not defined", "Command used to syncing")
	syncFolder   = pflag.StringP("sync-folder", "", "~/tmp/api-gateway-config", "The folder to watch for changes.")
	reloadCmd    = pflag.StringP("reload-cmd", "", "echo reload-cmd not defined", "Command used to reload the gateway")
	httpAddr     = pflag.StringP("http-addr", "", "127.0.0.1:8888", "Http Address exposing a /health-check for the sync process")
	// when was the reload cmd executed last time
	lastReload = time.Now()
	// when did the last change occur
	lastChange = time.Now()
)

func syntaxError() {
	fmt.Fprintf(os.Stderr, `Execute a sync command and watch a folder for changes.`)
}

// ParseFlags parses the command line flags
func ParseFlags() {
	pflag.Usage = syntaxError
	pflag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Setup profiling if desired
	if *cpuprofile != "" {
		log.Println(ansicolor.Red("Starting CPU Profiling"))
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}
}

func executeSyncCmd() {
	go sync.Execute(*syncCmd)
}

func executeReloadCmd() {
	log.Println(ansicolor.Red("Executing Reload Cmd"))
	go sync.Execute(*reloadCmd)
	lastReload = time.Now()
}

func checkForReload() {
	if time.Since(lastChange) < time.Since(lastReload) && time.Since(lastReload) > *syncInterval {
		lastReload = time.Now()
		executeReloadCmd()
	}
}

//watches for changes in the syncFolder, execute reloadCmd when there are changes
func watchForFSChanges() {
	c := sync.WatchFolderRecursive(*syncFolder)
	for {
		select {
		case file := <-c:
			if file == "" {
				continue
			}
			lastChange = time.Now()
			if time.Since(lastReload) > *syncInterval {
				lastReload = time.Now()
				go func() {
					// wait a little in case there are more changes to sync
					for time.Since(lastChange) < time.Second*1 {
						time.Sleep(1 * time.Second)
					}
					executeReloadCmd()
				}()
			}
		}
	}
}

func startWebServer() {
	go ws.RunWS(*httpAddr)
}

func startWatchingFS() {
	go watchForFSChanges()
	scheduler.Every(int(syncInterval.Seconds())).Seconds().Run(executeSyncCmd)
	scheduler.Every(int(syncInterval.Seconds())).Seconds().Run(checkForReload)
}

func waitToTerminate() {
	// Waiting for terminating (i use a sighandler like in vitess)
	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
}

func main() {
	ParseFlags()
	if *version {
		fmt.Printf("config-supervisor %s\n", "0.1")
		os.Exit(0)
	}

	startWebServer()    // expose a /health-check API on the localhost
	startWatchingFS()	// watch for file system changes
	waitToTerminate()	// wait until the program terminates

	log.Printf("Stopped ... ")
}
