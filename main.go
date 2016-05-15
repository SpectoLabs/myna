package main

import (
    "fmt"
    "os"
    "os/exec"
    "bytes"
    "syscall"
    "strings"
)

type CommandResult struct {
    Command []string
    Stdout []byte
    Stderr []byte
    ReturnCode int
}

func (cmd *CommandResult) Print() {
    fmt.Println("Command: " + strings.Join(cmd.Command, " "))
    fmt.Println("Stdout: " + string(cmd.Stdout))
    fmt.Println("Stderr: " + string(cmd.Stderr))
    fmt.Printf("Return code: %d\n", cmd.ReturnCode)
}

func main() {
    cmd := exec.Command(os.Args[1], os.Args[2:]...)
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    exitCode := 0
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            waitStatus := exitError.Sys().(syscall.WaitStatus)
            exitCode = waitStatus.ExitStatus()
        }
    }
    result := CommandResult{}
    result.Command = os.Args[1:]
    result.Stdout = stdout.Bytes()
    result.Stderr = stderr.Bytes()
    result.ReturnCode = exitCode
    result.Print()
}
