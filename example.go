package main

import "github.com/josa42/go-terminal-markdown/markdown"

func main() {
	markdown.MaxImageWidth = 400
	markdown.Print(`
# Hallo Welt
##Foobar
## Foobar2

> Blog quote
> here

Hallo Welt
==========

Foobar 2
--------

-----




` + "```" + `
function run() {
	console.log("Foo"); console.log("Foo"); console.log("Foo"); console.log("Foo"); console.log("Foo"); console.log("Foo"); console.log("Foo")
}
` + "```" + `

` + "```" + `
function run() {
	console.log("Foo")
}
function run() {
	console.log("Foo")
}
function run() {
	console.log("Foo")
}
function run() {
	console.log("Foo")
}
function run() {
	console.log("Foo")
}
function run() {
	console.log("Foo")
}
function run() {
	console.log("Foo")
}
` + "```" + `


<!--

Comment here

-->

Some text with **bold** parts.
Some text with __bold__ parts.

<!--

And here

-->


####### Foobar3

[Google](http://google.de)

![](https://www.huement.com/web/wp-content/uploads/2013/10/logo-1.jpg)
	`)
}
