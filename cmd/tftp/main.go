package main

import (
	"flag"
	"github.com/lantosgyuri/networking/udp/tftp"
	"io/ioutil"
	"log"
)

var(
	address = flag.String("a", "127.0.0.1:69", "listen address")
	payload = flag.String("p", "../../asset/go-fuzz.svg", "file to serve the clients")
)

// Run with root permission, and first build the binary

func main() {
	flag.Parse()

	p, err := ioutil.ReadFile(*payload)
	if err != nil {
		log.Fatalln(err)
	}

	s := tftp.Server{Payload: p}

	log.Fatal(s.ListenAdnServe(*address))
}
