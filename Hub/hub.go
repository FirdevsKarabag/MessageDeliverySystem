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
	c.Write([]byte(string("RLAY" + message)))
}

//Incoming Message Formats
//IDME : returns unique client id
//LIST : returns all connected client ids
//RLAY : relay message to specific clients list
//STOP : end connection
func handleConnection(c net.Conn) {
	clientId := count
	fmt.Println("new client connected d: ", clientId)
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			delete(clientsMap, clientId)
			c.Close()
			return
		}

		var index int = 0
		var responseMessage string

		// message format control
		if len(netData) < 5 {
			responseMessage = "Invalid Message Type"
		} else {
			messageType := netData[0:4]
			index += 4
			switch messageType {
			case "IDME":
				responseMessage = "IDME" + strconv.Itoa(clientId)
			case "LIST":
				var idList = ""
				for id := range clientsMap {
					if id == clientId {
						continue
					}
					idList += strconv.Itoa(id) + ","
				}
				if idList != "" {
					//remove last ","
					idList = idList[:len(idList)-1]
					responseMessage = "LIST" + idList
				} else {
					responseMessage = "No other clients!.."
				}

			//RELAY MESSAGE FORMAT
			//4 byte message type : RLAY
			//comma seperated receiver client ids
			//message header : MSG
			//message text
			case "RLAY":

				messageIndex := strings.Index(netData, "MSG")
				message2send := netData[messageIndex+3:]
				receiverClientIds := strings.Split(string(netData[index:messageIndex]), ",")
				fmt.Println(message2send)
				responseMessage = ""
				unFoundIds := ""

				for _, id := range receiverClientIds {
					i, err := strconv.Atoi(id)
					if err != nil {
						fmt.Println(err)
						//responseMessage = "non integer client id"
					} else {
						if connectionInfo, ok := clientsMap[i]; ok {
							go sendMessage(connectionInfo, message2send)
						} else {
							fmt.Printf("client id not found: %s\n", id)
							unFoundIds += id + " "
						}

					}
				}
				if len(unFoundIds) > 1 {
					responseMessage = "Message Relay Triggered! except clients: " + unFoundIds

				} else {
					responseMessage += "Message Relay Triggered"
				}

			case "STOP":
				delete(clientsMap, clientId)
				c.Close()
				return
			default:
				responseMessage = "Invalid Message Type"
			}

		}

		responseMessage += "\n"
		c.Write([]byte(string(responseMessage)))
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
		//fmt.Println("connectionInfo:", clientsMa)

	}
}
