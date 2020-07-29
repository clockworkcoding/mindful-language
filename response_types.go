package main

const (
	noResponse     int = -1
	threadResponse int = iota + 1
	channelResponse
	ephemeralResponse
	directMessageResponse
)
