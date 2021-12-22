# SSHZ

A command line tool wraps ssh client and adds rz/sz support for windows terminal in WSL.

## Prerequisite

> Tested with Windows 11 and WSL2. Should work in Windows 10 and WSL1 as well.

1. Make sure rz/sz installed in WSL, eg. `sudo apt install lrzsz` in Debian.
2. Make sure Golang installed in WSL. [https://go.dev/doc/install](https://go.dev/doc/install)

## Installation

Run the following command in WSL to install the tool:

`go install github.com/esonic/sshz@latest`

The tool will be installed in $GOPATH/bin.

## Usage

The tool can be used as same as `ssh`. eg. `sshz someone@127.0.0.1`

The file dialog will appear when call rz and sz.

![image](https://raw.githubusercontent.com/esonic/scp2remote/master/screenshot-20211222-175036.jpg)
