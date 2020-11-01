package main

import (
	"fmt"
	"log"
)

type Client struct {
	name string
	Callback ClientTrace
}

func (c *Client) PrintName() error {
	done := c.Callback.onPrintName()
	err := c.doPrintName()
	done(err)
	return err
}

func (c *Client) doPrintName() error {
	fmt.Printf("Hello, %s", c.name)
	return nil
}

type ClientTrace struct {
	onPrintName func() func(err error)
}

func (a ClientTrace) Compose(b ClientTrace) (c ClientTrace) {
	switch {
	case a.onPrintName == nil:
		c.onPrintName = b.onPrintName
	case b.onPrintName == nil:
		c.onPrintName = a.onPrintName
	default:
		c.onPrintName = func() func(err error) {
			doneA := a.onPrintName()
			doneB := b.onPrintName()
			switch {
			case doneA == nil:
				return doneB
			case doneB == nil:
				return doneA
			default:
				return func(err error) {
					doneA(err)
					doneB(err)
				}
			}
		}
	}
	return c
}


func main() {

	var trace ClientTrace

	trace = trace.Compose(
		ClientTrace{
			onPrintName: func() func(err error) {
				log.Println("start printing name")
				return func(err error) {
					log.Println("printing done", err)
				}
			},
		},
	)

	trace = trace.Compose(ClientTrace{
		onPrintName: func() func(err error) {
			log.Println("start metrics")
			return func(err error) {
				log.Println("metrics done")
			}
		},
	})

	client := Client{
		name:     "foo",
		Callback: trace,
	}

	client.PrintName()
}
