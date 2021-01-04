# GoLox Interpreter

An implementation of jlox, from Crafiting Interpreters, using Golang.

Implemented using Java in the book, the only difference is that I'm implementing it using Golang (while I'm still learning how to write Golang).

[Crafiting Interpreters](https://craftinginterpreters.com)

## Instructions

```
go build or go run lox.go
(TODO)
```
# Example Syntax

file.txt
```
print "## Scope";
var a = "global a";
var b = "global b";
var c = "global c";
{
  var a = "outer a";
  var b = "outer b";
  {
    var a = "inner a";
    print a;
    print b;
    print c;
  }
  print a;
  print b;
  print c;
}
print a;
print b;
print c;
print "## While";
var d = 0;
while (d <= 10) {
	print d;
	d = d + 2;
}
print "## For";
for (var e = 0; e <= 50; e = e + 5) {
  if (e == 10 or e == 20 or e == 30 or e == 50) {
	print e;
  }
}
```

Running:
```
go run lox.go file.txt
```
