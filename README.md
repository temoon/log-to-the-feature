# Log to the future

## About
This is a utility for "play" log as it was written: line by line, with delays between entries.

For example, you have a log file `test.log` like this:
```
[2019-04-09T10:17:48.000] Foo ...
[2019-04-09T10:17:48.342] Bar ...
[2019-04-09T10:17:49.012] Baz ...
```
Run `play-log -f '2006-01-02T15:04:05.999' -i 1 test.log` to play it.<br>
It will output line _Foo_ immediately, line _Bar_ just 342 milliseconds later and _Baz_ 670 milliseconds after _Bar_ is printed. 

## Installation and launch
```
git clone https://github.com/temoon/log-to-the-feature.git
cd log-to-the-feature
git submodule init
git submodule update
make install
./bin/play-log --help
```

## Usage of play-log
All examples based on log sections from the beginning of this document.

### Time settings
* `--time-index|-i` - timestamp position on line;
* `--time-format|-f` - format of timestamp:
  * As described at https://golang.org/pkg/time/;
  * `s` (default) - Unix timestamp in seconds (may be with milliseconds after `.`);
  * `ms` - Unix timestamp in milliseconds (last 3 digits will be interpreted as milliseconds);
* `--time-offset|-t` - begin/offset timestamp (must be the same format as log lines).

Print only last line:
```
play-log -f '2006-01-02T15:04:05.999' -i 1 -t '2019-04-09T10:17:48.500' test.log
```

### Input settings
* `--delimiter|-d` - line delimiter (default `\n`);
* `--buffer-size|-b` - buffer size (default 65536).

Read lines delimited by Windows-style and not longer than 10Kb:
```
play-log -d $'\r\n' -b $((10*1024))
```

### Output settings
* `--skip|-v` - number of lines to skip;
* `--limit|-n` - number of lines to print;
* `--speed|-x` - time multiplier (float number from 0 (no delay) to X (X-times slower), default 1).

Skip 10 and print next 20 lines 2x faster:
```
play-log -v 10 -n 20 -x 0.5 test.log
```

## Contacts
* Tema Novikov &lt;novikov.tema@gmail.com&gt;