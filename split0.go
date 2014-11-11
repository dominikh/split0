package main // import "honnef.co/go/split0"

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

func help() {
	fmt.Fprintf(os.Stderr, `Usage: %s <bytes> <command>

split0 splits stdin into chunks of N bytes, passing each chunk to a
new invocation of a command.

The command will be invoked as /bin/sh -c <command>, and an
environment variable $N will be set to the current iteration number.

Example:
$ head -c 12 /dev/zero | split0 5 'echo -n "$N: "; wc -c'
0: 5
1: 5
2: 2`+"\n", os.Args[0])
	os.Exit(0)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <bytes> <command>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		help()
	}
	if len(os.Args) != 3 {
		usage()
	}
	num, err := strconv.Atoi(os.Args[1])
	if err != nil || num < 1 {
		usage()
	}
	command := os.Args[2]

	var buf [1]byte
	br := bytes.NewReader(buf[:])
	i := 0
	for {
		// Check if there's any data left
		n, rerr := os.Stdin.Read(buf[:])
		if n == 0 || rerr != nil {
			break
		}

		br.Seek(0, os.SEEK_SET)
		os.Setenv("N", strconv.Itoa(i))
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdin = io.MultiReader(br, &io.LimitedReader{R: os.Stdin, N: int64(num - 1)})
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		i++
	}
}
