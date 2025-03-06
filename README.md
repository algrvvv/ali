# ali - your aliases

ali - cli utility for more convenient and quick work with routine teams

### Install

```shell
go install github.com/algrvvv/ali@latest

# init global config
ali setup
```

### Auto completion

```shell
ali completion your_shell > dir/for/completion

# for example
ali completion zsh > ~/tmp/ali
source ~/tmp/ali
# It's better to put this command in your shell configuration, for example ~/.zshrc
```

### Usage

The application has a global configuration and any number of local configurations.
To do this, go to the desired directory and use the command: `ali init`.
All local overlapping aliases have an advantage over global ones.

```bash
> ali help
    Usage:
      ali [flags]
      ali [command]

    Available Commands:
      completion  Generate completion script
      edit        Edit global or local config
      help        Help about any command
      init        Init new local config
      list        Get list aliases
      setup       Setup global config
      version     See app version and more information

    Flags:
      -D, --debug       print debug messages
      -h, --help        help for ali
      -L, --local-env   use only local env
```

### App configuration

The application configuration is stored in a file with the `.toml` extension.
The `aliases` section contains a list of the type `alias = command` and there is
also the `app` section, which contains so far the only setting that is responsible
for the default configuration editor.

To edit the global configuration, use: `ali edit`
To edit the local configuration, use: `ali edit --local`
By default, `vi` opens to edit the configuration.

Example:

```toml
[aliases]
# test alias
test = 'echo "hello world"'

[app]
editor = 'vim'
```

### Usage examples

Pass arguments inside a command:

```toml
# in configuration
hello = 'echo "hello, <user>"'
```

```bash
ali hello --user=$(whoami)
# equal: echo "hello, $(whoami)"
```

passing arguments after the command:

```toml
# in configuration
gl = 'git log -n '
```

```bash
ali gl 3
# equal: git log -n 3
```

### Multiple commands

Multiple commands are commands that will be executed in parallel in a single
terminal session.

Configuration example:

```toml
[[parallel.lara]]
# name of command
label='Laravel'
# color of the label in logs
color='red'
# command for execute
command='php artisan serve'
# work directory
path='.'

[[parallel.lara]]
label='Vue'
color='green'
command='npm run dev'
path='.'
```

To run parallel commands, use `ali CommandName -p` or `ali commandName --parallel`

There are also additional settings for parallel commands.
For example, the `--without-output` flag to disable command output.
It is also possible to change the output color by using flag `--output-color`.

Available colors:

- red
- green
- yellow
- blue
- magenta
- cyan
- gray
- orange
- pink
- lime
- white

At the moment, these commands will only be displayed in the list using the`-f` flag.

### Additionally

To get logs, use `--debug` or `-D`

To use only local aliases, use `-L` or `--local-env`.
This can be useful when using `ali list -L` to output only a list of local aliases.
