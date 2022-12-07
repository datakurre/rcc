package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/peercc"
	"github.com/robocorp/rcc/pretty"
)

var (
	domainId    string
	serverName  string
	serverPort  int
	versionFlag bool
	holdingArea string
)

func defaultHoldLocation() string {
	where, err := pathlib.Abs(filepath.Join(os.TempDir(), "peercchold"))
	if err != nil {
		return "temphold"
	}
	return where
}

func init() {
	flag.BoolVar(&common.DebugFlag, "debug", false, "Turn on debugging output.")
	flag.BoolVar(&common.TraceFlag, "trace", false, "Turn on tracing output.")

	flag.BoolVar(&versionFlag, "version", false, "Just show peercc version and exit.")
	flag.StringVar(&serverName, "hostname", "localhost", "Hostname/address to bind server to.")
	flag.IntVar(&serverPort, "port", 4653, "Port to bind server in given hostname.")
	flag.StringVar(&holdingArea, "hold", defaultHoldLocation(), "Directory where to put HOLD files once known.")
	flag.StringVar(&domainId, "domain", "personal", "Symbolic domain that this peer serves.")
}

func ExitProtection() {
	status := recover()
	if status != nil {
		exit, ok := status.(common.ExitCode)
		if ok {
			exit.ShowMessage()
			common.WaitLogs()
			os.Exit(exit.Code)
		}
		common.WaitLogs()
		panic(status)
	}
	common.WaitLogs()
}

func showVersion() {
	common.Stdout("%s\n", common.Version)
	os.Exit(0)
}

func process() {
	if versionFlag {
		showVersion()
	}
	pretty.Guard(common.SharedHolotree, 1, "Shared holotree must be enabled and in use for peercc to work.")
	common.Log("Peer for rcc starting (%s) ...", common.Version)
	peercc.Serve(serverName, serverPort, domainId, holdingArea)
}

func main() {
	defer ExitProtection()
	pretty.Setup()

	flag.Parse()
	common.UnifyVerbosityFlags()
	process()
}
