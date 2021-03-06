package event

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	CONNECTED         string = "004"
	ERR_NICKNAMEINUSE string = "433"
	DISCONNECTED      string = "DISCONNECTED"
	RAW_MESSAGE       string = "RAW_MESSAGE"
	PRIVMSG           string = "PRIVMSG"
	PING              string = "PING"
	PONG              string = "PONG"
	JOIN              string = "JOIN"
	QUIT              string = "QUIT"
	NICK              string = "NICK"
	PASS              string = "PASS"
	USER              string = "USER"
	ACTION            string = "IRKKI_ACTION"
	MESSAGE           string = "IRKKI_MESSAGE"
	CTCP_VERSION      string = "IRKKI_CTCP_VERSION"
)

type Event struct {
	Raw     string
	Source  string
	User    string
	Command string
	Args    []string
}

func ParseEvent(raw string) (*Event, error) {
	evt := &Event{Raw: raw}

	var command string
	var user string
	var args []string

	if raw[0] == ':' {
		parts := strings.SplitN(raw[1:], " ", 2)
		evt.Source = parts[0]
		parts = strings.Split(parts[1], " ")
		command = parts[0]

		idx := 1
		exclamationMarkPos := strings.Index(evt.Source, "!")
		if _, err := strconv.Atoi(command); err == nil {
			user = parts[idx]
			idx++
		} else if exclamationMarkPos != -1 {
			user = evt.Source[0:exclamationMarkPos]
		}
		args = parts[idx:]
	} else {
		parts := strings.Split(raw, " ")
		command = parts[0]
		args = parts[1:]
	}
	evt.User = user
	evt.Command = command
	evt.Args = args
	return evt, nil
}

func ParseAdditionalEvents(baseEvent Event) []*Event {
	events := []*Event{}

	if baseEvent.Command == "PRIVMSG" {
		target := baseEvent.Args[0]
		message := strings.Join(baseEvent.Args[1:], " ")[1:]

		if message[0] == 1 && message[len(message) - 1] == 1 {
			message = message[1:len(message) - 1]
			parts := strings.Split(message, " ")

			if parts[0] == "ACTION" {
				baseEvent.Command = ACTION
				baseEvent.Args = []string{target, strings.Join(parts[1:], " ")}
				events = append(events, &baseEvent)
				return events
			} else if parts[0] == "VERSION" {
				baseEvent.Command = CTCP_VERSION
				baseEvent.Args = []string{target}
				events = append(events, &baseEvent)
				return events
			}
		}

		// Parse channel message from a PRIVMSG
		if target[0] == '#' {
			baseEvent.Command = MESSAGE
			baseEvent.Args = []string{target, message}
			events = append(events, &baseEvent)
			return events
		}
	}
	return events
}

func (e Event) String() string {
	return fmt.Sprintf("\x1b[32;1mSource='%s', \n\t\x1b[33;1mUser='%s', \n\t\x1b[34;1mCommand='%s', \n\t\x1b[35;1mArgs='%s', \n\t\x1b[36;1mRaw='%s'\x1b[0m",
		e.Source, e.User, e.Command, e.Args, e.Raw)
}
