package habu

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// Habu is command group for your application.
type Habu struct {
	commands map[string]*cobra.Command
}

// AddCommand is add command to path location.
func (h *Habu) AddCommand(cmd *cobra.Command, path string) error {
	if cmd == nil {
		return fmt.Errorf("`cmd` CANNOT be nil")
	}
	if path == "" {
		return fmt.Errorf("`path` CANNOT be empty")
	}
	if h.commands == nil {
		h.commands = make(map[string]*cobra.Command)
	}

	path = clean(path)

	fullPath := join(path, cmd.Name())
	_, found := h.commands[fullPath]
	if found {
		return fmt.Errorf("command already exists: %v", split(fullPath))
	}

	h.commands[fullPath] = cmd

	return nil
}

// MustAddCommand is add command to path location. If error happen, raise panic.
func (h *Habu) MustAddCommand(cmd *cobra.Command, path string) {
	err := h.AddCommand(cmd, path)
	if err != nil {
		panic(err)
	}
}

// ToCobra is convert to []*cobra.Command
func (h *Habu) ToCobra(opts ...Option) ([]*cobra.Command, error) {
	var opt options
	for _, o := range opts {
		o(&opt)
	}

	if opt.createIntermediateCommand {
		for _, path := range h.sortedKeys() {
			for {
				dirPath := dir(path)
				if dirPath == "/" {
					break
				}

				_, found := h.commands[dirPath]
				if !found {
					h.commands[dirPath] = stub(commandName(dirPath))
				}

				path = dirPath
			}
		}
	}

	var root []*cobra.Command
	for _, path := range h.sortedKeys() {
		cmd := h.commands[path]
		dirPath := dir(path)

		if dirPath == "/" {
			root = append(root, cmd)
			continue
		}

		parent, found := h.commands[dirPath]
		if !found {
			return nil, fmt.Errorf("parent command %v not found", split(dirPath))
		}
		parent.AddCommand(cmd)
	}

	return root, nil
}

// MustToCobra is convert to []*cobra.Command, If error happen, raise panic.
func (h *Habu) MustToCobra(opts ...Option) []*cobra.Command {
	root, err := h.ToCobra(opts...)
	if err != nil {
		panic(err)
	}
	return root
}

// Execute is convert to cobra.Command and run first command.
func (h *Habu) Execute(opts ...Option) error {
	c, err := h.ToCobra(opts...)
	if err != nil {
		return err
	}

	switch len(c) {
	case 0:
		return fmt.Errorf("command not found")
	case 1:
		return c[0].Execute()
	default:
		return fmt.Errorf("too many root command")
	}
}

// ExecuteContext is convert to cobra.Command and run first command.
func (h *Habu) ExecuteContext(ctx context.Context, opts ...Option) error {
	c, err := h.ToCobra(opts...)
	if err != nil {
		return err
	}

	switch len(c) {
	case 0:
		return fmt.Errorf("command not found")
	case 1:
		return c[0].ExecuteContext(ctx)
	default:
		return fmt.Errorf("too many root command")
	}
}

func (h *Habu) sortedKeys() []string {
	keys := make([]string, 0, len(h.commands))
	for key := range h.commands {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		a := split(keys[i])
		b := split(keys[j])

		if len(a) != len(b) {
			return len(a) < len(b)
		}

		for i := range a {
			if a[i] != b[i] {
				return a[i] < b[i]
			}
		}

		return false
	})

	return keys
}

func join(p ...string) string {
	if len(p) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(trim(p[0]))
	for _, path := range p[1:] {
		sb.WriteString("/")
		sb.WriteString(trim(path))
	}
	return clean(sb.String())
}

func trim(s string) string {
	return strings.Trim(s, "/")
}

func clean(s string) string {
	return "/" + trim(s)
}

func dir(s string) string {
	return clean(s[:strings.LastIndex(s, "/")])
}

func commandName(s string) string {
	s = clean(s)
	return s[strings.LastIndex(s, "/")+1:]
}

func split(s string) []string {
	return strings.Split(clean(s), "/")
}

func help(c *cobra.Command, _ []string) error {
	return c.Help()
}

func stub(name string) *cobra.Command {
	return &cobra.Command{
		Use:  name,
		RunE: help,
	}
}
