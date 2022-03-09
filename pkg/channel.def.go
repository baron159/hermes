package pkg

type NotifyContractor chan<- Contract
type ContractorListener <-chan Contract

type ContractResult chan<- Contract // Held inside the Contract, and will be used when all conditions are met
type Mailbox <-chan Contract        // used by consumers wait for there responses to Contract request

type OutsourceNotify chan<- Contract
type OutsourceConsumer <-chan Contract

// CreateContractorPair intended to be the long term listener into the Job-Loop
func CreateContractorPair() (NotifyContractor, ContractorListener) {
	ch := make(chan Contract)
	return ch, ch // Looks silly but helps with clarity
}

// CreateContractPair intended to be used as short-term query/result type relation
//	Mailbox is intended to be used by the provolker to get the results of the submitted
// 	Contract. ContractResult is intended to be nested inside the Contract it is tied too
func CreateContractPair() (ContractResult, Mailbox) {
	ch := make(chan Contract)
	return ch, ch
}
