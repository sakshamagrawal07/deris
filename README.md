# Deris - A Redis Clone in Go

Deris is a lightweight in-memory key-value store written in Go, inspired by Redis. It supports basic data operations, expiration, persistence, and pub/sub functionality.

## Features
- **Key-Value Store**: Supports setting and retrieving key-value pairs.
- **Expiration**: Keys can have a TTL (time-to-live) and expire automatically.
- **Persistence**: Periodic backups using Append-Only File (AOF) and JSON snapshots.
- **Pub/Sub**: Publish-subscribe messaging system.
- **Command Queuing**: Ensures sequential execution of commands.
- **Multi-Client Support**: Handles multiple clients on a single thread with the help of an event loop.

## Installation
```sh
git clone https://github.com/yourusername/deris.git
cd deris
go build -o deris
```

## Usage
```sh
./deris
```

## Available Commands
| Command           | Description                                      | Example                   |
|-------------------|--------------------------------------------------|---------------------      |
| `set key value`   | Stores a key-value pair in memory                | `set name John`           | 
| `get key`         | Retrieves the value of a key                     | `get name`                |
| `lpush key value` | Pushes a value to a key from left                | `lpush name John`         | 
| `rpush key value` | Pushes a value to a key from right               | `rpush name John`         | 
| `lpop key`        | Pops a value from a key from left                | `lpop name`               | 
| `rpop key`        | Pops a value from a key from right               | `rpop name`               | 
| `delete key`      | Deletes a key from memory                        | `delete name`             |
| `expire key sec`  | Sets an expiration time (TTL) for a key          | `expire name 10`          |
| `publish key msg` | Publishes a message to a channel                 | `publish news "Hello!"`   |
| `subscribe key`   | Subscribes to a channel for messages             | `subscribe news`          |
| `unsubscribe key` | Unsubscribes to a channel                        | `unsubscribe news`        |
| `bgsave`          | Triggers a background save to a JSON file        | `bgsave`                  |
| `multi`           | Input multiple commands one by one for execution | `multi`                   |
| `execute`         | Execute the multiple queued commands one by one  | `execute`                 |
| `discard`         | Discard the multiple queued commands             | `discard`                 |

## Flags
| Flag                      | Description                                           | Default   |
|---------------------------|-------------------------------------------------------|-----------|
| `--port`                  | Host for the deris server                             | `0.0.0.0` |
| `--host`                  | Sets the server port                                  | `7379`    |
| `--apend-only`            | Starts the backup using AOF                           | `false`   |
| `--clear-aof`             | Clear the AOF file and start a fresh server           | `true`    |
| `--expire-key-cron-timer` | Interval at which the expire key cron job will run    | `10`      | 

## Example Usage
```sh
# Start the server
./deris --port 8080

# Set a key-value pair
set name Alice

# Retrieve the key
get name

# Set expiration
expire name 10

# Subscribe to a channel
sub news

# Publish a message
pub news "Breaking News!"
```

## Contributing
1. Fork the repository
2. Create a new branch (`git checkout -b feature-branch`)
3. Commit your changes (`git commit -m 'Add new feature'`)
4. Push to the branch (`git push origin feature-branch`)
5. Open a Pull Request

## License
MIT License

`