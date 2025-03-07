# Chat App Over TCP

Chat application that uses TCP protocol.

## Features

- Public message broadcast (by default)
- Private messages
- Clients list
- Command system
- Kicking out by username
- Message of the day
- Contains server and a client in one package
- Usernames are unique (e.g. user#12ef)

## How to build

```
git clone https://github.com/silentstranger5/chat
cd chat
go build .
./chat
```

## How to use

How to use application:

```
./chat -help
```

Connect to address:

```
./chat -address "<chat.address>:<port>" -user <username>
```

Set up a server:

```
./chat -mode server -address ":<port>"
```

Both client and server contain a command system. 
To see a list of commands, type `help` once connected.