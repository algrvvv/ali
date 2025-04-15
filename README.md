# ali - your aliases

ali - cli utility for more convenient and quick work with routine teams

### Install

intall:
```shell
go install github.com/algrvvv/ali@latest
```

setup global config
```shell
# init global config
ali setup
```

you can also use the `--editor` flag to specify the default config editor.

in the example below, all configurations will be edited using 
the `ali edit` command using the `vs code` code editor
```shell
ali setup --editor=code
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
  ali - cli app for your aliases
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
    -D, --debug                 print debug messages
    -h, --help                  help for ali
    -L, --local-env             use only local env
        --output-color string   color of the ouput of the parallel command
    -p, --parallel              do parallel command
        --print                 print result command before start exec
        --without-output        dont show parallel commands output

  Use "ali [command] --help" for more information about a command.
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

or yaml configuration example (you can read about dynamic file expansion [here](#local-configuration-dynamic-extensions)):
```yaml
aliases:
  test: 'hello world'
parallel:
  lara:
    - label: server
      color: red
      command: php artisan serve
      path: .
    - label: vue
      color: green
      command: npm run dev 
      path: .
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

### Local configuration (dynamic extensions)

Recently, the ability to use multiple file formats for the local configuration has been added. `toml`, `yaml`, `json`. the file should also be named `.ali`, the utility recognizes the file extension on its own.


> global configuration not supported it

the default file extension is `toml`. 
to create in a different format, you can use 
```shell
ali init --format yaml
``` 

or using short flag and other config type.
```shell
ali init -F json
```

> at any time, you can rewrite your local configuration file to a different format and it won't break anything. 
> To do this, you can use online converters or do it manually.

### Variables

since version `v1.6.3` it is now possible to create and use variables. 
they are stored in the configuration under the tag `vars`. 
They can also be divided into global or local ones.

below is an example from the local `yaml` config:
```yaml
vars:
  execute: go run main.go
  roomcode: myroomcode
  user: myuser

aliases:
  # golang
  single: '{{execute}} --user={{user}}'
  host: '{{execute}} --host --code {{roomcode}} --user={{user}}'
  guest: '{{execute}} --guest --code {{roomcode}} --user={{user}}'
```

It is also possible to reassign variables for the duration of execution, that is, not in the file, 
but only once at the time of command execution.

you can follow the example below:

```shell
ali single -V_execute=mybuild --print
```

using the `--print` flag, we will see the final command to execute. 
the output before execution will be as follows:
```shell
mybuild --user=myuser
```

how does it work? 
to reassign a variable, you need to use `V_` at the beginning of the flag,
and then use the variable name, followed by its value.

> important! it is necessary to pass the value through the `=` sign.


To see their list, you can use the `ali list -v` command.

### throwing flags or values

throwing flags that are not used directly or by substituting an argument into a command works as follows.
We have an alias:
```yaml
aliases:
  run: go run main.go
```

and we want to add a flag that is missing there, but we don't want to add it to the file, since we won't need it later. 
in this case, we will use the following command (we will also use `--print` to see the final command)
```shell
ali run --myarg=myvalue --print
```

ouput:
```text
command:  go run main.go --myarg=myvalue
```

> important! it is necessary to pass the value through the `=` sign.

### Additionally

To get logs, use `--debug` or `-D`

To use only local aliases, use `-L` or `--local-env`.
This can be useful when using `ali list -L` to output only a list of local aliases.

### TODO

прокидывание флагов, которые не используются напрямую или подменой аргумента в команду.
добавить логирование итоговой команды перед ее запуском.
