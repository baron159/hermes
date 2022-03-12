package pkg

import (
	"flag"
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

var maxLoopBuff int
var maxWorkers int

const DefaultLoopBuffSize = 3

func init() {
	flag.IntVar(&maxLoopBuff, "mq", 5, "The max size of the queue for the Job Loop")
	flag.IntVar(&maxWorkers, "mx", 5, "The max number of workers that a Contractor can have, should not be bigger than the queue")
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

// Start  notify loop off, also kicking off the workers. Takes in optional int params that control the queue size and number of workers
/*
 Current args:
	[0] - Worker Count: the number of separate worker loops are notified on a first-available worker
	[1] - Queue Length: the size of the buffer for the worker. Note: this should never be less than WC
*/
func (cr *Contractor) Start(args ...int) error {
	cr.once.Do(func() {
		// How many workers
		wkr := 1
		if len(args) > 0 && int(math.Abs(float64(args[0]))) <= maxWorkers {
			wkr = int(math.Abs(float64(args[0])))
		}
		// Set up the channel queue
		queue := DefaultLoopBuffSize
		if len(args) > 1 && int(math.Abs(float64(args[1]))) <= maxLoopBuff {
			queue = int(math.Abs(float64(args[1])))
			if queue < wkr {
				println("The queue should not be smaller than the total worker pool")
				queue = wkr
			}
		}
		jLoop := make(chan Contract, queue)
		cr.jobQueue = jLoop
		cr.startTime = time.Now().UnixMilli()
		for w := 0; w < wkr; w++ {
			// Set-up Job Loop Listener that handles the fetching and calling of the handler functions
			go func(ct *Contractor, wk int) {
				// Hire the number of workers we need
				for c := range ct.jobQueue {
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

// CheckQueue returns the number of things in queue, and a bool on where the queue is filled
//  it should be noted that both these values can change from one instance to the next
func (cr *Contractor) CheckQueue() (qLen int, filled bool) {
	return len(cr.jobQueue), len(cr.jobQueue) == cap(cr.jobQueue)
}
