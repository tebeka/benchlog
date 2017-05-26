# benchlog - Keep Log of Benchmarks

When optimizing you do a lot of trail and error. And you need to keep track of
performance as you do that. `benchlog` comes to help here, it's intended to run
with `go test -bench` and will keep log of how the benchmarks did.


## Install

    go get github.com/tebeka/benchlog

## Usage

Then when running benchmarks do

     go test -exec benchlog -v -run XXX -bench .

Logs will be save in `.benchlog` (markdown format) as well as some metadata and
git diff. You can change the log directory by setting `BENCHLOG_DIR`
environment variable.
