package main

import (
	"fmt"
)

/*
support resume from break point need server:
 writer.Header().Set("Accept-Ranges","bytes")
 writer.Header().Set("Content-Length",YourFileSize)
 writer.Header().Set("Content-Disposition","attachment; filename=YourFileName")

so we check respone.Header if  "Accept-Ranges" == "bytes" know whether
server support resume or not.

we can download park of file using:
 request.Header.Set("Range","bytes=start-end")

this will download file's size[start,end] park.
*/

func main() {
	// url := "https://dl.google.com/go/go1.11.1.src.tar.gz"
	err := CMD()
	if err != nil {
		fmt.Println(err)
	}
}
