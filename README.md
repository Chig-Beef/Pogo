# POGO
## The Python transpiler to Golang
I am working on a transpiler, and am documenting my process for educational reasons.

## Keep track of the blog posts
To read what I'm doing on this project, and to ask questions.
https://dev.to/chigbeef_77/compiling-python-to-go-pogo-pt1-3lah

## How to run
To run Pogo it's a good idea to have Go installed.
Once you've done that, use Go to compile Pogo.
You can do this using `go build` while in the src directory.
To run Pogo you need to give it a file to compile.
An example of running Pogo would be `Pogo test.py`.
It should write "test.go" to a folder called "Output" (it doesn't it ends up in TypingSystem, my).