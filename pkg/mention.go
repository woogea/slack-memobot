package mention

type Mention struct {
	storage Storage
}

func NewMention(namespace string) *Mention {
	return &Mention{storage: NewMemoryStorage(namespace)}
}

func (m *Mention) Exec(cmd string, text string) (string, error) {
	i := findFirst(command, cmd)
	if i != -1 {
		return command[i].Action(m, text)
	} else {
		//exec get when the command is not found.
		//return "", &NotfoundError{}
		i := findFirst(command, "get")
		return command[i].Action(m, cmd+" "+text)
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
