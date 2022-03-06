package hermes

type overseer struct {
	startStamp int64
	serviceMap map[Service]*contractor
}

type Service interface {
	String() string
	GetNotifyChannel()
}

type contractor struct {
}