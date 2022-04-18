package pkg

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

// Contract is the core data structure for the communication between services
type Contract struct {
	id         string
	startStamp int64 // Used for bench-marking, also used to know if its started?
	Params     []string
	Payload    interface{} // Dynamic Storage for the Contract
	Err        error
	completed  bool
	rtnChan    ContractResult // When all nexts are clear, send the Contract here
	nexts      []communicationStamp
	prevs      []communicationStamp
	// History? anytime anything is wiped. put the old data in a cache?
	// Asking ID/Service? this would allow for an additional check to make sure outside
	//		msg don't make it in
}

type communicationStamp struct {
	sendTo ServiceID
	doThis ActionID
}

func CreateContract(frtSer ServiceID, act ActionID, args ...string) *Contract {
	return &Contract{
		id:     uuid.New().String(),
		Params: args,
		nexts:  []communicationStamp{{frtSer, act}},
	}
}

// Send is the only thing func call required to start a Contract off. **AFTER ALL META & SERVICE DATA LOADED**
/*
	Mailbox is a chan type to listen for the response from the service,
	and error will be thrown if the Contract does not have any next's, effectively meaning there
	is no place for the Contract to go
*/
func (c *Contract) Send() (Mailbox, error) {
	if len(c.nexts) == 0 {
		return nil, fmt.Errorf("nexts can not be zero to start a contract")
	}
	n, l := CreateContractPair()
	c.rtnChan = n
	c.next()
	return l, nil
}

func (c *Contract) AppendNext(toService ServiceID, reqAction ActionID) {
	c.nexts = append(c.nexts, communicationStamp{
		sendTo: toService,
		doThis: reqAction,
	})
}

// AttachPayload takes in anything, and will over-write the exiting payload if any
func (c *Contract) AttachPayload(i interface{}) {
	c.Payload = i
}

// SafelyAttachPayload takes in anything, and will attach a payload only if nothing else is present
func (c *Contract) SafelyAttachPayload(i interface{}) error {
	if c.Payload != nil {
		return fmt.Errorf("already a payload present")
	}
	c.Payload = i
	return nil
}

// SetParams Over-writes the existing params if any exist
func (c *Contract) SetParams(args ...string) {
	c.Params = args
}

// ParseParams will only work if the params are using one of the valid sep[:;=]
func (c *Contract) ParseParams() (rtn map[string]string, err error) {
	rtn = make(map[string]string)
	if len(c.Params) == 0 {
		return
	}
	for _, s := range []string{":", ";", "="} {
		if len(strings.Split(c.Params[0], s)) == 2 {
			for _, p := range c.Params {
				ps := strings.Split(p, s)
				if len(ps) != 2 {
					err = fmt.Errorf("param used that doesn't match split")
				}
				rtn[ps[0]] = ps[1]
			}
			break
		}
	}
	return
}

// AppendParams leaves the existing params
func (c *Contract) AppendParams(args ...string) {
	c.Params = append(c.Params, args...)
}

// next How to send contracts on their way
func (c *Contract) next() {
	if c.startStamp == 0 {
		c.startStamp = time.Now().UnixMilli()
	}
	if c.hasNext() {
		step, _ := c.popNext()
		step.sendTo.GetNotify() <- *c
	} else {
		if c.rtnChan == nil {
			panic("No channel to return the contract for")
		}
		if c.Err == nil {
			c.completed = true
		}
		c.rtnChan <- *c
	}
}
func (c *Contract) popNext() (nxt communicationStamp, err error) {
	if !c.hasNext() {
		err = fmt.Errorf("no next to grab")
	} else {
		nxt = c.nexts[0]
		if len(c.nexts) > 1 {
			c.nexts = c.nexts[1:]
		} else {
			c.nexts = nil
		}
		c.prevs = append(c.prevs, nxt)
	}
	return
}
func (c Contract) hasNext() bool {
	return c.nexts != nil && len(c.nexts) > 0 && c.Err == nil
}
func (c Contract) GetID() string {
	return c.id
}
func (c Contract) IsCompleted() bool {
	return c.completed && c.Err == nil
}
func (c Contract) CurrentStep() (ServiceID, ActionID) {
	if c.prevs == nil || len(c.prevs) == 0 {
		return ServiceID{}, ""
	}
	el := c.prevs[len(c.prevs)-1]
	return el.sendTo, el.doThis
}
func (c Contract) String() string {
	cSer, cAct := c.CurrentStep()
	return fmt.Sprintf(`
	Current Step | Service: %s RequestingAction: %d
	Contract ID: %s  Timestamp:%v
	Params:%s
	Payload: %v
	Error?: %s
	completed?: %t
	nexts: %v
	pasts: %v
`, cSer.String(), cAct, c.id, time.UnixMilli(c.startStamp), c.Params,
		c.Payload, c.Err, c.completed, c.nexts, c.prevs,
	)
}
func (c Contract) StartTime() int64 { return c.startStamp }
