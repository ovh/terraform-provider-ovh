# rfc3339
rfc3339 is a library that allows to manipulate durations as rfc3339

## Limitations

This library handles RFC3339 up to the weeks duration. Since the months and years notions are not clear enough to me
at the moment I write the library, I rather not implement it and adds the week than write an invalid code.

## Examples

```golang
    package main

    import (
        "time"
        "fmt"

        "github.com/ybriffa/rfc3339"
    )

    func main() {
        s := rfc3339.FormatDuration(time.Hour)
        fmt.Println("My duration is ", s)
    }
```

```golang
    package main

    import (
        "time"
        "fmt"

        "github.com/ybriffa/rfc3339"
    )

    func main() {
        d, err := rfc3339.ParseDuration("PT42M")
        if err != nil {
            fmt.Println("I got an error: ", err)
            return
        }
        fmt.Println("My duration is ", d)
    }
```