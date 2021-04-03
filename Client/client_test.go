package main

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

var hubConnInfo string

var clientsMap map[int]net.Conn

func TestNewClient(t *testing.T) {
	hubConnInfo = "127.0.0.1:5200"
	clientsMap = make(map[int]net.Conn)

	for i := 0; i < 300; i++ {
		c, err := connectHub(hubConnInfo)
		if err != nil {
			t.Errorf("connection error occured :%v", err)
		}
		clientsMap[i] = c

		fmt.Println("TCP/IP Client Started", i)

		go listener(c)

	}

	//time.Sleep(2 * time.Second)
	// client id values
	for i := 0; i < 300; i++ {
		msgIDME(clientsMap[i])
	}

	// full client list
	msgIDME(clientsMap[0])
	msgLIST(clientsMap[0])
	//wait for response
	time.Sleep(1 * time.Second)
	fmt.Println("****")
	fmt.Println(clientsList)
	fmt.Println("****")

	//message relay tests
	relayMessage := fmt.Sprintf("RLAY%d,%d,%d,%dMSGFirst message--for 4 clients☻\n",
		clientsList[1], clientsList[2], clientsList[3], clientsList[4])
	//fmt.Println(relayMessage)
	msgRLAY(clientsMap[0], relayMessage)
	time.Sleep(3 * time.Second)
	relayMessage = fmt.Sprintf("RLAY%d,%d,%d,%d,%d,%d,%d,%dMSGSecond message♦For 8 clients\n",
		clientsList[1], clientsList[2], clientsList[3], clientsList[4], clientsList[5], clientsList[6], clientsList[7], clientsList[8])
	msgRLAY(clientsMap[0], relayMessage)
	time.Sleep(3 * time.Second)
	relayMessage = fmt.Sprintf("RLAY%d,%d,%d,%d,%d,%dMSGThird message♣For 6 clients\n",
		clientsList[11], clientsList[21], clientsList[13], clientsList[14], clientsList[15], clientsList[16])
	msgRLAY(clientsMap[0], relayMessage)
	time.Sleep(3 * time.Second)

	//relay message format test
	relayMessage = fmt.Sprintf("%d,%d,%d,%d,%d,%dMSGNo messgae type test message♣For 6 clients\n",
		clientsList[11], clientsList[21], clientsList[13], clientsList[14], clientsList[15], clientsList[16])
	msgRLAY(clientsMap[0], relayMessage)
	time.Sleep(3 * time.Second)

	//relay message size test 1023 KB
	relayMessage = fmt.Sprintf("RLAY%dMSGLorem ipsum dolor sit amet, consectetur adipiscing elit. Ut lacinia metus ac volutpat iaculis. Fusce imperdiet id neque vitae dapibus. Duis urna massa, semper sed nisl et, vestibulum dapibus nunc. Donec finibus tempus lorem, sed efficitur nisl euismod quis. Sed faucibus, nibh et tincidunt imperdiet, leo tellus placerat neque, sed vulputate ligula orci id magna. In volutpat elit vitae mattis posuere. Praesent a augue ut neque posuere consequat nec in magna.Vivamus id neque augue. Nulla dolor leo, elementum vel tristique sit amet, mattis non risus. Fusce malesuada a erat in aliquet. Ut sed erat eget mauris aliquet blandit. Aenean egestas risus et neque mollis tincidunt. Proin eget justo porta purus efficitur vestibulum. Nullam mauris enim, convallis et enim mattis, tristique blandit sem.Suspendisse nisi ipsum, viverra nec massa ut, tincidunt fringilla orci. Nunc leo metus, sollicitudin in tempus venenatis, ullamcorper quis ex. In hac habitasse platea dictumst. Vivamus urna dui, rhoncus eget ultrices eget nulla\n",
		clientsList[1])
	msgRLAY(clientsMap[0], relayMessage)
	time.Sleep(3 * time.Second)

	//relay message size test 1024 KB
	relayMessage = fmt.Sprintf("RLAY%dMSGLorem ipsum dolor sit amet, consectetur adipiscing elit. Ut lacinia metus ac volutpat iaculis. Fusce imperdiet id neque vitae dapibus. Duis urna massa, semper sed nisl et, vestibulum dapibus nunc. Donec finibus tempus lorem, sed efficitur nisl euismod quis. Sed faucibus, nibh et tincidunt imperdiet, leo tellus placerat neque, sed vulputate ligula orci id magna. In volutpat elit vitae mattis posuere. Praesent a augue ut neque posuere consequat nec in magna.Vivamus id neque augue. Nulla dolor leo, elementum vel tristique sit amet, mattis non risus. Fusce malesuada a erat in aliquet. Ut sed erat eget mauris aliquet blandit. Aenean egestas risus et neque mollis tincidunt. Proin eget justo porta purus efficitur vestibulum. Nullam mauris enim, convallis et enim mattis, tristique blandit sem.Suspendisse nisi ipsum, viverra nec massa ut, tincidunt fringilla orci. Nunc leo metus, sollicitudin in tempus venenatis, ullamcorper quis ex. In hac habitasse platea dictumst. Vivamus urna dui, rhoncus eget ultrices eget nulla.\n",
		clientsList[1])
	msgRLAY(clientsMap[0], relayMessage)
	time.Sleep(3 * time.Second)

	//relay message receiver client count control -- 255 clients
	client_ids := ""
	for i := 0; i < 255; i++ {
		client_ids += fmt.Sprintf("%d,", clientsList[i+1])
	}
	client_ids = strings.TrimSuffix(client_ids, ",") //remove last separator of id list
	relayMessage = fmt.Sprintf("RLAY%sMSGThis message will relay to 255 clients!☻\n", client_ids)
	msgRLAY(clientsMap[0], relayMessage)
	time.Sleep(10 * time.Second)

	//relay message receiver client count control -- no more than 255
	client_ids = ""
	for i := 0; i < 256; i++ {
		client_ids += fmt.Sprintf("%d,", clientsList[i+1])
	}
	client_ids = strings.TrimSuffix(client_ids, ",") //remove last separator of id list
	relayMessage = fmt.Sprintf("RLAY%sMSGThis message will relay to 255 clients!☻\n", client_ids)
	msgRLAY(clientsMap[0], relayMessage)
	time.Sleep(5 * time.Second)

	time.Sleep(5 * time.Second)

}
