# dupless

`dupless` checks for patterns in function and package names that we want forbid.


Default forbidden package names are:
* `^util[s]$`
* `^helper[s]$`
* `^base$`
* `^interfaces`

All patterns are regular expressions so you have a huge flexibility.

Default forbidden function names is empty (for now).

You can configure it using `functionNames`, `variableNames` and `packageNames` parameters. Examples:

```
dupless -packageNames 'Manager$' -packageNames '^util$' # no packages with suffix `Manager` or that it's name is exactly `util`.
dupless -functionNames 'BadWord' -functionNames '^foo' # no functions that contain `BadWord` in any place and no functions with `foo` prefix 
dupless -variableNames '^xyz$' # no variable names that it is exactly `xyz`
```

As you can see, you can define the parameter multiple times and all patterns will apply.

## Installation

```
go get github.com/bkielbasa/dupless/cmd/dupless
```
