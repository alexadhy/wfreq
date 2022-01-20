# Wfreq

Count top 10 most used words from input file / string

There are 2 projects / binary here, one for running the service and the other is a frontend to the service (in the form of cli)

## Quick Start

- Build it first

```bash
$ # build the service and cli
$ make build
$ # run the service
$ ./wfreq-svc
$ # run the wfreq-cli with crimes_and_punishment.txt from Project Guthenberg
$ # use the `-i` flag to input string / filename, and `-wl` option to specify minimum word length 
$ ./wfreq-cli -i ./tests/crimes_and_punishment.txt -wl 3
```

## Details

The project contains 2 binary in monorepo fashion, one is the service, the other is a client in the form of cli app

### wfreq-Svc

Contains 2 endpoints:

``` 
/       => accepts json input, ex: curl -X POST -H 'Content-Type: application/json' -d '{"min_word_length": 3, "content": "re six cucumber abuse cucumber yeah dude whatever mate yeah dude whatever mate"}' "http://localhost:3334/" 
/upload => accepts multipart form file, with text/plain mimetype and utf-8 charset, ex: curl -X POST -F "file=@./tests/crimes_and_punishment.txt" "http://localhost:3334/upload?min_word=3"
```

### wfreq-cli

Usage:

```bash
Usage of ./wfreq-cli:
  -i string
        input string or input file in plain text utf-8 format (default "some input string")
  -u string
        wfreq-svc url (default to http://localhost:3334) (default "http://localhost:3334/")
  -wl int
        minimum word length to be counted (default 2)
```

Example:

- Count top-10 most used words in crimes and punishment by Fyodor Dostoevsky, with 4 as the minimum word length

```bash
$ ./wfreq -i ./tests/crimes_and_punishment.txt -wl 4
$ # will output:
[
    {
        "word": "that",
        "occurence": 2529
    },
    {
        "word": "with",
        "occurence": 1698
    },
    {
        "word": "have",
        "occurence": 1092
    },
    {
        "word": "from",
        "occurence": 709
    },
    {
        "word": "what",
        "occurence": 684
    },
    {
        "word": "were",
        "occurence": 678
    },
    {
        "word": "your",
        "occurence": 631
    },
    {
        "word": "been",
        "occurence": 560
    },
    {
        "word": "they",
        "occurence": 548
    }
]
```


## Stats

```bash
$ bombardier -H "Content-Type: multipart/form-data; boundary=Asrf456BGe4h" -c50 -m POST -f ./workspace/mangtas/wfreq/tests/crimes_and_punishment.txt http://localhost:3334/upload
Bombarding http://localhost:3334/upload for 10s using 50 connection(s)
[=============================================================================================================================================================================================] 10s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec       907.84     734.03    3296.93
  Latency       55.49ms    25.53ms   154.77ms
  HTTP codes:
    1xx - 0, 2xx - 0, 3xx - 0, 4xx - 8947, 5xx - 0
    others - 0
  Throughput:     0.98GB/s```
