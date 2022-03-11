package pkg

import (
	"flag"
	"fmt"
	"math"
	"sync"
	"time"
)

// Contractor requires a ServiceDef to create
type Contractor struct {
	startTime         int64
	purpose           ServiceDef
	overflowFn        ContractActionHandler
	outsourcePipeline []OutsourceNotify
	jobQueue          chan Contract
	listener          ContractorListener
	once              sync.Once
}

var jobLoopBuff int
var maxWorkers int

func init() {
	flag.IntVar(&jobLoopBuff, "lp", 3, "The max size of the queue for the Job Loop")
	flag.IntVar(&maxWorkers, "mx", 10, "The max number of workers that a Contractor can have, should not be bigger than the queue")
}

func CreateContractor(pur ServiceDef) *Contractor {
	n, l := CreateContractorPair()
	pur.ID().setNotify(n)
	return &Contractor{
		purpose:    pur,
		overflowFn: defaultOverflowHandler,
		listener:   l,
	}
}

func (cr *Contractor) Start(args ...int) error {
	cr.once.Do(func() {
		jLoop := make(chan Contract, jobLoopBuff)
		cr.jobQueue = jLoop
		cr.startTime = time.Now().UnixMilli()
		// How many workers to hire
		wkr := 1
		if len(args) > 0 && int(math.Abs(float64(args[0]))) <= maxWorkers {
			wkr = int(math.Abs(float64(args[0])))
		}
		for w := 0; w < wkr; w++ {
			// Set-up Job Loop Listener that handles the fetching and calling of the handler functions
			go func(ct *Contractor, wk int) {
				// Hire the number of workers we need
				for c := range ct.jobQueue {
					fmt.Printf("Ctr: %s worker: %d handling %s\n", ct.purpose.ID().Id, wk, c.id)
					_, actId := c.CurrentStep()
					if act, ok := ct.purpose.GetHandlers()[actId]; !ok {
						panic("contract given to job loop without know handler")
					} else {
						act.GetHandler()(&c)
						c.next()
					}
				}
			}(cr, w)
		}

		// Set-up Contractor listener listener
		go func(ct *Contractor) {
			for c := range ct.listener {
				ser, actId := c.CurrentStep()
				if ser.String() != ct.purpose.ID().String() {
					panic("Given Contract for different service")
				}
				if _, ok := ct.purpose.GetHandlers()[actId]; ok { // Case 1 - We have an Action handler
					ct.jobQueue <- c
				} else if len(ct.outsourcePipeline) > 0 { // Case 2 - We have a listener
					// Do a broadcast to all listeners, which should handle the contract next() step call
					panic("not implimented")
				} else { // Case 3 - Use Overflow Handler [Last Resort]
					ct.overflowFn(&c)
					c.next()
				}
			}
		}(cr)
	})
	return nil
}

func (cr *Contractor) Term() error {
	// Closes all channels and sends shut-down signal now the pipe
	panic("fast term")
}

func (cr *Contractor) CheckQueue() (qLen int, filled bool) {
	return len(cr.jobQueue), len(cr.jobQueue) == cap(cr.jobQueue)
}
