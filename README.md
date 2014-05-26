# Gomilk

Gomilk is a faster Gmilk written by Go.
Gmilk is the command line tool for [Milkode](https://github.com/ongaeshi/milkode).

Search the source code of **40000 files in 0.5 seconds**.

## Installation

Download from here or "go build".

- [gomilk-0.1.0-darwin-amd64.zip](https://dl.dropboxusercontent.com/u/28734381/gomilk-0.1.0-darwin-amd64.zip)
- [gomilk-0.1.0-windows-amd64.zip](https://dl.dropboxusercontent.com/u/28734381/gomilk-0.1.0-windows-amd64.zip)

Put the "gomilk" binary into PATH directory.

And, install [Milkode](https://github.com/ongaeshi/milkode#installation).

## How to use

### 1. Create Milkode database and Add packages

[Milkode - Usage](https://github.com/ongaeshi/milkode#usage)

### 2. Start the Milkode web server with gomilk mode

```
$ milk web -g
```

### 3. Search

The basic usage same as [Gmilk](https://github.com/ongaeshi/milkode#search-command-line).

```
$ gomilk search_keyword
.
lib/a.txt:1: test aaa
test/b.txt:1: test bbb
```

## Performance Test

OSX 10.7.5, Core2 Duo 3.06 GHz, 8 GB RAM.

Test on Ruby-2.1.2. (0.281 seconds 4722 files)

```
$ cd ./ruby-2.1.2

# 4722 files
$ find . | wc
    4722    4765  148103

# Search in 0.281 seconds 4722 files 
$ time gomilk "performance test"
test/bigdecimal/test_bigdecimal.rb:1501:    # this is mainly a performance test (should be very fast, not the 0.3 s)

real	0m0.281s
user	0m0.062s
sys	0m0.024s
```

Test on linux-3.10-rc4. (0.584 seconds 45752 files)

```
$ cd ./linux-3.10-rc4

# 45752 files
$ find . | wc
   45752   45752 1552272

# Search in 0.584 seconds 45752 files 
$ time gomilk "performance test"
arch/ia64/mm/init.c:582: * useful for performance testing, but conceivably could also come in handy for debugging
Documentation/fb/udlfb.txt:155:			performance tests to start and finish in a very short
Documentation/networking/eql.txt:339:  Although you may already done this performance testing, here
Documentation/scsi/FlashPoint.txt:149:embedded system.  I am presently working on some performance testing and
Documentation/serial/sx.txt:202:per second to transmit.  If you do any performance testing in this
.
.
real	0m0.584s
user	0m0.427s
sys	0m0.091s
```

Search all packages. (0.638 seconds 50474 files)

```
$ cd ~

# Search in 0.638 seconds all files
$ gomilk "performance test" -a
ruby-2.1.2/test/bigdecimal/test_bigdecimal.rb:1501:    # this is mainly a performance test (should be very fast, not the 0.3 s)
linux-3.10-rc4/arch/ia64/mm/init.c:582: * useful for performance testing, but conceivably could also come in handy for debugging
.
.
real	0m0.638s
user	0m0.492s
sys	0m0.109s
```

