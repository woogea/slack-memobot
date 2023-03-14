package mention

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	enabled bool
	client  *openai.Client
	rbuf    *RingBuffer
}

// chat struct for string and boolean
type Chat struct {
	chat  string
	isBot bool
}

// ring buffer struct for Chat
type RingBuffer struct {
	buf  []Chat
	head int
	tail int
}

// create new ring buffer
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buf: make([]Chat, size),
	}
}

// add new chat to ring buffer
func (r *RingBuffer) Add(chat Chat) {
	r.buf[r.head] = chat
	r.head = (r.head + 1) % len(r.buf)
	if r.head == r.tail {
		r.tail = (r.tail + 1) % len(r.buf)
	}
}

// get length of from head to tail
func (r *RingBuffer) Length() int {
	if r.head >= r.tail {
		return r.head - r.tail
	}
	return len(r.buf) - (r.tail - r.head)
}

// get last chat from ring buffer
func (r *RingBuffer) Last() Chat {
	return r.buf[r.tail]
}

// get last n chats from ring buffer
func (r *RingBuffer) LastN(n int) []Chat {
	if n > len(r.buf) {
		n = len(r.buf)
	}
	if n > r.head {
		n = r.head
	}
	if n == 0 {
		return []Chat{}
	}
	if r.head > n {
		return r.buf[r.head-n : r.head]
	}
	return append(r.buf[r.head-n:], r.buf[:r.head]...)
}

// get all chats from ring buffer from head to tail
func (r *RingBuffer) All() []Chat {
	if r.head > r.tail {
		return r.buf[r.tail:r.head]
	}
	return append(r.buf[r.tail:], r.buf[:r.head]...)
}

// clear ring buffer
func (r *RingBuffer) Clear() {
	r.head = 0
	r.tail = 0
}

// initialize and create openai client
func EnableOpenapi() (o *OpenAI) {
	o = &OpenAI{
		enabled: true,
		client:  openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		rbuf:    NewRingBuffer(8),
	}
	return
}

func (o *OpenAI) DisableOpenapi() {
	o.enabled = false
	o.client = nil
	o.rbuf = nil
}

// return true if openai is enabled
func (o *OpenAI) CanUseOpenai() bool {
	return o.enabled
}

// call openai api and return response. if client is not initialized return error
func (o *OpenAI) CallOpenai(s string, isbot bool) (string, error) {
	if !o.enabled {
		return "", fmt.Errorf("openai is not enabled")
	}
	o.rbuf.Add(Chat{s, isbot})
	if isbot {
		return "", nil
	}
	// create messages
	messages := []openai.ChatCompletionMessage{}
	for _, chat := range o.rbuf.All() {
		if chat.chat == "" {
			continue
		}
		if !isbot {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: chat.chat,
			})
		} else {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: chat.chat,
			})
		}
	}
	fmt.Printf("messages: %v\n", messages)
	// create request
	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		fmt.Printf("Caht completion error: %v\n", err)
		o.rbuf.Clear()
		return "エラーが発生したので、履歴を忘れました", err
	}
	return resp.Choices[0].Message.Content, nil
}
