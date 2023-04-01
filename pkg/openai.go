package mention

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const bufferLength = 100

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

// erase half of ring buffer content length from head to tail
func (r *RingBuffer) EraseHalf() {
	if r.head > r.tail {
		r.tail = r.head - (r.head-r.tail)/2
	} else {
		r.tail = r.head - (len(r.buf)-(r.tail-r.head))/2
		if r.tail < 0 {
			r.tail += len(r.buf)
		}
	}
}

// initialize and create openai client
func EnableOpenapi() (o *OpenAI) {
	o = &OpenAI{
		enabled: true,
		client:  openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		rbuf:    NewRingBuffer(bufferLength),
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

// constのstring配列をeraseListとして定義
var eraseList = []string{"忘れて", "わすれて", "忘れてください", "わすれてください"}

// call openai api and return response. if client is not initialized return error
func (o *OpenAI) CallOpenai(s string, isbot bool) (string, error) {
	if !o.enabled {
		return "", fmt.Errorf("openai is not enabled")
	}
	// create messages
	messages := []openai.ChatCompletionMessage{}
	for _, chat := range o.rbuf.All() {
		if chat.chat == "" {
			continue
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: chat.chat,
		})
	}
	o.rbuf.Add(Chat{s, isbot})
	if isbot {
		return "", nil
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: s,
	})
	// eraseキーワードリストの中にsが一致するものがあれば履歴を忘れる
	for _, erase := range eraseList {
		if s == erase {
			o.rbuf.Clear()
			return "履歴を忘れました", nil
		}
	}

	fmt.Printf("messages: %v\n", messages)
	// create request
	// 3回までリトライする
	cnt := 0
	var err error
	var resp openai.ChatCompletionResponse
	for cnt < 3 {
		cnt++
		// 30秒でタイムアウトするコンテキストを作成
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		fmt.Println("call api")
		resp, err = o.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: messages,
			},
		)
		if cancel != nil {
			cancel()
		}
		if err == nil {
			break
		}
	}
	if cnt == 3 {
		return "タイムアウト", err
	}
	if err != nil {
		// errに"Please reduce the length of the messages"が含まれていたら履歴を忘れる
		if condition := strings.Contains(err.Error(), "Please reduce the length of the messages"); condition {
			o.rbuf.EraseHalf()
			return "tokenが長すぎるので履歴を半分消しました", err
		}
		fmt.Printf("Caht completion error: %v\n", err)
	}
	return resp.Choices[0].Message.Content, nil
}
