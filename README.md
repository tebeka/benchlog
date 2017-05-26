# benchlog - Keep Log of Benchmarks

When optimizing you do a lot of trail and error. And you need to keep track of
performance as you do that. `benchlog` comes to help here, it's intended to run
with `go test -bench` and will keep log of how the benchmarks did.


## Install

    go get github.com/tebeka/benchlog

## Usage

Then when running benchmarks do

     go test -exec benchlog -v -run XXX -bench .

Logs will be save in `.bench.log` (markdown format). You can change the log
file by setting `BENCHLOG_FILE` environment variable.
