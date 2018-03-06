package discordbot

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Commands      map[string]Command
	caseSensative bool
	Num           int
	Discord       *discordgo.Session
	MessagePtr    *discordgo.MessageCreate
	TestChannelID string
	LogChannelID  string
}

type CreateMessage func(*discordgo.Session, *discordgo.MessageCreate)
type Command func(*Bot, string)

var idx int = 0

func New(token string) *Bot {
	b := Bot{}
	b.Commands = make(map[string]Command, 0)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	dg.AddHandler(b.messageCreate)
	dg.Open()
	b.Discord = dg
	b.caseSensative = false
	return &b
}

func (b *Bot) Close() {
	b.SendLogMessage("Shutting down bot...")
	b.Discord.Close()
}

func (b *Bot) EvaluateMessage(msg string) (Command, string) {
	for prefix, fn := range b.Commands {
		found := true
		prfxwords := strings.Split(prefix, " ")
		msgwords := strings.Split(msg, " ")
		var i int
		var word string
		for i, word = range prfxwords {
			m := msgwords[i]
			if !b.caseSensative {
				m = strings.ToLower(m)
				word = strings.ToLower(word)
			}
			if word != m {
				found = false
				break
			}
		}
		if found {
			return fn, strings.Join(msgwords[i+1:], " ")
		}
	}
	return nil, ""
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	fn, args := b.EvaluateMessage(m.Content)
	if fn != nil {
		b.Discord = s
		b.MessagePtr = m
		fn(b, args)
	}

}

func (b *Bot) AddCommand(key string, cmd Command) {
	if b.Commands[key] == nil {
		m := make(map[string]Command, len(b.Commands))
		m[key] = cmd
		for k, v := range b.Commands {
			m[k] = v
		}
		b.Commands = m
	}
}

func (b *Bot) RemoveCommand(key string) {
	if b.Commands[key] != nil {
		delete(b.Commands, key)
	}
}

func (b *Bot) SendMessage(msg string) {
	b.Discord.ChannelMessageSend(b.MessagePtr.ChannelID, msg)
}

func (b *Bot) SetTestChannel(id string) {
	b.TestChannelID = id
}

func (b *Bot) SetLogChannel(id string) {
	b.LogChannelID = id
}

func (b *Bot) SendTestMessage(msg string) {
	b.Discord.ChannelMessageSend(b.TestChannelID, msg)
}

func (b *Bot) SendLogMessage(msg string) {
	logMsg := time.Now().Format("[Mon Jan 2 15:04:05.999 MST 2006]: ") + msg
	b.Discord.ChannelMessageSend(b.TestChannelID, logMsg)
}

func (b *Bot) Error(err error) {
	b.SendLogMessage(err.Error())
}

func (b *Bot) GetMessageAuthorID() string {
	return b.MessagePtr.Author.ID
}

func (b *Bot) IsCaseSensative(state bool) {
	b.caseSensative = state
}
