package main

import (
	"flag"
	"fmt"
	"github.com/mattn/go-xmpp"
	"github.com/pebbe/zmq3"
	"github.com/russellchadwick/jabberbot/contract"
	"github.com/russellchadwick/zmq"
	"log"
	"os"
)

func main() {
	server, username, password, replyAddress, publisherAddress := parseArgs()

	sendChatCommandChannel := make(chan contract.SendChatCommand, 4)
	chatReceivedEventChannel := make(chan contract.ChatReceivedEvent, 4)

	talk := connectToTalk(server, username, password)
	publisher := connectToPublisher(publisherAddress)

	go zeroMqLoop(replyAddress, sendChatCommandChannel)
	go xmppLoop(talk, chatReceivedEventChannel)
	chanLoop(sendChatCommandChannel, chatReceivedEventChannel, publisher, talk)
}

func zeroMqLoop(replyAddress *string, sendChatCommandChannel chan contract.SendChatCommand) {
	reply, err := zmq3.NewSocket(zmq3.REP)
	if err != nil {
		log.Fatal(err)
	}
	defer reply.Close()
	reply.Bind(*replyAddress)

	for {
		var sendChatCommand contract.SendChatCommand
		err := zmq.RecvJson(reply, &sendChatCommand)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Received request [%s]\n", sendChatCommand.To)
		sendChatCommandChannel <- sendChatCommand
		reply.Send("", 0)
	}
}

func xmppLoop(talk *xmpp.Client, chatReceivedEventChannel chan contract.ChatReceivedEvent) {
	for {
		chat, err := talk.Recv()
		if err != nil {
			log.Fatal(err)
		}

		switch v := chat.(type) {
		case xmpp.Chat:
			chatReceivedEvent := contract.ChatReceivedEvent{
				From: v.Remote,
				Text: v.Text,
			}

			chatReceivedEventChannel <- chatReceivedEvent

			log.Printf("Remote: [%s] Text: [%s]\n", v.Remote, v.Text)
		}
	}
}

func chanLoop(sendChatCommandChannel chan contract.SendChatCommand, chatReceivedEventChannel chan contract.ChatReceivedEvent, publisher *zmq3.Socket, talk *xmpp.Client) {
	for {
		select {
		case sendChatCommand := <-sendChatCommandChannel:
			talk.Send(xmpp.Chat{
				Remote: sendChatCommand.To,
				Type:   "chat",
				Text:   sendChatCommand.Text,
			})
		case chatReceivedEvent := <-chatReceivedEventChannel:
			zmq.SendJson(publisher, chatReceivedEvent)
		}
	}
}

func parseArgs() (server *string, username *string, password *string, replyAddress *string, publisherAddress *string) {
	server = flag.String("server", "talk.google.com:443", "")
	username = flag.String("username", "", "Username including hostname, for google this should be @gmail.com")
	password = flag.String("password", "", "")
	replyAddress = flag.String("replyAddress", "ipc://sendChatCommand.ipc", "Address where REP zeromq socket will be opened. Its expecting to receive work in Json format of SendChatCommand in contracts.")
	publisherAddress = flag.String("publisherAddress", "ipc://chatReceivedEvent.ipc", "Address where PUB zeromq socket will be opened. It will publish each time a chat is received in Json format of ChatReceivedEvent in contracts.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: Jabberbot [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()
	if *username == "" || *password == "" {
		flag.Usage()
	}

	return
}

func connectToTalk(server *string, username *string, password *string) (talk *xmpp.Client) {
	var err error

	talk, err = xmpp.NewClient(*server, *username, *password)
	if err != nil {
		log.Fatal(err)
	}

	return talk
}

func connectToPublisher(publisherAddress *string) (publisher *zmq3.Socket) {
	publisher, err := zmq3.NewSocket(zmq3.PUB)
	if err != nil {
		log.Fatal(err)
	}
	publisher.Bind(*publisherAddress)

	return publisher
}
