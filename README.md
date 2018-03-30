# NBTNav
View NBT files within a terminal window.

## Installation
I'm too lazy to provide download links, but luckily Golang
makes compiling from source really easy.

1. [install Golang and set your GOPATH](https://golang.org/doc/install).
2. (Optional) Add $GOPATH/bin to your system path
3. Run `go get github.com/deanveloper/nbtnav`
4. The executable is now located at $GOPATH/bin/nbtnav. If you skipped step 2, either
move or symlink this file to one of the directories in your PATH

## Usage
Use `nbtnav <file>` to view an nbt file! Easy as that.

## Commands
Once you have run `nbtnav <file>`, you will be given a command prompt. 
The command names in NBTNav are inspired by those found in bash.

* `help`: Lists all commands
* `cd <compound>`: Moves into an NBT Compound
* `ls [compound]`: Lists all elements within the compound you are in,
or the one that you supply
* `tree [compound]`: Similar to `ls`, but does a deep search, showing the entire tree
* `cat <tag>`: Displays the value at a given tag
*  save \[compress\] \[output\] : Saves the current NBT tree to output. `compress`
can be any of `gzip`, `zlib`, or `none`(default). `output` is the output file and
defaults to the original file name.
* `set <tag> <type> [value]`: Sets the tag to the give type and value
* `exit`: Exits NBTNav

## Future Features
* Sort NBT tags in the order they appear in (not a necessary feature, but could possibly be useful somehow)
