package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fredsa/c64u"
)

func main() {
	args := os.Args[1:]

	debug := false
	if len(args) > 0 && (args[0] == "--debug" || args[0] == "-d") {
		debug = true
		args = args[1:]
	}

	if len(args) < 2 {
		usage()
	}

	host := args[0]
	cmd := args[1]
	rest := args[2:]

	client := c64u.NewClient(host)
	client.Debug = debug

	var err error
	switch cmd {
	// About
	case "version":
		err = cmdVersion(client)

	// Runners
	case "sidplay":
		err = cmdSIDPlay(client, rest)
	case "modplay":
		err = cmdMODPlay(client, rest)
	case "load":
		err = cmdLoadPRG(client, rest)
	case "run":
		err = cmdRunPRG(client, rest)
	case "crt":
		err = cmdRunCRT(client, rest)

	// Machine
	case "reset":
		err = cmdSimple(client, client.Reset)
	case "reboot":
		err = cmdSimple(client, client.Reboot)
	case "pause":
		err = cmdSimple(client, client.Pause)
	case "resume":
		err = cmdSimple(client, client.Resume)
	case "poweroff":
		err = cmdSimple(client, client.PowerOff)
	case "writemem":
		err = cmdWriteMem(client, rest)
	case "readmem":
		err = cmdReadMem(client, rest)
	case "debugreg":
		err = cmdDebugReg(client, rest)

	// Config
	case "categories":
		err = cmdCategories(client)
	case "config":
		err = cmdConfig(client, rest)
	case "setconfig":
		err = cmdSetConfig(client, rest)
	case "saveconfig":
		err = cmdSimple(client, client.SaveConfigToFlash)
	case "loadconfig":
		err = cmdSimple(client, client.LoadConfigFromFlash)
	case "resetconfig":
		err = cmdSimple(client, client.ResetConfigToDefault)

	// Drives
	case "drives":
		err = cmdDrives(client)
	case "mount":
		err = cmdMount(client, rest)
	case "unmount":
		err = cmdUnmount(client, rest)
	case "driveon":
		err = cmdDriveOnOff(client, rest, true)
	case "driveoff":
		err = cmdDriveOnOff(client, rest, false)
	case "drivemode":
		err = cmdDriveMode(client, rest)
	case "drivereset":
		err = cmdDriveReset(client, rest)
	case "loadrom":
		err = cmdLoadROM(client, rest)

	// Streams
	case "stream-start":
		err = cmdStreamStart(client, rest)
	case "stream-stop":
		err = cmdStreamStop(client, rest)

	// Files
	case "fileinfo":
		err = cmdFileInfo(client, rest)
	case "create-d64":
		err = cmdCreateD64(client, rest)
	case "create-d71":
		err = cmdCreateD71(client, rest)
	case "create-d81":
		err = cmdCreateD81(client, rest)
	case "create-dnp":
		err = cmdCreateDNP(client, rest)

	default:
		fatalf("Unknown command: %s\n", cmd)
	}

	if err != nil {
		fatalf("Error: %v\n", err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: c64u [-d] <host> <command> [args...]

Commands:
  About:
    version                         Show API version

  Runners:
    sidplay <file> [songnr]         Play a SID file
    modplay <file>                  Play a MOD file
    load <file>                     Load a PRG (no auto-run)
    run <file>                      Load and run a PRG
    crt <file>                      Start a CRT cartridge file

  Machine:
    reset                           Reset the machine
    reboot                          Reboot the machine
    pause                           Pause the CPU
    resume                          Resume the CPU
    poweroff                        Power off (U64 only)
    writemem <addr> <hexdata|file>  Write hex data or file to memory
    readmem <addr> [length]         Read memory (default 256 bytes)
    debugreg [value]                Read/write debug register (U64)

  Configuration:
    categories                      List config categories
    config [category [item]]        Get config values
    setconfig <category> <item> <value>  Set a config value
    saveconfig                      Save config to flash
    loadconfig                      Load config from flash
    resetconfig                     Reset config to defaults

  Drives:
    drives                          List drive info
    mount <drive> <image> [mode]    Mount a disk image (d64|g64|d71|g71|d81; mode: readwrite|readonly|unlinked)
    unmount <drive>                 Remove disk from drive
    driveon <drive>                 Turn drive on
    driveoff <drive>                Turn drive off
    drivemode <drive> <mode>        Set drive mode (1541|1571|1581)
    drivereset <drive>              Reset drive
    loadrom <drive> <file>          Load a drive ROM (16K or 32K)

  Streams (U64 only):
    stream-start <name> <ip[:port]> Start a stream (video|audio|debug)
    stream-stop <name>              Stop a stream

  Files:
    fileinfo <path>                 Get file info
    create-d64 <path> [tracks] [diskname]  Create a .d64 image
    create-d71 <path> [diskname]    Create a .d71 image
    create-d81 <path> [diskname]    Create a .d81 image
    create-dnp <path> <tracks> [diskname]  Create a .dnp image

Examples:
    c64u 192.168.1.100 version
    c64u 192.168.1.100 run /Usb0/games/commando.prg
    c64u 192.168.1.100 writemem D020 05
    c64u 192.168.1.100 mount a /Usb0/disks/game.d64 readonly
`)
	os.Exit(1)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func requireArgs(args []string, min int, usage string) {
	if len(args) < min {
		fatalf("Usage: %s\n", usage)
	}
}

// putOrPost checks if the path is a local file. If so, it reads and POSTs it.
// Otherwise it treats it as a remote path on the Ultimate and PUTs it.
func putOrPost(path string, putFn func(string) (*c64u.ErrorResponse, error), postFn func([]byte) (*c64u.ErrorResponse, error)) (*c64u.ErrorResponse, error) {
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return postFn(data)
	}
	return putFn(path)
}

func printJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func printErrors(resp *c64u.ErrorResponse) {
	if resp != nil && len(resp.Errors) > 0 {
		for _, e := range resp.Errors {
			fmt.Fprintf(os.Stderr, "Error: %s\n", e)
		}
	}
}

func cmdSimple(_ *c64u.Client, fn func() (*c64u.ErrorResponse, error)) error {
	resp, err := fn()
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdVersion(client *c64u.Client) error {
	resp, err := client.Version()
	if err != nil {
		return err
	}
	printJSON(resp)
	return nil
}

func cmdSIDPlay(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> sidplay <file> [songnr]")
	songNr := 0
	if len(args) > 1 {
		n, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid song number: %s", args[1])
		}
		songNr = n
	}
	var resp *c64u.ErrorResponse
	if _, statErr := os.Stat(args[0]); statErr == nil {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}
		resp, err = client.SIDPlayData(data, songNr, nil)
		if err != nil {
			return err
		}
	} else {
		var err error
		resp, err = client.SIDPlay(args[0], songNr)
		if err != nil {
			return err
		}
	}
	printErrors(resp)
	return nil
}

func cmdMODPlay(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> modplay <file>")
	resp, err := putOrPost(args[0], client.MODPlay, client.MODPlayData)
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdLoadPRG(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> load <file>")
	resp, err := putOrPost(args[0], client.LoadPRG, client.LoadPRGData)
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdRunPRG(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> run <file>")
	resp, err := putOrPost(args[0], client.RunPRG, client.RunPRGData)
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdRunCRT(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> crt <file>")
	resp, err := putOrPost(args[0], client.RunCRT, client.RunCRTData)
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdWriteMem(client *c64u.Client, args []string) error {
	requireArgs(args, 2, "c64u <host> writemem <addr> <hexdata | file>")
	addr := args[0]
	var resp *c64u.ErrorResponse
	var err error
	if _, statErr := os.Stat(args[1]); statErr == nil {
		data, readErr := os.ReadFile(args[1])
		if readErr != nil {
			return readErr
		}
		resp, err = client.WriteMemData(addr, data)
	} else {
		resp, err = client.WriteMem(addr, args[1])
	}
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdReadMem(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> readmem <addr> [length]")
	length := 0
	if len(args) > 1 {
		n, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid length: %s", args[1])
		}
		length = n
	}
	data, err := client.ReadMem(args[0], length)
	if err != nil {
		return err
	}
	// Print as hex dump
	for i, b := range data {
		if i > 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X ", b)
	}
	fmt.Println()
	return nil
}

func cmdDebugReg(client *c64u.Client, args []string) error {
	if len(args) > 0 {
		resp, err := client.WriteDebugReg(args[0])
		if err != nil {
			return err
		}
		printJSON(resp)
	} else {
		resp, err := client.ReadDebugReg()
		if err != nil {
			return err
		}
		printJSON(resp)
	}
	return nil
}

func cmdCategories(client *c64u.Client) error {
	resp, err := client.ListCategories()
	if err != nil {
		return err
	}
	printJSON(resp)
	return nil
}

func cmdConfig(client *c64u.Client, args []string) error {
	if len(args) == 0 {
		resp, err := client.ListCategories()
		if err != nil {
			return err
		}
		printJSON(resp)
	} else if len(args) == 1 {
		resp, err := client.GetConfig(args[0])
		if err != nil {
			return err
		}
		printJSON(resp)
	} else {
		resp, err := client.GetConfigItem(args[0], args[1])
		if err != nil {
			return err
		}
		printJSON(resp)
	}
	return nil
}

func cmdSetConfig(client *c64u.Client, args []string) error {
	requireArgs(args, 3, "c64u <host> setconfig <category> <item> <value>")
	resp, err := client.SetConfigItem(args[0], args[1], strings.Join(args[2:], " "))
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdDrives(client *c64u.Client) error {
	resp, err := client.ListDrives()
	if err != nil {
		return err
	}
	printJSON(resp)
	return nil
}

func cmdMount(client *c64u.Client, args []string) error {
	requireArgs(args, 2, "c64u <host> mount <drive> <image> [mode]")
	var mode c64u.MountMode
	if len(args) > 2 {
		mode = c64u.MountMode(args[2])
	}
	drive := args[0]
	image := args[1]
	var resp *c64u.ErrorResponse
	var err error
	if _, statErr := os.Stat(image); statErr == nil {
		data, readErr := os.ReadFile(image)
		if readErr != nil {
			return readErr
		}
		resp, err = client.MountImageData(drive, data, filepath.Base(image), "", mode)
	} else {
		resp, err = client.MountImage(drive, image, "", mode)
	}
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdUnmount(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> unmount <drive>")
	resp, err := client.RemoveDisk(args[0])
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdDriveOnOff(client *c64u.Client, args []string, on bool) error {
	requireArgs(args, 1, "c64u <host> driveon|driveoff <drive>")
	var resp *c64u.ErrorResponse
	var err error
	if on {
		resp, err = client.DriveOn(args[0])
	} else {
		resp, err = client.DriveOff(args[0])
	}
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdDriveMode(client *c64u.Client, args []string) error {
	requireArgs(args, 2, "c64u <host> drivemode <drive> <mode>")
	resp, err := client.SetDriveMode(args[0], c64u.DriveMode(args[1]))
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdDriveReset(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> drivereset <drive>")
	resp, err := client.ResetDrive(args[0])
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdLoadROM(client *c64u.Client, args []string) error {
	requireArgs(args, 2, "c64u <host> loadrom <drive> <file>")
	drive := args[0]
	file := args[1]
	var resp *c64u.ErrorResponse
	var err error
	if _, statErr := os.Stat(file); statErr == nil {
		data, readErr := os.ReadFile(file)
		if readErr != nil {
			return readErr
		}
		resp, err = client.LoadDriveROMData(drive, data)
	} else {
		resp, err = client.LoadDriveROM(drive, file)
	}
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdStreamStart(client *c64u.Client, args []string) error {
	requireArgs(args, 2, "c64u <host> stream-start <name> <ip[:port]>")
	resp, err := client.StartStream(c64u.StreamName(args[0]), args[1])
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdStreamStop(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> stream-stop <name>")
	resp, err := client.StopStream(c64u.StreamName(args[0]))
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdFileInfo(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> fileinfo <path>")
	resp, err := client.FileInfo(args[0])
	if err != nil {
		return err
	}
	printJSON(resp)
	return nil
}

func cmdCreateD64(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> create-d64 <path> [tracks] [diskname]")
	tracks := 0
	diskName := ""
	if len(args) > 1 {
		n, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid tracks: %s", args[1])
		}
		tracks = n
	}
	if len(args) > 2 {
		diskName = args[2]
	}
	resp, err := client.CreateD64(args[0], tracks, diskName)
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdCreateD71(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> create-d71 <path> [diskname]")
	diskName := ""
	if len(args) > 1 {
		diskName = args[1]
	}
	resp, err := client.CreateD71(args[0], diskName)
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdCreateD81(client *c64u.Client, args []string) error {
	requireArgs(args, 1, "c64u <host> create-d81 <path> [diskname]")
	diskName := ""
	if len(args) > 1 {
		diskName = args[1]
	}
	resp, err := client.CreateD81(args[0], diskName)
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}

func cmdCreateDNP(client *c64u.Client, args []string) error {
	requireArgs(args, 2, "c64u <host> create-dnp <path> <tracks> [diskname]")
	tracks, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid tracks: %s", args[1])
	}
	diskName := ""
	if len(args) > 2 {
		diskName = args[2]
	}
	resp, err := client.CreateDNP(args[0], tracks, diskName)
	if err != nil {
		return err
	}
	printErrors(resp)
	return nil
}
