package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var clientsList []int

func listener(c net.Conn) {
	for {
		message, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(message) > 4 {
			//message received
			fmt.Printf("Message received on client:%v\n", c)
			messageType := message[0:4]
			switch messageType {
			case "IDME":
				fmt.Println("\nMY ID->: " + message[4:])
			case "LIST":
				var newList []int
				strIds := string(message[4 : len(message)-1])
				clientIds := strings.Split(strIds, ",")
				for _, id := range clientIds {
					i, err := strconv.Atoi(id)
					if err != nil {
						fmt.Println(err)
					} else {
						newList = append(newList, i)
					}
				}
				clientsList = newList
				fmt.Print("Client list-> ")
				fmt.Println(clientsList)
			case "RLAY":
				fmt.Println("\nINCOMING RELAY MESSAGE->: " + message[4:])
			default:
				fmt.Println("\nHUB MESSAGE->: " + message)
			}
		} else {
			fmt.Println("\nHUB MESSAGE->: " + message)
		}

	}

}

func msgIDME(c net.Conn) {
	fmt.Fprintf(c, "IDME\n")
}

func msgLIST(c net.Conn) {
	fmt.Fprintf(c, "LIST\n")
}

//example message data
//relay "Hello" to client id = 1 and client id = 2
//message : RLAY1,2MSGHello
func msgRLAY(c net.Conn, message string) {
	if len(message) < 5 {
		fmt.Println("Invalid Message")
		return
	}
	messageType := message[0:4]
	if messageType != "RLAY" {
		fmt.Println("Invalid Message Type")
		return
	}
	msgIndex := strings.Index(message, "MSG")
	if msgIndex < 5 {
		fmt.Println("Invalid MSG index")
		return
	}
	receiverClientIds := strings.Split(string(message[4:msgIndex]), ",")
	//fmt.Println(len(receiverClientIds))
	if len(message[msgIndex+3:]) > 1024 {
		fmt.Println("max message size should be 1024 kilobytes!")
		return
	}
	if len(receiverClientIds) > 255 {
		fmt.Println("max messgae receiver client count should be 255!")
		return
	}

	fmt.Fprintf(c, message)

}

func msgSTOP(c net.Conn) {
	fmt.Fprintf(c, "STOP\n")
}

func connectHub(CONNECT string) (net.Conn, error) {

	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return c, nil
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	//c, err := net.Dial("tcp", CONNECT)
	c, err := connectHub(CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("TCP/IP Client Started")
	go listener(c)

	fmt.Println("Press [1]	to get ID from HUB")
	fmt.Println("Press [2]	to List All Connected Clients from HUB")
	fmt.Println("Press [3]	to Relay Message")
	fmt.Println("Press [4]	to STOP Connection")
	fmt.Print("Please select message type [1-2-3-4] -> ")
	for {
		messageType, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		switch messageType[0] {
		case '1':
			fmt.Println("Get ID from HUB")
			msgIDME(c)
		case '2':
			fmt.Println("List All Connected Clients from HUB")
			msgLIST(c)
		case '3':
			fmt.Println("Relay Message")
			if len(clientsList) < 1 {
				fmt.Println("There is no other clients to Relay message")

			} else {

				fmt.Print("Clients Ids on Hub: ")
				fmt.Println(clientsList)
				fmt.Println("Create message using client ids on list, with below format ")
				fmt.Println("FORMAT: RLAY1,2,3MSGexample message")
				relayMessage, _ := bufio.NewReader(os.Stdin).ReadString('\n')
				msgRLAY(c, relayMessage)
			}

		case '4':
			fmt.Println("STOP Connection")
			msgSTOP(c)
			return
		default:
			fmt.Println("Invalid selection!")
			fmt.Println("Please select message type [1-2-3-4] ->")
		}
	}
}
