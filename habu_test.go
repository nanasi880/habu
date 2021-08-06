package habu_test

import (
	"testing"

	"github.com/nanasi880/habu"
	"github.com/spf13/cobra"
)

func TestHabu(t *testing.T) {
	commands := []struct {
		name string
		path string
	}{
		{
			name: "root",
			path: "/",
		},
		{
			name: "sub1",
			path: "/root",
		},
		{
			name: "sub2",
			path: "/root/sub1",
		},
		{
			name: "sub3",
			path: "/root",
		},
		{
			name: "sub4",
			path: "/root/sub3",
		},
	}

	inst := new(habu.Habu)
	for _, cmd := range commands {
		c := &cobra.Command{
			Use: cmd.name,
			RunE: func(c *cobra.Command, _ []string) error {
				return c.Help()
			},
		}
		err := inst.AddCommand(c, cmd.path)
		if err != nil {
			t.Fatal(err)
		}
	}

	roots, err := inst.ToCobra()
	if err != nil {
		t.Fatal(err)
	}
	root := roots[0]

	subs := root.Commands()
	if len(subs) != 2 {
		t.Fatal()
	}
	if subs[0].Name() != "sub1" {
		t.Fatal(subs[0].Name())
	}
	if subs[1].Name() != "sub3" {
		t.Fatal(subs[1].Name())
	}
	sub1 := subs[0]
	if len(sub1.Commands()) != 1 {
		t.Fatal()
	}
	if sub1.Commands()[0].Name() != "sub2" {
		t.Fatal(sub1.Commands()[0].Name())
	}
	sub3 := subs[1]
	if len(sub3.Commands()) != 1 {
		t.Fatal()
	}
	if sub3.Commands()[0].Name() != "sub4" {
		t.Fatal(sub3.Commands()[0].Name())
	}
}

func TestHabu_CreateIntermediateCommands(t *testing.T) {
	c := &cobra.Command{
		Use: "command",
		RunE: func(c *cobra.Command, _ []string) error {
			return c.Help()
		},
	}

	inst := new(habu.Habu)
	err := inst.AddCommand(c, "/path/to/")
	if err != nil {
		t.Fatal(err)
	}

	roots, err := inst.ToCobra(habu.CreateIntermediateCommands(true))
	if err != nil {
		t.Fatal(err)
	}
	root := roots[0]

	if root.Name() != "path" {
		t.Fatal(root.Name())
	}
	if len(root.Commands()) != 1 {
		t.Fatal()
	}
	if root.Commands()[0].Name() != "to" {
		t.Fatal(root.Commands()[0].Name())
	}
	if len(root.Commands()[0].Commands()) != 1 {
		t.Fatal()
	}
	if root.Commands()[0].Commands()[0].Name() != "command" {
		t.Fatal(root.Commands()[0].Name())
	}
}
