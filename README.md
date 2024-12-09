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
# строку ниже стоит добавить в ~/.zshrc или другой файл вашей оболочки
source ~/tmp/ali
```

### TODO

- [x] прокидывание команд
- [ ] написать доку
