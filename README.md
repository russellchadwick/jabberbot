Jabberbot
=========

Connects to Jabber and relays all messages through ØMQ. Handling of incoming messages can change without reconnecting to Jabber. Multiple interested applications to subscribe to incoming messages.

    Usage: Jabberbot [options]
      -password="": 
      -publisherAddress="ipc://chatReceivedEvent.ipc": Address where PUB zeromq socket will be opened. It will publish each time a chat is received in Json format of ChatReceivedEvent in contracts.
      -replyAddress="ipc://sendChatCommand.ipc": Address where REP zeromq socket will be opened. Its expecting to receive work in Json format of SendChatCommand in contracts.
      -server="talk.google.com:443": 
      -username="": Username including hostname, for google this should be @gmail.com

Client code to send a new message

    package main

    import (
        "github.com/russellchadwick/jabberbot/contract"
        "github.com/russellchadwick/zmq"
        "github.com/pebbe/zmq3"
    )

    func main() {
    	requester, _ := zmq3.NewSocket(zmq3.REQ)
    	defer requester.Close()
    	requester.Connect("ipc://sendChatCommand.ipc")
    	sendChatCommand := contract.SendChatCommand {
    		To:		"3w32owa6l9a9y02iqzr9w15n4e@public.talk.google.com",
    		Text:	"Hi",
    	}
    	length, _ := zmq.SendJsonNoReply(requester, sendChatCommand)
    }
