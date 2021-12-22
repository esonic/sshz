package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func DoRzSz(message []byte, sh io.ReadWriter) (bool, error) {
	msg := fmt.Sprintf("%q", message)
	switch {
	case strings.Contains(msg, "**\\x18B00000000000"):
		return true, doSz(sh)
	case strings.Contains(msg, "rz waiting to receive.**\\x18B0100"):
		return true, doRz(sh)
	}
	return false, nil
}

func writeCancelSz(r io.Writer) {
	_, _ = r.Write([]byte{'\x18', '\x18', '\x18', '\x18', '\x18'})
}

func doSz(r io.ReadWriter) error {
	dir, err := selectDir()
	// log.Printf("dir=%q, err=%v", dir, err)
	if err != nil {
		return err
	} else if dir == "" {
		writeCancelSz(r)
		// _, _ = r.Write([]byte("# Cancel Choose Dir\r"))
		return nil
	}

	cmd := exec.Command("rz", "-b", "-e", "-y")
	cmd.Dir = dir
	cmd.Stdin = r
	cmd.Stdout = r
	cmd.Stderr = os.Stderr
	go func() {
		for {
			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				_, _ = r.Write([]byte{'\r'})
				return
			}
			time.Sleep(time.Second / 10)
		}
	}()
	if err = cmd.Run(); err != nil {
		return err
	}
	// _, _ = r.Write([]byte(fmt.Sprintf("# Sent -> %q\r", dir)))

	return nil
}

func doRz(r io.ReadWriter) error {
	file, err := selectFile()
	// log.Printf("file=%q, err=%v", file, err)
	if err != nil {
		return err
	} else if file == "" {
		writeCancelSz(r)
		// _, _ = r.Write([]byte("# Cancel Choose File\r"))
		return nil
	}

	cmd := exec.Command("sz", file, "-e", "-b")
	cmd.Stdin = r
	cmd.Stdout = r
	cmd.Stderr = os.Stderr
	go func() {
		for {
			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				_, _ = r.Write([]byte{'\r'})
				return
			}
			time.Sleep(time.Second / 10)
		}
	}()
	if err = cmd.Run(); err != nil {
		return err
	}
	// _, _ = r.Write([]byte(fmt.Sprintf("# Received %q\r", file)))

	return nil
}

func selectDir() (string, error) {
	cmd := exec.Command("powershell.exe",
		`Add-Type -AssemblyName System.windows.forms;$f=New-Object System.Windows.Forms.`+
			`FolderBrowserDialog;if($f.ShowDialog()){wsl wslpath $f.SelectedPath.Replace("\", "\\")}`)
	res, err := cmd.Output()
	return strings.TrimSpace(string(res)), err
}

func selectFile() (string, error) {
	cmd := exec.Command("powershell.exe",
		`Add-Type -AssemblyName PresentationFramework;$f=New-Object Microsoft.Win32.`+
			`OpenFileDialog;if($f.ShowDialog()){wsl wslpath $f.FileName.Replace("\", "\\")}`)
	res, err := cmd.Output()
	return strings.TrimSpace(string(res)), err
}
