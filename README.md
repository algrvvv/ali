# ali - your aliases

ali - утилита командной строки для более удобной и быстрой работы с рутинными командами

### Установка

```shell
go install https://github.com/algrvvv/ali.git
```

# Заверешение команд

```shell
ali completion your_shell > dir/for/completion

# например
ali completion zsh > ~/tmp/ali
source ~/tmp/ali
```

После использование команды source могут быть ошибки, но не смотря на них
все может отработать правильно.

### TODO

- [ ] прокидывание команд
