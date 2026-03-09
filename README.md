# CYRIL [WIP]
_Currently a work in progress_

The CLI butler you didn't know you needed. Get system-wide access to your notes with a minimalistic TUI :) Create, read, edit and give aliases to your notes so you never have to dig through hundreds of notes trying to find the one you need - trust me, I've been there which is why I built this. No more juggling through open applications and tabs trying to find the note you need, it'll be right where you work - in your terminal.

## Installation
Requires go version >=1.25.4

Just type `go install github.com/yendelevium/cyril@latest` and it will be ready to use (make sure you've exported go to $PATH)


## Working
Your notes will be organized under various `topics`, which is just a fancy way of saying folders. This will let you distinguish between different notes with the same name -> putting them in different topics. In case you forgot where your note is, you can list all your available notes or list them topic wise.

cyril uses your default editor for editing and creating notes, letting you use the editor you're already comfortable with. cyril looks for a config file under `~/.config/cyril/config.yml`. Here's how a sample config file looks like:

```yml
editor: nvim
store: $HOME/Documents/cyrilStore
defaultTopic: general
dbPAth: $HOME/Documents/cyrilStore
```

It falls back to defaults if you don't have the file so no worries.

## Available commands

### `create`
Available commands
- cyril create {filename}
- cyril create {filename} -t {topicname}

This creates a file under the specified topic (default if not passed) -> also lets you add content to it upon creation.

### `read`
Available commands
- cyril read {filename/aliasname}
- cyril read {filename/aliasname} -t {topicname} _[Need to implement]_

Displays the given note. You can access it via the filename or the aliasname (strictly speaking, technically everything is an aliasname, even the original filename). If there are multiple files with the same name across topics/multiple prefix-matches for that file, you can choose which one to read via a selector. 

### `edit`
Available commands
- cyril edit {filename/aliasname}
- cyril edit {filename/aliasname} -t {topicname} _[Need to implement]_

Lets you edit the file mentioned. Will give the option to create a new file if it doesn't exist _[Need to implement]_. Similarly to the read command, multiple files trigger a file selector.

### `alias`
Available commands
- cyril alias {new name} {filename to alias}

Just adds a new alias name for the given file (or an existing aliasname) in the same topic where the file belongs.

### `list`
Available commands
- cyril list
- cyril list -t {topicname}

Lists available files system wide or under a specific topic.

## Future work
There are so many more features I want to add here

- a delete command to delete the file and/or an aliasname
- a way to create repeating notes (probably through cron jobs)
- other QOL features like better user flows and error(edge case?) handling
- moving files between topics
- a prettier TUI (some sort of a side-by-side filename and file content display? idk)
- more customizability, like letting the user choose colour schemes based on the config
- eliminating the need for commands all togther and just invoking stuff via `cyril {filename} -{flags}`
- versioned github releases
- markdown rendering (idk)

and a lot more but these are the few off the top of my head.

## Soliloquy
I made this for me to solve a problem I faced, which is why I'm super proud because this is actually a tool I'm actively using. And since I actually wanted to make this (not for my resume or to attract job recruiters), I had a lot of fun. And this project is **AI free**. No vibe coded slop here. I've written (or copy-pasted from stack overflow/blogs) every line of code in this project. Maybe this sounds stupid to you but in a world where everyone wants to just make "design decisions" and let the AI write code, I want to write code because well, I love to write code. I don't see it as a menial task, I enjoy writing code (especially golang) as much as I like designing systems. And since this is a personal project, there's no deadline I need to meet so I can take my sweet time writing whatever I wanna write :) 

Anyways, give this a try if you're interested. Let me know if you like it, byeeee.