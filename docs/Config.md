# Configuration

This document covers the configuration parameters available to the server.

## Table of Contents

* [Running](#running)
* [Configuration Parameters](#configuration-parameters)
    + [Logging](#logging)
        - [type](#type)
        - [path](#path)
    + [Server](#server)
        - [host](#host)
        - [port](#port)
        - [shutdown_timeout](#shutdown-timeout)
    + [Database](#database)
        - [type](#type-1)
        - [path](#path-1)
    + [Discord](#discord)
        - [bot_token](#bot-token)
        - [guild_id](#guild-id)
        - [admin_ids](#admin-ids)
        - [registration_channel_id](#registration-channel-id)
        - [registered_role_id](#registered-role-id)



## Running

At bare minimum you will need a configuration file with the following options (an example `server.conf.example` is
provided in the root directory).

```toml
[server]
host = "localhost"
port = "8080"
```

After this file has been created, you can reference it while running the server (by default it searches for `server.conf`):


## Configuration Parameters

### Logging

Provides configuration settings for the types of logging systems available.

Example:

```toml
[logging]
type = "console"   # supported options at the moment are "disabled", "console" and "file"
path = "tmp/log/"       # only used if type is "file", note this is a directory, not a file
```

#### type

The type of logging system to use.

Available options are:

* `disabled`: No logging will be done.
* `console`: Logs will be printed to the console.
* `file`: Logs will be written to a file.

#### path

Only used if type is `file`. This is the path to the directory where the logs will be written.

### Server

Provides configuration settings for the server.

Example:
 
```toml
[server]
host = "localhost"
port = "8080"
shutdown_timeout = 10 # Seconds to wait before terminating the server
```

#### host

The hostname or IP address to bind the server to.

#### port

The port to bind the server to.

#### shutdown_timeout

The number of seconds to wait before terminating the server.

### Database

Provides configuration settings for the database.

Example:

```toml
[database]
type = "bolt"      # specify bolt as the database type
path = "eden.db"   # path to the database file
```

#### type

The type of database to use.

Available options are:

* `bolt`: Bolt is a file db that is used for storing user data. The current bolt driver uses [storm](https://github.com/asdine/storm)

#### path

The path to the database file.

### Discord

Provides configuration settings for the discord bot.

Example

```toml
[discord]
bot_token = "bot user token goes here"
guild_id  = "main guild id goes here"
admin_ids = ["00000000001", "00000000002", "00000000003"]
registration_channel_id = "channel id for registration updates goes here"
registered_role_id = "role id for registered users goes here"
```

#### bot_token

The bot token for the discord bot.

#### guild_id

The id of the discord guild to use.

#### admin_ids

The ids of the discord users that are admins in a list of strings.

#### registration_channel_id

The id of the discord channel to use for registration updates.

#### registered_role_id

The id of the discord role to use for registered users.