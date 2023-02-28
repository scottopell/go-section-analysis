Go Elf Section Analysis
-------------------------

Based heavily on https://stackoverflow.com/a/70777803
Basically taken as-is and added some formatting options.

Currently hard-coded to look at `.rodata`, search in `main.go` to change

Usage:
```
$ go run main.go ../datadog-agent/bin/dogstatsd/dogstatsd | head -n 20
github.com/stretchr/testify/assert:
               .text:     1456 bytes (  1.5 kB)
          .noptrdata:      192 bytes (   192 B)
                .bss:      240 bytes (   240 B)
             .rodata:        4 bytes (     4 B)

github.com/gogo/protobuf/jsonpb:
               .text:      448 bytes (   448 B)
          .noptrdata:      128 bytes (   128 B)
                .bss:       40 bytes (    40 B)
             .rodata:        8 bytes (     8 B)

github.com/containerd/ttrpc:
               .text:     5504 bytes (  5.5 kB)
          .noptrdata:      686 bytes (   686 B)
                .bss:       40 bytes (    40 B)
               .data:      176 bytes (   176 B)
           .noptrbss:       12 bytes (    12 B)
             .rodata:        8 bytes (     8 B)
```
