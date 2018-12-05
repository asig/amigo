# amigo

`amigo` is a program that allows you to control an Arduino from 
Minecraft, heavily inspired by https://www.instructables.com/id/Arduino-Minecraft-Interface/
The main reason for writing this was that my son got a "Programmieren mit 
Minecraft" advent calendar, but the shipped software didn't work...

# Build
```bash
go get github.com/tarm/serial
go build
```

# Installation
Just copy the binary to somewhere on your disk. Feel free to copy it into the 
logs directory as described on instructables.com

# Usage
`amigo [-logfile <path/to/latest.log>] [-serial <com-port>] [-baud <baudrate>] [-verbose]`


## Flags

| Flag         | Meaning                                                                  |
|--------------|--------------------------------------------------------------------------|
| `-logfile`   | Path to Minecraft's logfile. Default is 'latest.log'                     |
| `-serial`    | COM-Port to use, e.g. COM3                                               |
| `-baud`      | Baud rate. Default is 9600                                               |
| `-verbose`   | If set, `amigo` is quite chatty and lets you know whatever it is doing   |

Example: `amigo -logfile %APPDATA%\.minecraft\logs\latest.log -serial COM3 -baud 9600`