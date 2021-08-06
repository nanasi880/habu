package habu

type options struct {
	createIntermediateCommand bool
}

// Option is habu option.
type Option func(opt *options)

// CreateIntermediateCommands is create intermediate command just like `mkdir -p`.
func CreateIntermediateCommands(create bool) Option {
	return func(opt *options) {
		opt.createIntermediateCommand = create
	}
}
