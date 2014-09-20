// Debug features for the Fab Lab Locksmith
package main

import (
  "io"
  "log" 
)

var (
  TRACE   *log.Logger
  INFO    *log.Logger
  WARNING *log.Logger
  ERROR   *log.Logger
)

func logInit(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

  /*
   Usage: TRACE.Println(Error Type)
   ERROR.Println("Couldn't prepare the card authenticator\n")
   WARNING.Println("Be careful!\n")
   INFO.Println("Something good just happened!\n")
  */

	TRACE = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	INFO = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	WARNING = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ERROR = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}