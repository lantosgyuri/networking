package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	destinationServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		fmt.Printf("Error %v",err)
	}

	defer destinationServer.Close()

	go func() {
		conn ,err := destinationServer.Accept()
		if err != nil {
			fmt.Printf("Error %v",err)
			return
		}

		defer conn.Close()

		buffer := make([]byte, 1024)
		for {
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Printf("Error %v",err)
				return
			}

			_, err = conn.Write(buffer)
			if err != nil {
				fmt.Printf("Error %v",err)
				return
			}
		}

	}()

	sender, err := net.Dial("tcp", destinationServer.Addr().String())
	if err != nil {
		fmt.Printf("Error %v",err)
	}

	defer sender.Close()

	buffer := make([]byte, 1024)

	_, err = sender.Write([]byte("Ping and piond"))
	n, err := sender.Read(buffer)
	if err != nil {
		fmt.Printf("Error %v",err)
		return
	}
	if n > 0 {
		fmt.Printf("n is: %v \n", n)
		fmt.Printf("message is:  %v \n", string(buffer[:n]))
	}

	_, err = sender.Write([]byte("Pklahsjglkjsadhg!!"))
	_, err = sender.Write([]byte("!!!!!!!!!!!Psagmiésdgsdian"))

	time.Sleep(1 * time.Second)
		_, err = sender.Read(buffer)
		if err != nil {
			fmt.Printf("Error %v",err)
			return
		}
		if n > 0 {
			fmt.Printf("n is: %v \n", n)
			fmt.Printf("message is:  %v \n", string(buffer[:n]))
		}
}


