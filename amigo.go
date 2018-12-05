/*
 * amigo: An "Arduino Minecraft Interface" inspired by
 * https://www.instructables.com/id/Arduino-Minecraft-Interface/
 *
 * Copyright (c) 2018 Andreas Signer <asigner@gmail.com>
 *
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/tarm/serial"
)

var (
	logfileFlag = flag.String("logfile", "latest.log", "Path to Minecraft's logfile.")
	serialFlag  = flag.String("serial", "", "Arduino's serial port, e.g. COM3")
	baudFlag    = flag.Int("baud", 9600, "Baud rate for communicating with Arduino.")
	verboseFlag = flag.Bool("verbose", false, "If true, lots of debugging messages will be printed.")

	chatRegexp = regexp.MustCompile(`\[\d\d:\d\d:\d\d\] +\[[^]]*\]: +\[[^]]*\] +\[[^]]*\] +!([^!]+)!.*$`)
)

func tail(file *os.File, lines chan string) {
	// First, skip to the end of the file
	if _, err := file.Seek(0, 2); err != nil {
		log.Fatalf("Can't skip to the end of the file: %s", err)
	}

	// Now, keep reading lines
	curLine := ""
	buf := make([]byte, 1)
	for {
		read, err := file.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("Can't read from logfile: %s", err)
		}
		if read != 0 {
			if buf[0] == '\r' || buf[0] == '\n' {
				if *verboseFlag {
					log.Printf("Read line %q from logfile", curLine)
				}
				lines <- curLine
				curLine = ""
			} else {
				curLine = curLine + string(buf[0])
			}
		}
	}
}

func process(ser *serial.Port, lines chan string) {
	for {
		line := <-lines
		if *verboseFlag {
			log.Printf("Received line %q to process", line)
		}
		m := chatRegexp.FindStringSubmatch(line)
		if len(m) != 2 {
			if *verboseFlag {
				log.Printf("Line is not a CHAT line, ignoring it")
			}
		}
		if len(m) == 2 {
			msg := m[1]
			if *verboseFlag {
				log.Printf("Sending message %q to serial port", msg)
			}
			ser.Write([]byte(msg + "\n"))
		}
	}
}

func main() {
	flag.Parse()

	if *verboseFlag {
		log.Printf("Using logfile %q", *logfileFlag)
		log.Printf("Using serial port %q at %d baud.", *serialFlag, *baudFlag)
	}

	logfile, err := os.Open(*logfileFlag)
	if err != nil {
		log.Fatalf("Can't open log file %q: %s", *logfileFlag, err)
	}

	config := &serial.Config{Name: *serialFlag, Baud: *baudFlag}
	ser, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("Can't open serial port %q: %s", *serialFlag, err)
	}

	lines := make(chan string, 100)
	go tail(logfile, lines)
	process(ser, lines)
}
