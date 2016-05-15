package main

import (
    "fmt"
    "os"
    "os/exec"
    "bytes"
    "syscall"
    "strings"
    "encoding/json"
    "github.com/boltdb/bolt"
    "encoding/base64"
)

type CommandResult struct {
    Command []string
    Stdout []byte
    Stderr []byte
    ReturnCode int
}

type CommandResultJson struct {
    Command []string
    Stdout string
    Stderr string
    ReturnCode int
}

func encodeBase64(str []byte) string {
    return base64.StdEncoding.EncodeToString(str)
}

func (cmd *CommandResult) Json() []byte {
    payload := CommandResultJson{}
    payload.Command = cmd.Command
    payload.ReturnCode = cmd.ReturnCode
    payload.Stdout = encodeBase64(cmd.Stdout)
    payload.Stderr = encodeBase64(cmd.Stderr)
    result, _ := json.Marshal(payload)
    fmt.Println((string)(result))
    return result
}

func (cmd *CommandResult) Key() []byte {
    return ([]byte)(strings.Join(cmd.Command, " "))
}

func (cmd *CommandResult) Print() {
    fmt.Println("Command: " + strings.Join(cmd.Command, " "))
    fmt.Println("Stdout: " + string(cmd.Stdout))
    fmt.Println("Stderr: " + string(cmd.Stderr))
    fmt.Printf("Return code: %d\n", cmd.ReturnCode)
}

func Save(cmd *CommandResult) error {
    db, err := bolt.Open("processes.db", 0600, nil)
    if err != nil {
        return err
    }
    defer db.Close()
    err = db.Update(func(tx *bolt.Tx) error {
        b, err := tx.CreateBucketIfNotExists([]byte("processes"))
        if err != nil {
            return fmt.Errorf("create bucket: %s", err)
        }
        return b.Put(cmd.Key(), cmd.Json())
    })
    if err != nil {
        return err
    }
    return nil
}

func main() {
    cmd := exec.Command(os.Args[1], os.Args[2:]...)
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    result := CommandResult{}
    result.Command = os.Args[1:]
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            waitStatus := exitError.Sys().(syscall.WaitStatus)
            result.ReturnCode = waitStatus.ExitStatus()
        } else {
            fmt.Println(err.Error())
            return
        }
    }
    result.Stdout = stdout.Bytes()
    result.Stderr = stderr.Bytes()
    result.Print()
    err = Save(&result)
    if err != nil {
        fmt.Println(err.Error())
    }

}
