package pkg

type ServiceDef interface {
	String() string // Should provide read-able insight about the service and current state
	ID() *ServiceID
	GetHandlers() ActionMap
	GetActions() ActionList /* Get a list of all the Action a service can handle */
}

type ActionMap map[ActionID]Action /* Map type  */
type ActionList []ActionID         /* List of ActionID */

type ServiceID struct {
	Id    string
	_noty NotifyContractor
	//maybe the addition of uuid field for absolute uniqueness
}

func (s *ServiceID) String() string {
	return s.Id
}
func (s *ServiceID) setNotify(ch NotifyContractor) {
	s._noty = ch
}

// GetNotify get the channel that is used to notify the Contractor that is running the ServiceDef
func (s *ServiceID) GetNotify() NotifyContractor {
	return s._noty
}

func (am ActionMap) Add(a Action) {
	am[a.GetId()] = a
}
