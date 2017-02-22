# fbcount

Simple tool that returns the byte count for a given file, provided the line and character position.

This is a good tool to use with go's guru command.


## To install:

`go get github.com/gdey/fbcount/cmd/fbcount`

## To run:

`fbcount $filename.go#$line:$col` 

This is will return a line that can be used with go guru.

i.e.:

`guru describe $(fbcount main.go#80:25)`

