/*
Copyright 2020-2021 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package main

import "sync"
import "net/http"

import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type client struct {
	sync.RWMutex

	conns map[*websocket.Conn]chan []*UnreadArticle

	messages chan []*UnreadArticle // For sending pre-marshaled JSON
}

func (c *client) newBabysitter(l *SessionLogger, conn *websocket.Conn, user string) {
	incoming := make(chan []*UnreadArticle)
	l.I.Println("Creating new conn baby sitter.")

	// Send "hello" packet
	unread := GetUnread(l, user)
	if unread != nil {
		err := conn.WriteJSON(unread)
		if err != nil {
			l.W.Println("Closed connection when trying to send hello packet: ", err)
			conn.Close()
			return
		}
	}

	c.Lock()
	c.conns[conn] = incoming
	c.Unlock()

	for {
		msg := <-incoming
		err := conn.WriteJSON(msg)
		if err != nil {
			l.W.Println("Closed connection when trying to send update packet: ", err)
			conn.Close()
			c.Lock()
			delete(c.conns, conn)
			c.Unlock()
			break
		}
	}
	l.I.Println("Conn baby sitter going away.")
}

func (c *client) Broadcast(unread []*UnreadArticle) {
	c.RLock()
	defer c.RUnlock()

	for _, comm := range c.conns {
		comm <- unread
	}
}

var Feeds = &dispatcher{
	clients: map[string]*client{},
}

type dispatcher struct {
	sync.RWMutex

	// keyed by ID
	clients map[string]*client
}

func (d *dispatcher) Upgrade(l *SessionLogger, w http.ResponseWriter, r *http.Request, user string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.E.Printf("Error upgrading websocket, error: %v\n", err)
		return
	}

	// Get or create client for ID
	d.Lock()
	c, ok := d.clients[user]
	if !ok {
		c = &client{
			conns:    make(map[*websocket.Conn]chan []*UnreadArticle),
			messages: make(chan []*UnreadArticle),
		}
		d.clients[user] = c
	}
	d.Unlock()

	c.newBabysitter(l, conn, user)
}

// The feed updater keeps track of what users have new content, this then sends payloads to those users.
func (d *dispatcher) BroadcastLatest(l *SessionLogger, updated map[string]bool) {
	d.RLock()
	defer d.RUnlock()

	l.I.Printf("Update broadcast initiated.\n")
	for user, c := range d.clients {
		if !updated[user] {
			l.I.Printf("No updates for user %v.\n", user)
			continue
		}
		l.I.Printf("Broadcasting updates to user %v.\n", user)

		unread := GetUnread(l, user)
		if unread == nil {
			continue
		}

		c.Broadcast(unread)
	}
}

func (d *dispatcher) BroadcastTo(l *SessionLogger, user string) {
	d.RLock()
	defer d.RUnlock()

	c, ok := d.clients[user]
	if !ok {
		return
	}

	l.I.Printf("Broadcasting special update to user %v.\n", user)

	unread := GetUnread(l, user)
	if unread == nil {
		return
	}

	c.Broadcast(unread)
}
