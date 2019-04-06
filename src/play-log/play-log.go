package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/spf13/pflag"
)

const UnixFormat = "s"
const UnixMsFormat = "ms"

const TimeFormatHelp = `
Time format (https://golang.org/pkg/time/) and:
  s    Seconds since epoch
  ms   Milliseconds since epoch
`

type Opts struct {
    TimeIndex     uint
    TimeFormat    string
    TimeOffset    string
    LineDelimiter string
    BufferSize    uint
    SkipLines     uint
    LimitLines    uint
    SpeedModifier float64
}

func main() {
    var err error

    filename, opts := parseArgs()
    file := os.Stdin

    if filename != "" {
        if file, err = os.Open(filename); err != nil {
            fatal("Log open error:", err)
        }
    }

    scanner := bufio.NewScanner(file)
    scanner.Split(makeSplitFunc(opts))
    scanner.Buffer([]byte{}, int(opts.BufferSize))

    var data []byte
    var ts time.Time
    var prevTs time.Time
    var delay time.Duration

    if opts.TimeOffset != "" {
        if prevTs, err = parseTimestamp([]byte(opts.TimeOffset), opts.TimeFormat, 0); err != nil {
            fatal("Time offset parse error:", err)
        }
    }

    offset := opts.SkipLines
    count := uint(0)

    for scanner.Scan() {
        if offset > 0 {
            offset--

            continue
        }

        data = scanner.Bytes()

        if ts, err = parseTimestamp(data, opts.TimeFormat, int(opts.TimeIndex)); err != nil {
            fatal(err)
        }

        if !prevTs.IsZero() {
            if ts.Equal(prevTs) || ts.After(prevTs) {
                delay = time.Duration(float64(ts.Sub(prevTs)) * opts.SpeedModifier)

                time.Sleep(delay)
            } else {
                warning("Next timestamp is lesser:", prevTs, "<", ts)
            }
        }
        time.Now()

        prevTs = ts

        if _, err = os.Stdout.Write(data); err != nil {
            fatal("Data print error:", err)
        }

        count++

        if opts.LimitLines != 0 && opts.LimitLines == count {
            break
        }
    }

    if err = scanner.Err(); err != nil {
        fatal("Log read error:", err)
    }
}

func parseArgs() (string, *Opts) {
    opts := &Opts{}

    flags := pflag.NewFlagSet("play-log", pflag.ContinueOnError)
    flags.SortFlags = false

    // Time settings
    flags.UintVarP(&opts.TimeIndex, "time-index", "i", 0, "Time token offset")
    flags.StringVarP(&opts.TimeFormat, "time-format", "f", UnixFormat, "Time format")
    flags.StringVarP(&opts.TimeOffset, "time-offset", "t", "", "Time offset")
    // Input settings
    flags.StringVarP(&opts.LineDelimiter, "delimiter", "d", "\n", "Line delimiter")
    flags.UintVarP(&opts.BufferSize, "buffer-size", "b", bufio.MaxScanTokenSize, "Buffer size")
    // Output settings
    flags.UintVarP(&opts.SkipLines, "skip", "v", 0, "Skip lines from output")
    flags.UintVarP(&opts.LimitLines, "limit", "n", 0, "Limit lines to output")
    flags.Float64VarP(&opts.SpeedModifier, "speed", "x", 1, "Speed modifier")

    if err := flags.Parse(os.Args[1:]); err != nil {
        if err == pflag.ErrHelp {
            fmt.Print(TimeFormatHelp)

            os.Exit(0)
        }

        fatal("Arguments parse error:", err)
    }

    return flags.Arg(0), opts
}

func makeSplitFunc(opts *Opts) bufio.SplitFunc {
    delimiter := []byte(opts.LineDelimiter)
    offset := 0

    return func(data []byte, atEOF bool) (int, []byte, error) {
        if atEOF && len(data) == 0 {
            return 0, nil, nil
        }

        if i := bytes.Index(data, delimiter); i >= 0 {
            offset = i + len(delimiter)

            return offset, data[0:offset], nil
        }

        if atEOF {
            return len(data), data, nil
        }

        return 0, nil, nil
    }
}

func parseTimestamp(data []byte, format string, index int) (ts time.Time, err error) {
    switch format {
    case UnixFormat, UnixMsFormat:
        var sec uint64
        var secString = make([]byte, 0, 10)
        var nsec uint64
        var nsecString = make([]byte, 0, 3)

        var i int

        for i = index; i < len(data); i++ {
            if data[i] < '0' || data[i] > '9' {
                break
            }

            secString = append(secString, data[i])
        }

        if sec, err = strconv.ParseUint(string(secString), 10, 64); err != nil {
            return time.Time{}, err
        }

        if i < len(data) && data[i] == '.' {
            i++
        }

        for ; i < len(data) && len(nsecString) <= 3; i++ {
            if data[i] < '0' || data[i] > '9' {
                break
            }

            nsecString = append(nsecString, data[i])
        }

        if len(nsecString) != 0 {
            if nsec, err = strconv.ParseUint(string(nsecString), 10, 64); err != nil {
                return time.Time{}, err
            }
        }

        if format == UnixMsFormat {
            nsec = sec
            sec = 0
        }

        return time.Unix(int64(sec), int64(nsec*uint64(time.Millisecond))), nil
    default:
        return time.Parse(format, string(data[index:index+len(format)]))
    }
}

func fatal(message ...interface{}) {
    warning(message...)
    os.Exit(2)
}

func warning(message ...interface{}) {
    if _, err := fmt.Fprintln(os.Stderr, message...); err != nil {
        panic(err)
    }
}
