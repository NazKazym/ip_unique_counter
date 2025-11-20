# Unique IPv4 Counter

This project is a Go program for counting **unique IPv4 addresses** in
very large text files.\
It is designed to handle inputs that are too big to fit into memory
(tens or even hundreds of gigabytes).

The goal of the program is:

> **Read a huge file with one IP address per line and report how many
> unique IPv4 addresses it contains.**

Along the way, it also reports:

-   total lines processed\
-   number of valid IPv4 addresses\
-   number of invalid lines (logged separately)\
-   time spent and processing speed

The program has been tested on multi-billion-line files and runs
entirely in streaming mode.

------------------------------------------------------------------------

## How It Works (short explanation)

-   The file is read line by line using a buffered reader.
-   A lightweight IPv4 parser checks if each line contains a valid IPv4
    address.
-   Valid IPs are mapped into a **512 MB bitmap** (one bit per IPv4).
-   512MB is enough to cover all 2<sup>32</sup> possible values of IPv4. 
-   The bitmap is split across several workers so different parts of the
    bitmap are updated in parallel.
-   After the whole file is processed, the bitmap is scanned to count
    all set bits (i.e., unique IPs).

Invalid lines are written to `errors.log` along with their line number.

------------------------------------------------------------------------

## Running the Program

Build:

``` bash
go build -o ipcounter .
```

Then run:

``` bash
./ipcounter
```

There is no special file format required --- just one IPv4 address per
line.

------------------------------------------------------------------------

## Example Output

Here is a real run on the test file:

    Using 16 parallel workers with single 512 MiB bitmap
    Progress :: 7,980,000,000 lines | 6.56 M/s | 7,980,000,000 IPs | elapsed 20m16s

    Finished!
       Total lines processed : 8,000,000,000
       Valid IPv4 addresses  : 8,000,000,000
       Total time            : 20m18.968s

    Counting unique IPs...
       Average speed         : 6.56 M lines/sec
       Count time            : 23.53ms

    Unique IPv4 addresses : 1000000000
    Time elapsed          : 20m18.996s
    Memory used           : 592.36 MB


------------------------------------------------------------------------
## Input
There is a [config.yaml](config.yaml) file with parameters e.g. source file's path, size of buffer 
## Input Format

Each line should contain a single IPv4 address:

    192.168.0.1
    10.0.0.5
    8.8.8.8

Newline variations (`` or ``) are fine. Extra spaces are ignored.

If the parser finds something wrong in a line, it writes it to
`errors.log`, for example:

    1521 | "999.999.0.1" | octet > 255

------------------------------------------------------------------------

## When to Use This Tool

This project is a good fit if you need to:

-   process very large files where loading everything into memory is
    impossible\
-   quickly count unique IPv4s from logs, datasets, or network captures\
-   experiment with performance, concurrency, and memory-efficient data
    structures in Go

It is not meant to be a general-purpose log parser --- it does exactly
one task, but does it efficiently.

------------------------------------------------------------------------

## Project Structure

    .
    ├── bitmap.go        # Bitmap data structure
    ├── counter.go       # Main processing loop
    ├── parser.go        # IPv4 parser
    ├── progress.go      # Progress reporting
    ├── main.go          # CLI entry point
    └── README.md

------------------------------------------------------------------------

## Notes on Testing

The program was tested with downloaded [file](https://ecwid-vgv-storage.s3.eu-central-1.amazonaws.com/ip_addresses.zip) with size about 106GB.

Max performance I got: **6--7 million lines per second**.

------------------------------------------------------------------------

