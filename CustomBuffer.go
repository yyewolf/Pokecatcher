package main

import "io"

type buf []byte

func (b *buf) Read(p []byte) (n int, err error) {
    n = copy(p, *b)
    *b = (*b)[n:]
    if n == 0 {
        return 0, io.EOF
    }
    return n, nil
}

func (b *buf) Write(p []byte) (n int, err error) {
    *b = append(*b, p...)
    return len(p), nil
}

func (b *buf) String() string { return string(*b) }