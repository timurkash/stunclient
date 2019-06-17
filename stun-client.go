package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gortc/stun"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, os.Args[0], "stun.l.google.com:19302")
		fmt.Fprintln(os.Stderr, os.Args[1], "timur:realm:pass")
	}
	flag.Parse()
	addr := flag.Arg(0)
	if addr == "" {
		//addr = "stun.l.google.com:19302"
		addr = "localhost:3479"
	}
	lta := flag.Arg(1)
	if lta == "" {
		lta = "timur:realm:pass"
	}
	xorAddr, software, err := doDial(addr, lta)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(xorAddr)
	fmt.Println(software)
}

func doDial(addr, lta string) (string, string, error) {
	c, err := stun.Dial("udp", addr)
	if err != nil {
		return "", "", err
	}
	var (
		result string
		result2 string
	)
	userName := stun.NewUsername(lta)
	m := stun.MustBuild(
		stun.TransactionID,
		stun.BindingRequest,
		userName,
	)
	f := func(res stun.Event) {
		if res.Error != nil {
			log.Fatalln(err)
		}
		var software stun.Software
		if getErr := software.GetFrom(res.Message); getErr != nil {
			log.Fatalln(getErr)
		}
		result2 = software.String()
		//fmt.Println(software)
		var xorAddr stun.XORMappedAddress
		if getErr := xorAddr.GetFrom(res.Message); getErr != nil {
			log.Fatalln(getErr)
		}
		result = xorAddr.String()
	}
	if err != nil {
		return "", "", err
	}
	if err = c.Do(m, f); err != nil {
		fmt.Println(err)
		return "", "", err
	}
	if err := c.Close(); err != nil {
		return "", "", err
	}
	return result, result2, nil
}