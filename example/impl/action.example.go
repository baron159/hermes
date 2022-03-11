package impl

import (
	"fmt"
	"github.com/baron159/hermes/pkg"
	"math/rand"
	"time"
)

type StandardAction string

func (s StandardAction) String() string {
	return string(s)
}

func (s StandardAction) GetId() pkg.ActionID {
	return pkg.ActionID(s)
}

func (s StandardAction) GetHandler() pkg.ContractActionHandler {
	return func(c *pkg.Contract) {
		si, _ := c.CurrentStep()
		c.AppendParams(fmt.Sprintf("%s:%s:passed", si.Id, s.GetId()))
		dur := rand.Intn(10) * int(time.Second)
		time.Sleep(time.Duration(dur)) // Simulate some load
	}
}
