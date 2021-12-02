# Escaping Eden

Escaping Eden is a MUD/Roguelike written in Go.

It implements a unique xterm-256 based UI over telnet, and as such, supported terminals are limited.

## Table of Contents

  * [Introduction](#introduction)
  * [Eden](#Eden)
  * [Roadmap](#roadmap)
  * [Development](#development)
  * [Connecting](#connecting)
  * [Building](#building)
  * [Running](#running)
  * [Contributing](#contributing)
  * [License](#license)


## Introduction

Escaping Eden was born out of a desire to create a MUD/Roguelike that was a combination of Dwarf Fortress Adventure Mode
style gaming as well as the classic MUDs of the past. 

I have always wanted the ability to telnet into a game server, and to play through a terminal. However, finding a MUD with
even half of the features I wanted was difficult enough, let alone something in the vein of a Roguelike.

While there have been numerous attempts in the past at creating a MUD with a Roguelike-style interface, I have always found
then lacking in the actual gameplay I was hoping to find.

Thus, Escaping Eden was born.

## Eden 

The goal of Escaping Eden is to have a Roguelike with MUD-like features. That is to say, imagine a Dwarf Fortress-style
Map, with MUD-like combat, crafting, building, and other features.

Permanent death is a fixture of this game, as well as other Roguelike staples, such as hunger, thirst, and exhaustion.
Not all is lost, as your name may be remembered among the living, along with other secrets yet to be divulged.

Building will be a key feature of Escaping Eden, as it will be a way to create a world that is more than just a map.
Players will be able to collect resources, gaining skills along the way, in order to craft items and build towering structures.

While doing it alone will be possible, you'll find that playing with a group of others, perhaps in a settlement or town,
will allow you to achieve more. Don't be fooled by the safety of living in a town, as townspeople are more mouths to feed,
and the wealth of the town may attract bandits and other hostile creatures.

There are many more features to be found in Escaping Eden, and I hope that you'll stick around while I continue to develop
it! 

## Roadmap

The current project is still in its infancy, and almost all gameplay features are missing.

The UI is being stabilized, and the game engine component is being planned. I say planned, because there are 
different methods of communicating with the game engine that I am attempting before settling on one. 

This is important for the long term goals of EE, because I want to be able to be able to distribute the game processing
to a separate cluster in the future if necessary. This is a very long term goal, but it's essential that the game
engine and the networking/UI are completely separate entities. The connection manager can send messages to the game engine
on behalf of a console, but a console can never communicate directly with the game engine.

The game engine can _never_ have any notion of what a console is, nor should it ever care. 

Most of the work ongoing at the moment is to get the overall engine in a state where it can be tested by 
players in a live environment.

## Development

Escaping Eden is currently built on Go 1.17, and relies on xterm-color256 support in the terminal.

It relies heavily on ANSI escape sequences, and will not work on terminals that do not support them. If you are 
planning on contributing to any of the UI work, it is basically mandatory that you review 
the [ANSI escape code documentation](https://en.wikipedia.org/wiki/ANSI_escape_code) before you start.

## Connecting

Connect to the server with telnet, using a terminal with type xterm-256color. No other methods of connecting are
currently supported. There are plans to add other terminal types, and eventually a standalone client, but this is a long
way off.

By default, the development server launches on port 8080.

```bash
$ telnet localhost 8080
```

**Note that the official server runs on port 23 @ world.eden.sh**

### Windows

It may be possible to connect with PuTTY on Windows with the following steps, but results may vary:

`Settings -> Connection > Data > Terminal-type change to xterm-256color`

Assistance with adding putty-256color support would be appreciated!


## Building

Building and running the current version of Escaping Eden is relatively straightforward, however as the project matures
expect the configuration to grow more complex as customization features are added.

```bash
$ go mod download  
$ go build . 
```

## Running

At bare minimum you will need a configuration file with the following options (an example `server.conf.example` is 
provided in the root directory).

```toml
[server]
host = "localhost"
port = "8080"
```

After this file has been created, you can reference it while running the server (by default it searches for `server.conf`):

```bash
$ ./EscapingEden -c example-name.conf 
```

## Contributing

You can contribute to the project by submitting pull requests on [Github](https://github.com/yamamushi/EscapingEden/pulls),
or by joining our Discord community at [Discord](https://discord.gg/uMxZnjJGGu).

Most conversation about the development of Escaping Eden will be on Discord. 

## License

Unless otherwise noted, all content is licensed under the [GPL v3](https://www.gnu.org/licenses/gpl-3.0.html) license.

(c) 2021 Yamamushi
