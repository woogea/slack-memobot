package mention

import "strings"

type Mention struct {
	storage Storage
	command []Command
}

func NewMention(namespace string) *Mention {
	ret := &Mention{storage: NewMemoryStorage(namespace)}
	ret.initCommand()
	return ret
}

func (m *Mention) Exec(cmd string, text string) (string, error) {
	i := findFirst(m.command, cmd)
	if i != -1 {
		return m.command[i].Action(m, text)
	} else {
		//exec get when the command is not found.
		//return "", &NotfoundError{}
		i = findFirst(m.command, "choice")
		text = strings.TrimSpace(cmd + " " + text)
		return m.command[i].Action(m, text)
	}
}

func findFirst(cmd []Command, key string) int {
	var result = -1
	for i := range cmd {
		if cmd[i].Name == key {
			result = i
		}
	}
	return result
}
