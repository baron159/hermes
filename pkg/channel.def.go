package pkg

type NotifyContractor chan<- Intent
type ContractorListener <-chan Intent

type IntendSend chan<- Intent
type ResultIntent <-chan Intent

func CreateContractorPair() (NotifyContractor, ContractorListener) {
	ch := make(chan Intent)
	return ch, ch // Looks silly but helps with clarity
}

func CreateIntentPair() (IntendSend, ResultIntent) {
	ch := make(chan Intent)
	return ch, ch
}
