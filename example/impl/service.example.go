package impl

import (
	"fmt"
	"github.com/baron159/hermes"
	"github.com/baron159/hermes/pkg"
)

type SimpleService struct {
	id          *pkg.ServiceID
	actionGroup pkg.ActionMap
}

func CreateService(name string) *SimpleService {
	rtnService := &SimpleService{actionGroup: make(pkg.ActionMap)}
	if i, err := hermes.CreateServiceID(name); err != nil {
		panic(err)
	} else {
		rtnService.id = &i
	}
	if err := hermes.RegisterService(rtnService); err != nil {
		panic(err)
	}
	return rtnService
}

func (s SimpleService) String() string {
	return fmt.Sprintf(`
	ID: %s - ActGrpCnt: %d
	%v
`, s.id.Id, len(s.actionGroup), s.actionGroup)
}

func (s *SimpleService) ID() *pkg.ServiceID {
	return s.id
}

func (s *SimpleService) GetHandlers() pkg.ActionMap {
	var rtnMp = make(pkg.ActionMap)
	fn := func(n string) {
		a := StandardAction(n)
		rtnMp[a.GetId()] = a
	}
	fn("foo")
	fn("bar")
	fn("ping")
	fn("dabb")
	return rtnMp
}

func (s SimpleService) GetAction() pkg.ActionList {
	return pkg.ActionList{
		"foo", "bar", "ping", "dabb",
	}
}
