# ics
This is a simple module to generate iCalendar (ics) files. The API is subject
to change. Currently supports `VCALENDAR` and `VEVENT`.

## Example program
```go
import (
    "github.com/anders/ics"
)

func main() {
    cal := ics.NewCalendar()
    cal.Add(ics.Event{
        "DTSTART": time.Now(),
        "DTEND": time.Now().Add(45*time.Minute),
        "SUMMARY": "Hello World",
    })
    cal.Encode(os.Stdout)
}
```

output:

```
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//github.com/anders/ics
CALSCAL:GREGORIAN
BEGIN:VEVENT
DTEND:20200127T191212Z
DTSTART:20200127T182712Z
SUMMARY:Hello World
END:VEVENT
END:VCALENDAR
```
