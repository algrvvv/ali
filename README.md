# ali - your aliases

ali - cli utility for more convenient and quick work with routine teams

## about v2

- Improved alias settings. added more flexible settings.
- The second version introduced plugins, templates, and a complete switch to using `yaml` configurations. support for `toml` and `json` configurations has been removed.

### Install

install by `go install`:

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

The application configuration is stored in a file with the `.yaml` extension.
The `aliases` section contains a list of the type `alias = command` and there is
also the `app` section, which contains so far the only setting that is responsible
for the default configuration editor.

To edit the global configuration, use: `ali edit`
To edit the local configuration, use: `ali edit --local`
By default, `vi` opens to edit the configuration.

Example:

```yaml
aliases:
  # test alias
  test: 'echo "hello world"'
app:
  editor: 'vim'
```

### More flexibility for aliases

Example of a more flexible setup:

```yaml
aliases:
  lara:
    parallel: true # commands in alias exec parallel
    aliases: [l, laravel] # also can use: ali laravel
    desc: start laravel server + vuejs # desc for alias
    dir: ~/path/to/project # dir for exec alias
    env: 
      SERVER_PORT: 8888 # env variable for this alias
    cmds:
      - php artisan serve --port $SERVER_PORT
      - npm run dev
```

### More settings

Example of additional settings.

#### includes

Loading a different configuration. it is possible to download several by specifying several paths.

Recursive loading is currently not supported.

```yaml
include:
  - ~/some/path/
  - ~/path/some/.ali

aliases:
  t: echo 'test alias'
```

#### env

You can also add both global environment variables and local ones, that is, only for a specific alias.

```yaml
env:
  NAME: name

aliases:
  t: 
    env: 
      NAME: another_name
    cmds:
      - echo "hello, $NAME"
```

### Usage examples

Pass arguments inside a command:

```yaml
# in configuration
hello: echo "hello, <user>"
```

```bash
ali hello --user=$(whoami)
# equal: echo "hello, $(whoami)"
```

passing arguments after the command:

```yaml
# in configuration
gl: git log -n 
```

```bash
ali gl 3
# equal: git log -n 3
```

### Parallel commands

Parallel commands are commands that will be executed in parallel in a single
terminal session.

Configuration example:

```yaml
aliases:
  lara:
    parallel: true # use for parallel exec
    cmds:
      - php artisan serve --port $SERVER_PORT
      - npm run dev
```

To run parallel commands, use `ali CommandName` or `ali commandName`

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
  single: "{{execute}} --user={{user}}"
  host: "{{execute}} --host --code {{roomcode}} --user={{user}}"
  guest: "{{execute}} --guest --code {{roomcode}} --user={{user}}"
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

### Templates

Templates are prepared examples of `.ali` configurations that can be created or overwritten into an existing local config.

create new template:
```shell
ali templ testTempl --new
```

edit template:
```shell
ali templ testTempl --edit
```

import template and create local config:
```shell
ali templ testTempl
```

if you want to rewrite existed config:
```shell
ali templ testTempl --force
```

show templ list:
```shell
ali templ --list
```

### Plugins

Plugins are slightly different from templates. unlike them, they are provided so that you can run a script or any other action with a single command.

create new plugin:
```shell
ali plug testPlug --new
```

After creation, you will be given the path where the plugin should be placed. initially, there will be only its configuration file, in which it is enough to specify the description of the plugin, as well as how to launch it. This approach will add even more flexibility and will not be tied to a specific language.

plugin configuration example:
```yaml
desc: test python plugin
exec: python3 main.py
```

you can use anything in this directory, the main thing is to specify how your plugin will run.

show plugins list:
```shell
ali plug --list
```

### Additionally

To get logs, use `--debug` or `-D`

To use only local aliases, use `-L` or `--local-env`.
This can be useful when using `ali list -L` to output only a list of local aliases.

### TODO

прокидывание флагов, которые не используются напрямую или подменой аргумента в команду.
добавить логирование итоговой команды перед ее запуском.

### feature v2

- [x] добавить разные варианты алиасов;
- [x] добавить разделение команд на несколько, а не использование `&&`;
- [x] перейти в целом на yaml конфиг;
- [x] синонимы к алиасам;
- [x] правильный вывод всех алиасов;
- [x] можно добавить версию приложения в самое начало (нах);
- [x] переменные окружения;
- [x] поиск по синонимам - сначала по точным совпаденияем, а только потом по синонимам;
- [x] шаблоны;
- [x] директория выполнения команды (dir);
- [x] удалить возможность изменять путь к локальному конфигу;
- [x] include;
- [ ] почистить не используемый код;
- [ ] улучшить документацию команд именно в утилите, а не в доках
- [ ] дописать документацию в ридми
