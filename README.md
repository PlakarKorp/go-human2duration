# human2duration

`human2duration` is a Go package that parses human-readable time expressions into `time.Duration`.

It supports natural phrases (e.g. `"half an hour"`, `"2d 3h"`), relative time formats (e.g. `"in 5 minutes"`, `"3 days ago"`), and timestamps.

## Features

* Parses durations like `"2h30m"`, `"3 days"`, `"1.5 hours"`, `"2d 3h"`, etc...
* Parses golang time.Duration
* Supports fuzzy expressions like `"half a day"`, `"couple of hours"`, `"an hour and a half"`.
* Parses relative time expressions:

  * `"3 days ago"` → negative duration
  * `"in 2 hours"` or `"after 5 minutes"` → positive duration
* Supports common timestamp formats like:

  * RFC3339 (`2024-01-02T15:04:05Z`)
  * `2006-01-02 15:04:05`
  * `2006-01-02`

## Installation

```bash
go get github.com/PlakarKorp/human2duration
```

## Usage

```go
import (
    "fmt"
    "time"

    "github.com/PlakarKorp/go-human2duration"
)

func main() {
    d, err := human2duration.ParseDuration("2d 3h")
    if err != nil {
        panic(err)
    }
    fmt.Println(d) // 51h0m0s

    ago, _ := human2duration.ParseSinceDuration("5 minutes ago")
    fmt.Println(ago) // -5m0s

    after, _ := human2duration.ParseAfterDuration("in 2 hours")
    fmt.Println(after) // 2h0m0s
}
```

## Supported Units

| Unit   | Aliases                           |
| ------ | --------------------------------- |
| Second | `s`, `sec`, `second`, `seconds`   |
| Minute | `m`, `min`, `minute`, `minutes`   |
| Hour   | `h`, `hr`, `hour`, `hours`        |
| Day    | `d`, `day`, `days`                |
| Week   | `w`, `week`, `weeks`              |
| Month  | `mo`, `month`, `months` (30 days) |
| Year   | `y`, `year`, `years` (365 days)   |

## Fuzzy Phrases

The parser understands the following phrases:

* `half an hour`
* `an hour and a half`
* `half a day`
* `couple of minutes`
* `couple of hours`
* `couple of days`
* `an hour`, `a minute`, `a second`, `a day`, `a week`, `a month`

## License

ISC
