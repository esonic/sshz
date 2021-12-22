package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	args := []string{"-tt"}
	args = append(args, os.Args[1:]...)
	cmd := exec.Command("ssh", args...)
	cmd.Stderr = os.Stderr

	sh := &shell{}
	var err error
	sh.readC, err = cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	sh.writeC, err = cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	if err := setRawMode(true); err != nil {
		panic(err)
	}
	defer func() { _ = setRawMode(false) }()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := sh.Read(buf)
			if err != nil {
				return
			}
			data := buf[:n]
			if did, err := DoRzSz(data, sh); err != nil {
				fmt.Printf("rzsz error: %s\n", err)
				_ = setRawMode(false)
				os.Exit(1)
				return
			} else if did {
				// skip output
				continue
			}
			writeAll(os.Stdout, data)
		}
	}()

	go func() {
		_, _ = io.Copy(sh.writeC, os.Stdin)
	}()

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	_ = cmd.Wait()
}

func writeAll(out io.Writer, data []byte) {
	for len(data) > 0 {
		w, err := out.Write(data)
		if err != nil {
			fmt.Printf("out error: %s\n", err)
			_ = setRawMode(false)
			os.Exit(1)
			return
		}
		data = data[w:]
	}
}

func setRawMode(start bool) error {
	mode := "raw"
	echo := "-echo"
	if !start {
		mode = "-raw"
		echo = "echo"
	}

	rawMode := exec.Command("stty", mode, echo)
	rawMode.Stdin = os.Stdin
	if err := rawMode.Run(); err != nil {
		return err
	}

	return nil
}

type shell struct {
	readC  io.Reader
	writeC io.Writer
}

func (r *shell) Write(p []byte) (n int, err error) {
	return r.writeC.Write(p)
}

func (r *shell) Read(p []byte) (n int, err error) {
	return r.readC.Read(p)
}
