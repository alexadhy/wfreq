# Wfreq

Count top most used words from input file / string

## Usage

```
wfreq -i < -min [min_word_length] >  [input_file | input_string]
```

## Build

To build, run:

```
make build
```

## Stats

```bash
$ time $(wfreq -i ./tests/crimes_and_punishment.txt)

________________________________________________________
Executed in   63.41 millis    fish           external
   usr time   59.70 millis    0.10 millis   59.60 millis
   sys time    8.40 millis    1.61 millis    6.79 millis

```



