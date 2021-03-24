package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var count = 0

var clientsMap map[int]net.Conn

func sendMessage(c net.Conn, message string) {
	c.Write([]byte(string(message)))
}

func handleConnection(c net.Conn) {
	clientId := count
	fmt.Print(".")
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		//messageType := strings.TrimSpace(string(netData))
		messageFields := strings.Split(string(netData), "|")
		messageType := messageFields[0]
		fmt.Println("MESSGAE TYPE:" + messageType)
		if messageType == "STOP" {
			delete(clientsMap, clientId)
			break
		}

		responseMessage := "RESPONSE|"
		if messageType == "ID" {
			responseMessage += "ID " + strconv.Itoa(clientId)
			//c.Write([]byte(string(message)))
		}
		if messageType == "LIST" {
			var idList = ""
			for id := range clientsMap {
				if id == clientId {
					continue
				}
				idList += strconv.Itoa(id) + "|"
			}
			responseMessage += "LIST|" + idList
			//c.Write([]byte(string(message)))
		}
		if messageType == "RELAY" {
			fmt.Printf("%q", strings.Split(string(netData), "|"))

			receiverClientIds := strings.Split(string(messageFields[1]), ",")
			message2send := messageFields[2]

			for _, id := range receiverClientIds {
				fmt.Println("***:")
				fmt.Println(id)
				i, err := strconv.Atoi(id)
				if err != nil {
					// handle error
					fmt.Println(err)
					responseMessage += "non integer client id"
				} else {
					if connectionInfo, ok := clientsMap[i]; ok {
						fmt.Println("HH")
						fmt.Println(connectionInfo)
						fmt.Println(clientsMap[i])
						fmt.Println("HH")
						go sendMessage(connectionInfo, message2send)
					} else {
						fmt.Println("client id not found")
						responseMessage += "client id not found"
					}

				}

			}
			responseMessage += "RELAY|"
			//c.Write([]byte(string(message)))
		}

		responseMessage += "\n"
		c.Write([]byte(string(responseMessage)))
		//fmt.Println(messageType)
		//counter := strconv.Itoa(count) + "\n"
		//c.Write([]byte(string(counter)))
	}
	c.Close()
}

func main() {
	clientsMap = make(map[int]net.Conn)

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		count++
		clientsMap[count] = c
		go handleConnection(c)
		fmt.Println("connectionInfo:", clientsMap)

	}
}
