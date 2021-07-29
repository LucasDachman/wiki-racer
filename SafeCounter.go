package main

type SafeCounter struct {
	num  int
	ping chan byte
}

func NewSafeCounter(channelLen int) SafeCounter {
	return SafeCounter{
		ping: make(chan byte, channelLen),
	}
}

func (counter *SafeCounter) Add() {
	counter.ping <- 1
}

func (counter *SafeCounter) Start() {
	go (func() {
		for range counter.ping {
			counter.num++
		}
	})()
}

func (counter *SafeCounter) Stop() {
	close(counter.ping)
}
