package socketManager

// Messenger describes an entity that can parse messages from a client and optionally transform them.
type Messenger interface {
	OnOpen() (message Message, isError bool, shouldBroadcast bool)                     // OnOpen() sends a message when a client is opened. If the message is erroneous, isError will return true
	Parse(message Message) (unmarshalledMessage Message, isError bool, shouldBroadcast bool) // Parse() parses a message recieved from the client. After performing side effects, pass the message along. If the message is erroneous, isError will return true
}
