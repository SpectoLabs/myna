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

type Process struct {
    Command []string
    Stdout []byte
    Stderr []byte
    ReturnCode int
}

type ProcessJson struct {
    Command []string
    Stdout string
    Stderr string
    ReturnCode int
}

func encodeBase64(str []byte) string {
    return base64.StdEncoding.EncodeToString(str)
}
func decodeBase64(str string) []byte {
    data, _ := base64.StdEncoding.DecodeString(str)
    return data
}

func (p *Process) Json() []byte {
    payload := ProcessJson{}
    payload.Command = p.Command
    payload.ReturnCode = p.ReturnCode
    payload.Stdout = encodeBase64(p.Stdout)
    payload.Stderr = encodeBase64(p.Stderr)
    result, _ := json.Marshal(payload)
    return result
}

func (p *Process) FromJson(payload []byte) {
    result := ProcessJson{}
    json.Unmarshal(payload, &result)
    p.Stdout = decodeBase64(result.Stdout)
    p.Stderr = decodeBase64(result.Stderr)
    p.ReturnCode = result.ReturnCode
}

func (p *Process) Key() []byte {
    return ([]byte)(strings.Join(p.Command, " "))
}

func (p *Process) Print() {
    fmt.Println("Command: " + strings.Join(p.Command, " "))
    fmt.Println("Stdout: " + string(p.Stdout))
    fmt.Println("Stderr: " + string(p.Stderr))
    fmt.Printf("Return code: %d\n", p.ReturnCode)
}

func (p *Process) Save() error {
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
        return b.Put(p.Key(), p.Json())
    })
    return err
}

func (p *Process) Capture() {
    cmd := exec.Command(os.Args[1], os.Args[2:]...)
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    if err := cmd.Run() ; err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            waitStatus := exitError.Sys().(syscall.WaitStatus)
            p.ReturnCode = waitStatus.ExitStatus()
        } else {
            fmt.Println(err.Error())
            return
        }
    }
    p.Stdout = stdout.Bytes()
    p.Stderr = stderr.Bytes()
    err := p.Save()
    if err != nil {
        fmt.Println(err.Error())
    }
}

func (p *Process) Lookup() error {
    db, err := bolt.Open("processes.db", 0600, nil)
    if err != nil {
        return err
    }
    defer db.Close()
    err = db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("processes"))
        v := b.Get(p.Key())
        if v == nil {
            return fmt.Errorf("%s does not know this process. Run the command in capture mode first.", os.Args[0])
        }
        p.FromJson(v)
        return nil
    })
    return err
}

func (p *Process) Playback() error {
    os.Stdout.Write(p.Stdout)
    os.Stderr.Write(p.Stderr)
    os.Exit(p.ReturnCode)
    return nil
}

func main() {
    p := Process{}
    p.Command = os.Args[1:]
    p.Capture()
    if err := p.Lookup() ; err == nil {
        p.Playback()
    } else {
        fmt.Println(err.Error())
    }
}
