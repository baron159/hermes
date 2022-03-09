package hermes

import (
	"flag"
	"fmt"
	"github.com/baron159/hermes/pkg"
	"sync"
	"time"
)

type overseer struct {
	startStamp int64
	serviceMap map[*pkg.ServiceID]*pkg.Contractor
}

var _instance *overseer
var _once sync.Once

func _Overseer() *overseer {
	_once.Do(func() {
		flag.Parse()
		_instance = &overseer{
			startStamp: time.Now().UnixMilli(),
			serviceMap: make(map[*pkg.ServiceID]*pkg.Contractor),
		}
	})
	return _instance
}

// CreateServiceID is used to wrap your service names in a structure for internal use
func CreateServiceID(id string) (pkg.ServiceID, error) {
	for i, _ := range _Overseer().serviceMap {
		if i.Id == id {
			return pkg.ServiceID{}, fmt.Errorf("service with that id already present: %s", id)
		}
	}
	return pkg.ServiceID{
		Id: id,
	}, nil
}

// RegisterService takes in a pkg.ServiceDef and makes sure that one does not already exist. It then
// creates a pkg.Contractor for the Service Definition, adds it to the _Overseer Map and Starts the Contractor
func RegisterService(ser pkg.ServiceDef) error {
	for s, _ := range _Overseer().serviceMap {
		if ser.ID().String() == s.String() {
			return fmt.Errorf("service def w/ID: %s already in use", ser.ID())
		}
	}
	ctr := pkg.CreateContractor(ser)
	_Overseer().serviceMap[ser.ID()] = ctr
	if err := ctr.Start(); err != nil {
		return err
	}
	return nil
}

func ListServices() (rtnLt []pkg.ServiceID) {
	for s, _ := range _Overseer().serviceMap {
		rtnLt = append(rtnLt, *s)
	}
	return
}

func NewContract(frtSer pkg.ServiceID, act pkg.ActionID, args ...string) *pkg.Contract {
	// TODO make sure the service we are being passed is legit, validate action too?
	return pkg.CreateContract(frtSer, act, args...)
}
