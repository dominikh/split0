# split0

split0 splits stdin into chunks of N bytes, passing each chunk to a
new invocation of a command.

The command will be invoked as `/bin/sh -c <command>`, and an
environment variable `$N` will be set to the current iteration number.

Example:

```
$ head -c 12 /dev/zero | split0 5 'echo -n "$N: "; wc -c'
0: 5
1: 5
2: 2
```

# Install

    go get honnef.co/go/split0
