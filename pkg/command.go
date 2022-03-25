package mention

import (
	"strings"
)

var db = map[string]string{}

func echo(text string) (string, error) {
	return text, nil
}

func set(text string) (string, error) {
	items := strings.Split(text, " ")
	if len(items) < 1 {
		return "", &InvalidArgumentError{"empty value"}
	}
	db[items[0]] = items[1]
	return "", nil
}

func get(text string) (string, error) {
	value, ok := db[text]
	if ok {
		return value, nil
	}
	return "", &NotfoundError{}
}

type Command struct {
	Name      string
	ShortName []string
	Action    func(text string) (string, error)
	Help      string
}

var command = []Command{
	{
		Name:      "echo",
		ShortName: []string{},
		Action:    echo,
		Help:      "reutrn input text",
	},
	{
		Name:      "set",
		ShortName: []string{"s"},
		Action:    set,
		Help:      "set <key> <value>",
	},
	{
		Name:      "get",
		ShortName: []string{"s"},
		Action:    get,
		Help:      "get <key>",
	},
}
