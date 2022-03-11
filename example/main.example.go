package main

import (
	"fmt"
	"github.com/baron159/hermes"
	"github.com/baron159/hermes/example/impl"
	"github.com/baron159/hermes/pkg"
	"math/rand"
	"time"
)

func main() {
	s1 := impl.CreateService("S1")
	s2 := impl.CreateService("S2")
	s3 := impl.CreateService("S3")
	s4 := impl.CreateService("S4")
	//cnt := hermes.NewContract(*s1.ID(), "bar", "1st-contract-start")
	//cnt.AppendNext(*s3.ID(), "foo")
	//cnt.AppendNext(*s4.ID(), "ping")
	//cnt.AppendNext(*s2.ID(), "ping")
	//cnt.AppendNext(*s1.ID(), "foo")
	serLt := []pkg.ServiceID{*s1.ID(), *s2.ID(), *s3.ID(), *s4.ID()}
	actLt := pkg.ActionList{"foo", "ping", "bar", "dabb"}
	rndCommFn := func() (pkg.ServiceID, pkg.ActionID) {
		return serLt[rand.Intn(len(serLt))], actLt[rand.Intn(len(actLt))]
	}
	const clients = 20
	const clientCycles = 4000
	completed := make(chan bool, clients)
	for client := 0; client < clients; client++ {
		go func(clt int, done chan bool) {
			cycle := 0
			for {
				tasks := rand.Intn(10)
				initSer, initAct := rndCommFn()
				newCont := hermes.NewContract(initSer, initAct, fmt.Sprintf("c%d-Starting:%d", clt, cycle))
				for t := 0; t < tasks; t++ {
					newCont.AppendNext(rndCommFn())
				}
				if mb, err := newCont.Send(); err != nil {
					panic(err)
				} else {
					//go mailboxWatch(mb)
					mailboxWatch(mb)
				}
				dur := int(time.Millisecond) * rand.Intn(6)
				time.Sleep(time.Duration(dur))
				cycle += 1
				if cycle > clientCycles {
					completed <- true
					break
				}
			}
			fmt.Printf("Client-%d Finished Cycles!\n", clt)
		}(client, completed)
	}
	for cap(completed) != len(completed) {
	}
	//select {}
	println("Completed")
}

func mailboxWatch(m pkg.Mailbox) {
	select {
	case r := <-m:
		dur := time.Since(time.UnixMilli(r.StartTime()))
		fmt.Printf("Response Mail - Took %v\n%s", dur, r.String())
		//case <-time.After(time.Second * 8):
		//	fmt.Printf("8 Sec timeout was hit, reponse has not yet come")
	}
}
