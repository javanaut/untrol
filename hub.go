// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "bytes"
    "log"
)

const (
	// Max number of clients allowed in a room
	maxRoomSize = 100
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	rooms map[string][]*Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
        rooms:      make(map[string][]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
            if _, ok := h.rooms[client.room]; !ok {
                // add new room
                h.rooms[client.room] = make([]*Client, maxRoomSize)
            }
            // if len(h.rooms[client.room]) <= maxRoomSize
            //    add client to appropriate room 
            // else report error
            
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
            fields  := bytes.SplitN(message, []byte(":"), 2)
            msgRoom := string(fields[0])
            msg     := fields[1]
			for client := range h.clients {
                // only broadcast message to clients in same room
                if msgRoom == client.room {
                    log.Printf("sending message %s", msg)
                    select {
                    case client.send <- msg:
                    default:
                        close(client.send)
                        delete(h.clients, client)
                    }
                } else {
                    log.Printf("skipping client.room %s as message is for room %s", client.room, msgRoom)
                    log.Printf("original message is %s", message)
                }
            }
		}
	}
}
