# CRYPTOOL

## Описание

**Cryptool** — это утилита командной строки для шифрования и расшифровки файлов и директорий.
Инструмент предоставляет гибкий контроль над процессом шифрования и позволяет выбрать оптимальный алгоритм для ваших задач.

## Установка (Windows)

1. Перейдите в раздел **Realeses** репозитория на Github.

2. Скачайте последнюю версию файла:
```
cryptool-windows-amd64-vX.X.X.exe
```

3. Переименуйте файл в:
```
cryptool.exe
```

4. Переместите файл в отдельную папку для CLI приложений, например:
```
C:\Users\Name\bin
```

> Рекомендуется хранить пользовательские CLI-инструменты в отдельной папке.

Чтобы запускать cryptool из любой директории, нужно добавить путь в PATH.

1. Нажмите Win + R

2. Введите
```
sysdm.cpl
```

3. Откройте вкладку Advanced (Дополнительно)

4. Нажмите Enviroment Variables (Переменные среды)

5. В разделе User variables (Переменные пользователя) выберите переменную Path, нажмите изменить и добавте туда путь к папке, где находится cryptool.exe

6. Если в разделе User variables (Переменные пользователя) нет переменной Path, то нажмите Create (Создать). И создайте переменную Path с нужным путем.

### Проверка установки
```
cryptool --help
```
Если команда выполняется — установка прошла успешно.

## Использование

### Запуск программы:
```
cryptool
```

Вы увидите справочную информацию:
```
Cryptool is a CLI utility that allows you to encrypt and decrypt text files.This application helps you choose the optimal encryption algorithm and provides a user-friendly and intuitive interface for use.

Usage:
  cryptool [command]

Available Commands:
  benchmark   Run an encryption benchmark on a selected file.
  completion  Generate the autocompletion script for the specified shell
  decrypt     Decrypt previously encrypted files and directories.
  encrypt     Encrypt files and directories using the selected encryption algorithm.
  help        Help about any command

Flags:
  -h, --help   help for cryptool

Use "cryptool [command] --help" for more information about a command.
```

### Основные команды

- **encrypt** — шифрование файлов и директорий

- **decrypt** — расшифровка данных

- **benchmark** — тестирование производительности алгоритмов

Подробная информация о конкретной команде:
```
cryptool [command] --help
```

## Лицензия

Этот проект распространяется под лицензией MIT. <br>
Подробности см. в файле [LICENSE](LICENSE).