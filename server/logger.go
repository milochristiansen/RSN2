/*
Copyright 2020-2021 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package main

import "os"
import "io"
import "log"
import "time"

import "github.com/teris-io/shortid"

var (
	infoLog  io.Writer
	warnLog  io.Writer
	errorLog io.Writer
)

var logIDService <-chan string

func init() {
	err := os.MkdirAll("./logs", 0775)
	if err != nil {
		panic("Logger initialization failed. *shrug* Guess I'll die.\n" + err.Error())
	}

	f, err := os.Create("./logs/" + time.Now().UTC().Format("m01-d02-t150405") + ".log")
	if err != nil {
		panic("Logger initialization failed. *shrug* Guess I'll die.\n" + err.Error())
	}

	infoLog = io.MultiWriter(f, os.Stdout)
	warnLog = io.MultiWriter(f, os.Stdout)
	errorLog = io.MultiWriter(f, os.Stderr)

	ml = &SessionLogger{
		I: log.New(infoLog, "INFO@master: ", log.Ldate|log.Ltime|log.Lshortfile),
		W: log.New(warnLog, "WARNING@master: ", log.Ldate|log.Ltime|log.Lshortfile),
		E: log.New(errorLog, "ERROR@master: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	go func() {
		c := make(chan string)
		logIDService = c

		idsource := shortid.MustNew(16, shortid.DefaultABC, uint64(time.Now().UnixNano()))

		for {
			c <- idsource.MustGenerate()
		}
	}()
}

type SessionLogger struct {
	I, W, E *log.Logger
}

func newSessionLogger(endpoint string) *SessionLogger {
	id := <-logIDService
	return &SessionLogger{
		I: log.New(infoLog, "INFO@"+endpoint+":"+id+": ", log.Ldate|log.Ltime|log.Lshortfile),
		W: log.New(warnLog, "WARNING@"+endpoint+":"+id+": ", log.Ldate|log.Ltime|log.Lshortfile),
		E: log.New(errorLog, "ERROR@"+endpoint+":"+id+": ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

var ml *SessionLogger
