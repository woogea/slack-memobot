package mention

func Mention(cmd string, text string) (string, error) {
	i := findFirst(command, cmd)
	if i != -1 {
		return command[i].Action(text)
	} else {
		return "", &NotfoundError{}
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
