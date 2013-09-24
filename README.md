Jabberbot
=========

Connects to Jabber and relays all messages through ØMQ. Handling of incoming messages can change without reconnecting to Jabber. Multiple interested applications to subscribe to incoming messages.

    Usage: Jabberbot [options]
      -password="": 
      -publisherAddress="ipc://chatReceivedEvent.ipc": Address where PUB zeromq socket will be opened. It will publish each time a chat is received in Json format of ChatReceivedEvent in contracts.
      -replyAddress="ipc://sendChatCommand.ipc": Address where REP zeromq socket will be opened. Its expecting to receive work in Json format of SendChatCommand in contracts.
      -server="talk.google.com:443": 
      -username="": Username including hostname, for google this should be @gmail.com
