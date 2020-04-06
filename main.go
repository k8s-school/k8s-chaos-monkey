package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"sync"
	"time"
)

const ShellToUse = "bash"
const RdrPrefix = "qserv-dev-xrootd-redirector"
const RdrReplicas = 4

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func killProc(i int, messages chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	command := fmt.Sprintf("kubectl exec -i %s-%d -c xrootd -- kill 12", RdrPrefix, i)
	log.Printf("Launching %s\n", command)
	err, out, errout := Shellout(command)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	log.Println("--- stdout ---")
	log.Println(out)
	log.Println("--- stderr ---")
	log.Println(errout)
	command = fmt.Sprintf("kubectl exec -i %s-%d -c cmsd -- kill 11", RdrPrefix, i)
	log.Printf("Launching %s\n", command)
	err, out, errout = Shellout(command)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	log.Println("--- stdout ---")
	log.Println(out)
	log.Println("--- stderr ---")
	log.Println(errout)
	messages <- i

}

func main() {

	for {

		alive := rand.Intn(RdrReplicas)
		alive2 := rand.Intn(RdrReplicas)
		messages := make(chan int)
		var wg sync.WaitGroup
		wg.Add(RdrReplicas - 1)
		var i int

		for i = 0; i < RdrReplicas; i++ {
			if i != alive && i != alive2 {
				go killProc(i, messages, &wg)
			}
		}
		log.Printf("Sleep 15s")
		time.Sleep(25 * time.Second)

	}

}
