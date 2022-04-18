package hermes

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/baron159/hermes/pkg"
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
//  This must be used before RegisterService
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
/*
 	Current args:
	[0] - Worker Count: the number of separate worker loops are notified on a first-available worker
	[1] - Queue Length: the size of the buffer for the worker. Note: this should never be less than WC
*/
func RegisterService(ser pkg.ServiceDef, args ...int) error {
	for s, _ := range _Overseer().serviceMap {
		if ser.ID().String() == s.String() {
			return fmt.Errorf("service def w/ID: %s already in use", ser.ID())
		}
	}
	ctr := pkg.CreateContractor(ser)
	_Overseer().serviceMap[ser.ID()] = ctr
	if err := ctr.Start(args...); err != nil {
		return err
	}
	return nil
}

// ListServices returns a list of service IDs from currently active services
func ListServices() (rtnLt []pkg.ServiceID) {
	for s, _ := range _Overseer().serviceMap {
		rtnLt = append(rtnLt, *s)
	}
	return
}

// NewContract is the starting point for creating a new request from one or more services
func NewContract(frtSer pkg.ServiceID, act pkg.ActionID, args ...string) *pkg.Contract {
	// TODO make sure the service we are being passed is legit, validate action too?
	return pkg.CreateContract(frtSer, act, args...)
}

// FetchServiceID takes in a param and returns only if pkg.ServiceID can match the passed in string as in active Service Map
func FetchServiceID(n string) (pkg.ServiceID, error) {
	for s, _ := range _Overseer().serviceMap {
		if strings.EqualFold(n, s.Id) {
			return *s, nil
		}
	}
	return pkg.ServiceID{}, fmt.Errorf("no active service for: %s", n)
}

// LookUpServiceID takes in a param and returns only if pkg.ServiceID can match the passed in string as in active Service Map
//  Exact same method as FetchServiceID except uses the 'okay:bool'
func LookUpServiceID(n string) (pkg.ServiceID, bool) {
	for s, _ := range _Overseer().serviceMap {
		if strings.EqualFold(n, s.Id) {
			return *s, true
		}
	}
	return pkg.ServiceID{}, false
}
