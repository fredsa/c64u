# c64u

Go CLI for the [Ultimate 1541-II+](https://ultimate64.com/) REST API.

Requires Ultimate firmware 3.11 or later.

## API documentation

https://1541u-documentation.readthedocs.io/en/latest/api/api_calls.html

## Prerequisites

Enable the Web API on your Ultimate device:

1. Enter the Ultimate menu
2. Navigate to **Network Services & Timezone**
3. Set **Web Remote Control Service** to **Enabled**

## Install

```
go install github.com/fredsa/c64u@latest
```

Or build from source:

```
git clone https://github.com/fredsa/c64u.git
cd c64u
go install
```

## CLI Usage

```
c64u <host> <command> [args...]
```

File arguments that refer to a local file are automatically uploaded via POST.
Otherwise, the path is treated as a remote path on the Ultimate's filesystem.

### Commands

#### About
```
c64u 192.168.1.100 version
```

#### Runners
```
c64u 192.168.1.100 sidplay /Usb0/music/song.sid 3
c64u 192.168.1.100 modplay /Usb0/music/track.mod
c64u 192.168.1.100 load /Usb0/games/game.prg
c64u 192.168.1.100 run /Usb0/games/game.prg
c64u 192.168.1.100 run ./local-file.prg
c64u 192.168.1.100 crt /Usb0/carts/cart.crt
```

#### Machine Control
```
c64u 192.168.1.100 reset
c64u 192.168.1.100 reboot
c64u 192.168.1.100 pause
c64u 192.168.1.100 resume
c64u 192.168.1.100 poweroff
c64u 192.168.1.100 writemem D020 0504
c64u 192.168.1.100 readmem 0400 256
c64u 192.168.1.100 debugreg
c64u 192.168.1.100 debugreg FF
```

#### Configuration
```
c64u 192.168.1.100 categories
c64u 192.168.1.100 config "Drive A Settings"
c64u 192.168.1.100 config "Drive A Settings" "Drive Bus ID"
c64u 192.168.1.100 setconfig "Drive A Settings" "Drive Bus ID" 9
c64u 192.168.1.100 saveconfig
c64u 192.168.1.100 loadconfig
c64u 192.168.1.100 resetconfig
```

#### Drives
```
c64u 192.168.1.100 drives
c64u 192.168.1.100 mount a /Usb0/disks/game.d64
c64u 192.168.1.100 mount a /Usb0/disks/game.d64 readonly
c64u 192.168.1.100 unmount a
c64u 192.168.1.100 driveon a
c64u 192.168.1.100 driveoff a
c64u 192.168.1.100 drivemode a 1581
c64u 192.168.1.100 drivereset a
```

#### Streams (U64 only)
```
c64u 192.168.1.100 stream-start video 192.168.1.50
c64u 192.168.1.100 stream-start audio 192.168.1.50:6789
c64u 192.168.1.100 stream-stop video
```

#### Files
```
c64u 192.168.1.100 fileinfo /Usb0/disks/game.d64
c64u 192.168.1.100 create-d64 /Usb0/disks/new.d64
c64u 192.168.1.100 create-d64 /Usb0/disks/new.d64 40 "MY DISK"
c64u 192.168.1.100 create-d71 /Usb0/disks/new.d71
c64u 192.168.1.100 create-d81 /Usb0/disks/new.d81
c64u 192.168.1.100 create-dnp /Usb0/disks/new.dnp 100
```
