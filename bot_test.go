package discordbot_test

import (
	"fmt"
	"testing"

	discordbot "github.com/ltsnuggie/discord-bot"
)

func TestMain(t *testing.T) {
	b := discordbot.New()
	b.AddCommand("test", Cascade)
	b.AddCommand("testing", SampleCommand)

}

func Cascade(b *discordbot.Bot, args string) {
	fn, args := b.EvaluateMessage(args)
	if fn != nil {
		fn(b, args)
	}
}

func SampleCommand(b *discordbot.Bot, args string) {
	fmt.Println("SampleCommand")
	b.SendMessage(args)
}
