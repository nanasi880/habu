# Overview
The habu is [Cobra](https://github.com/spf13/cobra) 's very lightweight router.  
Cobra requires command instances to be bound with AddCommand, but it reduces this effort.  
If you want to separate packages for different subcommands, you don't have to worry about circular references to parent and child commands.

# Example
Example of matching a subcommand with a directory.

#### my/package/cmd/sub/xxx.go
```golang
package sub

import (
    "github.com/spf13/cobra"
    "github.com/nanasi880/habu"
)

var command = &cobra.Command {
    Use: "sub",
}

func init () {
    habu.AddCommand(command, "/root")
}
```

#### my/package/cmd/xxx.go
```golang
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/nanasi880/habu"
    _ "my/package/cmd/sub"
)

var rootCommand = &cobra.Command {
    Use: "root",
}

func init () {
    habu.AddCommand(rootCommand, "/")
}
```

#### my/package/main.go
```golang
package main

import (
    "os"
    
    "github.com/nanasi880/habu"
    _ "my/package/cmd"
)

func main() {
    err := habu.Execute()
    if err != nil {
        os.Exit(1)
    }
}
```