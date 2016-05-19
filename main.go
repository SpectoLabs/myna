package main

import (
    "fmt"
    "os"
    "os/exec"
    "io/ioutil"
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
func openBoltDb() (*bolt.DB, error) {
    location := os.Getenv("DATABASE_LOCATION")
    if location == "" {
        location = "processes.db"
    }
    return bolt.Open(location, 0600, nil)
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

func (p *Process) FromProcessJson(proc *ProcessJson) {
    p.Command = proc.Command
    p.Stdout = decodeBase64(proc.Stdout)
    p.Stderr = decodeBase64(proc.Stderr)
    p.ReturnCode = proc.ReturnCode
}

func (p *Process) FromJson(payload []byte) {
    result := ProcessJson{}
    json.Unmarshal(payload, &result)
    p.FromProcessJson(&result)
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
    db, err := openBoltDb();
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
    cmd := exec.Command(p.Command[0], p.Command[1:]...)
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
    p.Playback()
}

func (p *Process) Lookup() error {
    db, err := openBoltDb();
    if err != nil {
        return err
    }
    defer db.Close()
    err = db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("processes"))
        if b == nil {
            return fmt.Errorf("Nothing has been recorded yet. Try to capture some commands first")
        }
        v := b.Get(p.Key())
        if v == nil {
            return fmt.Errorf("%s does not know about this process. Run the command in capture mode first.", os.Args[0])
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

func Export() error {
    db, err := openBoltDb();
    if err != nil {
        return err
    }
    defer db.Close()
    payload := []ProcessJson{}
    db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("processes"))
        if b == nil {
            return nil
        }
        b.ForEach(func(k, v []byte) error {
            p := ProcessJson{}
            json.Unmarshal(v, &p)
            payload = append(payload, p)
            return nil
        })
        return nil
    })
    data, _ := json.Marshal(payload)
    os.Stdout.Write(data)
    return nil
}

func Import(file string) {
    payload, err := ioutil.ReadFile(file)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    parsed := []ProcessJson{}
    err = json.Unmarshal(payload, &parsed)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    for _, j := range parsed {
        p := Process{}
        p.FromProcessJson(&j)
        p.Save()
    }
}

func Usage() {
    fmt.Println(os.Args[0] + " [OPTS] [COMMAND] [[ARGS]]")
    fmt.Println("")
    fmt.Println(os.Args[0] + " can either capture or playback commands:")
    fmt.Println("")
    fmt.Println("Example:")
    fmt.Println("")
    fmt.Println("   $ " + os.Args[0] + " --capture ls -al /")
    fmt.Println("   $ " + os.Args[0] + " ls -al /")
    fmt.Println("")
    fmt.Println("Options:")
    fmt.Println("")
    fmt.Println("  --export                          Export database to json")
    fmt.Println("  --import [PATH]                   Import json to database")
    fmt.Println("  --capture [COMMAND] [[ARGS]]      Capture the output of running COMMAND [ARGS]")
}

func InCaptureMode() bool {
    return os.Getenv("CAPTURE") == "1"
}

func main() {
    if len(os.Args) == 1 {
        Usage()
    } else if os.Args[1] == "--export" {
        Export()
    } else if os.Args[1] == "--import" {
        Import(os.Args[2])
    } else if os.Args[1] == "--capture" {
        p := Process{}
        p.Command = os.Args[2:]
        p.Capture()
    } else {
        p := Process{}
        p.Command = os.Args[1:]
        if InCaptureMode() {
            p.Capture()
        } else if err := p.Lookup() ; err == nil {
            p.Playback()
        } else {
            fmt.Println(err.Error())
        }
    }
}
