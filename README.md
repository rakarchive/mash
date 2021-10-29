# `mash`
A simple shell written in golang.

# Installation
```bash
# Clone the repository
$ git clone https://github.com/raklaptudirm/mash.git

# Go to the repository directory
$ cd mash

# Build the executable
$ go build
```
Add the built executable to your `PATH`.

# Usage
Start the shell:
```bash
$ mash
```
Run a single command:
```bash
$ mash [ command ]
```

# Roadmap
- [x] Input and output.
- [x] Execute exe files.
- [x] `cd` command.
- [x] `exit` command.
- [x] Catch `ctrl+c` and `SIGTERM`.
- [x] Better argument parsing.
- [ ] Customization through root files.
