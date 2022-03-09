package main

import (
	"fmt"
	"github.com/baron159/hermes"
	"github.com/baron159/hermes/example/impl"
)

func main() {
	s1 := impl.CreateService("S1")
	s2 := impl.CreateService("S2")
	s3 := impl.CreateService("S3")
	s4 := impl.CreateService("S4")
	cnt := hermes.NewContract(*s1.ID(), "bar", "1st-contract-start")
	cnt.AppendNext(*s3.ID(), "foo")
	cnt.AppendNext(*s4.ID(), "ping")
	cnt.AppendNext(*s2.ID(), "ping")
	cnt.AppendNext(*s1.ID(), "foo")
	println("contract ready for sending")
	mb, err := cnt.Send()
	if err != nil {
		panic(err)
	}
	fmt.Printf("The Contract was filled!\n%s", (<-mb).String())
}
