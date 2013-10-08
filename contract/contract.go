package contract

type SendChatCommand struct {
	To   string
	Text string
}

type ChatReceivedEvent struct {
	From string
	Text string
}
