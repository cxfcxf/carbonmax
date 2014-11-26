package main

import (
        "fmt"
        "time"
        "strings"
        "strconv"
        "net"
        "log"
        "os/exec"
        "sync"
        "flag"
        "runtime"
        "github.com/vaughan0/go-ini"
)

var (
        inifile = flag.String("inifile","/etc/carbonmax.ini", "path to your ini config file")
        loop = flag.Bool("loop", false, "switch on if you want to loop the program for daemonization")
)

func iniParser(inifile string) (map[string]string, map[string]string){
        file, err := ini.LoadFile(inifile)
        if err != nil {
                log.Fatal(err)
        }
        return file["carbonlink"], file["resources"]
}

func feedcarbon(status map[string]string, carbonlink map[string]string) {

        conn, err := net.Dial("tcp", carbonlink["server"] + ":" + carbonlink["port"])
        if err != nil {
                log.Println("Can not connect the carbon-cache, Please check setting")
        }
        defer conn.Close()

        var message string
        for mn, ms := range status {
                message = message + fmt.Sprintf("%s.%s %s %d\n", carbonlink["client"], mn, ms, time.Now().Unix())
        }

        conn.Write([]byte(message))
        verbose, _ := strconv.ParseBool(carbonlink["verbose"])
        if verbose {
                fmt.Println(message)
        }
}

func cmdExec(name string, command string, timeout string, status map[string]string, wg *sync.WaitGroup) {
        defer wg.Done()
        ch := make(chan string, 1)
        to, _ := time.ParseDuration(timeout)
        go func() {
                result, _ := exec.Command("sh", "-c", command).Output()
                ch <- string(result)
        }()

        select {
        case result := <-ch:
                if result != "" {
                        res := strings.Trim(result,"\n")
                        status[name] = res
                }
        case <-time.After(to):
                log.Printf("%s execution timed out", name)
        }
}

func main() {
        flag.Parse()

        var wg sync.WaitGroup

        if *loop {
                log.Println("Now Looping Carbonmax at Given Interval")

                for {
                        carbonlink, resources := iniParser(*inifile)
                        cpus, err := strconv.Atoi(carbonlink["cpus"])
                        if err != nil {
                                runtime.GOMAXPROCS(1)
                        } else {
                                runtime.GOMAXPROCS(cpus)
                        }

                        status := make(map[string]string)

                        for k, v := range resources {
                                if k != "" && v != "" {
                                        wg.Add(1)
                                        go cmdExec(k, strings.Trim(v, "`"), carbonlink["exectimeout"], status, &wg)
                                }
                        }
                        wg.Wait()

                        feedcarbon(status, carbonlink)
                        interval, _ := time.ParseDuration(carbonlink["interval"])

                        time.Sleep(interval)
                }
        } else {
                carbonlink, resources := iniParser(*inifile)
                cpus, err := strconv.Atoi(carbonlink["cpus"])
                if err != nil {
                        runtime.GOMAXPROCS(1)
                } else {
                        runtime.GOMAXPROCS(cpus)
                }

                status := make(map[string]string)

                for k, v := range resources {
                        if k != "" && v != "" {
                                        wg.Add(1)
                                go cmdExec(k, strings.Trim(v, "`"), carbonlink["exectimeout"], status, &wg)
                        }
                }
                wg.Wait()

                feedcarbon(status, carbonlink)
        }
}
