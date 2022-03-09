package pkg

type Action interface {
	String() string // Should be useful insight about the state
	GetId() ActionID
	GetHandler() ContractActionHandler
}

type ActionID string
type ContractActionHandler func(c *Contract)
