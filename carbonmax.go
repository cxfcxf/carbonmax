package main

import (
    "os"
    "fmt"
    "net"
    "log"
    "time"
    "flag"
    "sync"
    "os/exec"
    "strings"
    "io/ioutil"
    "encoding/json"
)

var (
    f = flag.String("f","/etc/carbonmax/config.json", "path to your config file")
)

type Config struct {
    Carbonlink  *Carbonlink
    Metric      map[string]string
}

type Carbonlink struct {
    Server       string
    Client       string
    Exectimeout  time.Duration
    Interval     time.Duration
    Verbose      bool
    Daemonize    bool
}

func feedcarbon(status map[string]string, carbonlink *Carbonlink) {

    conn, err := net.Dial("tcp", carbonlink.Server)
    if err != nil {
        log.Println("Can not connect the carbon-cache, Please check setting")
    }
    defer conn.Close()

    var message string
    for mn, ms := range status {
        message = message + fmt.Sprintf("%s.%s %s %d\n", carbonlink.Client, mn, ms, time.Now().Unix())
    }

    conn.Write([]byte(message))

    if carbonlink.Verbose {
        log.Printf("\n%s", message)
    }
}

func cmdExec(name string, command string, timeout time.Duration, status map[string]string, wg *sync.WaitGroup) {

    defer wg.Done()
    ch := make(chan string, 1)
    go func() {
        result, _ := exec.Command("sh", "-c", command).Output()
        ch <- string(result)
    }()

    select {
    case result := <-ch:
        if result != "" {
            nsl := strings.Split(name, "|")
            metrics := strings.Split(strings.Trim(result,"\n"), "|")
            if len(nsl) == len(metrics) {
                for i, n := range nsl {
                    status[n] = metrics[i]
                }
            } else {
                log.Panic("number between metrics and name not match")
            }

        }
    case <-time.After(timeout * time.Second):
        log.Printf("%s execution timed out", name)
    }
}

func loadConf(config string) Config{
    f, err := ioutil.ReadFile(config)
    if err != nil {
        log.Panic(err)
    }

    var conf Config
    err = json.Unmarshal(f, &conf)
    if err != nil {
        log.Panic(err)
    }
    return conf
}

func main() {
    flag.Parse()

    var wg sync.WaitGroup
    status := make(map[string]string)

    for {
        conf := loadConf(*f)
        for k, v := range conf.Metric {
            if k != "" && v != "" {
                wg.Add(1)
                go cmdExec(k, v, conf.Carbonlink.Exectimeout, status, &wg)
            }
        }
        wg.Wait()

        feedcarbon(status, conf.Carbonlink)

        if !conf.Carbonlink.Daemonize {
            os.Exit(0)
        }

        time.Sleep(conf.Carbonlink.Interval * time.Second)
    }
}