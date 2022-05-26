package main

import (
	"fmt"
	"image/color"
	"strings"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/mattn/go-runewidth"
	"golang.org/x/image/font"
)

type Message struct {
	PlainText string
	Fg        color.RGBA
	Count     int
}

func NewMessage(text string, fg color.RGBA) *Message {
	return &Message{
		PlainText: text,
		Fg:        fg,
		Count:     1,
	}
}

func (m Message) FullText() string {
	if m.Count > 1 {
		return fmt.Sprintf("%s (x%d)", m.PlainText, m.Count)
	}
	return m.PlainText
}

type MessageLog struct {
	Messages []*Message
}

func NewMessageLog() *MessageLog {
	return &MessageLog{
		Messages: []*Message{},
	}
}

func (m *MessageLog) AddMessage(text string, fg color.RGBA, stack bool) {
	if stack && len(m.Messages) > 0 && text == m.Messages[len(m.Messages)-1].PlainText {
		m.Messages[len(m.Messages)-1].Count += 1
	} else {
		m.Messages = append(m.Messages, NewMessage(text, fg))
	}
}

func (m MessageLog) Render(screen *ebiten.Image, f font.Face, x, y, width, height int) {
	renderMessages(screen, f, x, y, width, height, m.Messages)
}

func wrap(str string, width int) chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		wrapped := strings.Split(runewidth.Wrap(str, width), "\n")
		for i := 0; i < len(wrapped)/2; i++ {
			wrapped[i], wrapped[len(wrapped)-i-1] = wrapped[len(wrapped)-i-1], wrapped[i]
		}
		for _, line := range wrapped {
			ch <- line
		}
	}()

	return ch
}

func renderMessages(screen *ebiten.Image, f font.Face, x, y, width, height int, messages []*Message) {
	yOffset := height - 1

	reversedMsg := make([]Message, 0, len(messages))
	for _, msg := range messages {
		reversedMsg = append(reversedMsg, *msg)
	}
	for i := 0; i < len(reversedMsg)/2; i++ {
		reversedMsg[i], reversedMsg[len(reversedMsg)-i-1] = reversedMsg[len(reversedMsg)-i-1], reversedMsg[i]
	}
	for _, msg := range reversedMsg {
		for line := range wrap(msg.FullText(), width) {
			text.Draw(screen, line, f, x*10, (y+yOffset)*10, msg.Fg)
			yOffset -= 1
			if yOffset < 0 {
				return
			}
		}
	}
}
