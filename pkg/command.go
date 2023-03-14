package mention

import (
	"fmt"
	"math/rand"
	"strings"
)

func echo(m *Mention, text string) (string, error) {
	return text, nil
}

func set(m *Mention, text string) (string, error) {
	items := strings.Split(text, " ")
	if len(items) < 1 {
		return "", &InvalidArgumentError{"empty value"}
	}
	val := strings.Join(items[1:], " ")
	val = strings.TrimSpace(val)
	return "", m.storage.Add(items[0], val)
}

func get(m *Mention, text string) (string, error) {
	text = strings.TrimSpace(text)
	return m.storage.Get(text)
}

func list(m *Mention, text string) (string, error) {
	l, err := m.storage.List(text)
	if err != nil {
		return "", err
	}
	// looking good with first newline.
	return "\n" + strings.Join(l, "\n"), nil
}

func rotate(m *Mention, text string) (string, error) {
	return "", m.storage.Rotate(text)
}

func choice(m *Mention, text string) (string, error) {
	l, err := m.storage.List(text)
	if err != nil {
		return "", err
	}
	// choice one item
	return l[rand.Intn(len(l))], nil
}

func remove(m *Mention, text string) (string, error) {
	items := strings.Split(text, " ")
	if len(items) < 1 {
		return "", &InvalidArgumentError{"empty value"}
	}
	return "", m.storage.Remove(items[0], strings.Join(items[1:], " "))
}

func removeAll(m *Mention, text string) (string, error) {
	return "", m.storage.RemoveAll(text)
}

func funcList(m *Mention, text string) (string, error) {
	ret := ""
	for _, x := range m.command {
		ret = ret + fmt.Sprintf("%s: %s\n", x.Name, x.Help)
	}
	return ret, nil
}

type Command struct {
	Name      string
	ShortName []string
	Action    func(m *Mention, text string) (string, error)
	Help      string
}

func (m *Mention) initCommand() {
	commanddef := []Command{
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
		{
			Name:      "list",
			ShortName: []string{"l"},
			Action:    list,
			Help:      "list <key>",
		},
		{
			Name:      "remove",
			ShortName: []string{"r"},
			Action:    remove,
			Help:      "remove <key> <value>",
		},
		{
			Name:      "rotate",
			ShortName: []string{"rot"},
			Action:    rotate,
			Help:      "rotate <key>",
		},
		{
			Name:      "choice",
			ShortName: []string{"c"},
			Action:    choice,
			Help:      "choice <key>",
		},
		{
			Name:      "removeAll",
			ShortName: []string{"ra"},
			Action:    removeAll,
			Help:      "removeAll <key>",
		},
		{
			Name:      "help",
			ShortName: []string{"ra"},
			Action:    funcList,
			Help:      "help",
		},
	}
	for _, x := range commanddef {
		m.command = append(m.command, x)
	}

}
