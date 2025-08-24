# ubuntu mainline kernel tool
Ubuntu mainline kernel tool предназначен для загрузки пакетов свежих ядер с https://kernel.ubuntu.com/mainline/. Новые версии ядер для Ubuntu выходят примерно синхронно с ядрами kernel.org. Утилита полностью консольная, предназначена для загрузки только свежих версий (а не любой из списка).

!!! На текущий момент - это MVP, работает минимальная функциональность. Кому надо больше или если нашли баг - заводите task. !!!

## как установить
Установить из PPA: https://launchpad.net/~alex19srv/+archive/ubuntu/ppa
```bash
sudo add-apt-repository ppa:alex19srv/ppa
sudo apt update
sudo apt install mainline-kernel-tool
```
или установить из исходников
```bash
git clone https://github.com/alex19srv/mainline-kernel-tool.git
cd mainline-kernel-tool
# нужно будет установить just и go
just build
# или без just
#  go build -o ./cmd/mainline_kernel/
```

## как использовать
```bash
mainline_kernel check [installed_version]
```
проверяет последнюю доступную версию и сравнивает с указанной в командной строке или среди установленных (`dpkg-query`). Если есть более свежая версия - выводит ее номер.

!!! в командной строке указывается установленная версия !!!
```bash
# Example:
mainline-kernel check 6.16.1
# will output:
# latest version 6.16.3
```

```bash
mainline_kernel fetch-latest [latest_installed_version] [--dir=target_dir] [--force]
```
проверяет, что доступна новая версия и загружает пакеты с ядром и модулями. Текущая установленная берется из командной строки или проверяется по списку установленных пакетов.

Если не указана директория `--dir` - загружает в текущую.

`--force` - загрузить последнюю версию без сравнения с установленными.

```bash
# example
mainline-kernel fetch-latest 6.16.1 --dir=/tmp/kernels --force
# выведет список загруженных файлов, для которых нужно сделать dpkg -i
```

## известные недоработки, которые я в ближайшее время исправлю
- При установке из пакета команда называется `mainline_kernel`. Из исходников `mainline-kernel`
- Парсер версии, если дать версию в неверном формате падает с непонятным выводом.
- удаление старых версий пока не реализовано.
