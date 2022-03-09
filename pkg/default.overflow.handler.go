package pkg

import "fmt"

// defaultOverflowHandler used when no other [Active] option is found, the default handler can be over-written
// but there can only ever be one
func defaultOverflowHandler(contract *Contract) {
	temp, _ := contract.CurrentStep()
	contract.Err = fmt.Errorf("contract %s hit the default handler of %s service. If this is expected please over ride the default handler",
		contract.id, temp.String())
}
