# Go Bard
This is a Bard reverse engineering project, made to let me familiarize with golang.

## Run it
Get help
```sh
$ go run main.go --help

[!] Usage: main <info / chat> [options]
Usage of info:
  -session1 string
        '__Secure-1PSID' cookie from google
  -session2 string
        '__Secure-1PSIDTS' cookie from google
Usage of chat:
  -session1 string
        '__Secure-1PSID' cookie from google
  -session2 string
        '__Secure-1PSIDTS' cookie from google
  -text string
        your message to bard (default "こんにちは")
exit status 1
```

```sh
go run main.go chat -session1 "<session_string>" -session2 "<session_string>" -text "<msg_to_Bard>"
```

## Caution
It is inappropriate to use this command util extensively to interact with bard, 
since this is a reverse engineering project and a project for learning purposes.
