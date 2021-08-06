package habu

import (
	"context"

	"github.com/spf13/cobra"
)

// Default is default habu.Habu instance.
var Default = new(Habu)

// AddCommand is add command to path location.
func AddCommand(cmd *cobra.Command, path string) error {
	return Default.AddCommand(cmd, path)
}

// MustAddCommand is add command to path location. If error happen, raise panic.
func MustAddCommand(cmd *cobra.Command, path string) {
	Default.MustAddCommand(cmd, path)
}

// ToCobra is convert to []*cobra.Command
func ToCobra(opts ...Option) ([]*cobra.Command, error) {
	return Default.ToCobra(opts...)
}

// MustToCobra is convert to []*cobra.Command, If error happen, raise panic.
func MustToCobra(opts ...Option) []*cobra.Command {
	return Default.MustToCobra(opts...)
}

// Execute is convert to cobra.Command and run first command.
func Execute(opts ...Option) error {
	return Default.Execute(opts...)
}

// ExecuteContext is convert to cobra.Command and run first command.
func ExecuteContext(ctx context.Context, opts ...Option) error {
	return Default.ExecuteContext(ctx, opts...)
}
