package main

import (
	"io"
	"log"
	"net"
	"sync"
)

type task struct {
	conn net.Conn
	epoller *epoll
}

type pool struct {
	workers   int
	maxTasks  int
	//taskQueue chan net.Conn
	taskQueue chan *task

	mu     sync.Mutex
	closed bool
	done   chan struct{}
}

func newPool(w int, t int) *pool {
	return &pool{
		workers:   w,
		maxTasks:  t,
		//taskQueue: make(chan net.Conn, t),
		taskQueue: make(chan *task, t),
		done:      make(chan struct{}),
	}
}

func (p *pool) Close() {
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

func (p *pool) addTask(conn net.Conn, epoller *epoll) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	p.taskQueue <- &task{conn:conn, epoller:epoller}
}

func (p *pool) start() {
	for i := 0; i < p.workers; i++ {
		go p.startWorker()
	}
}

func (p *pool) startWorker() {
	for {
		select {
		case <-p.done:
			return
		case t := <-p.taskQueue:
			if t != nil {
				handleConn(t)
			}
		}
	}
}

func handleConn(t *task) {
	_, err := io.CopyN(t.conn, t.conn, 8)
	if err != nil {
		if err := t.epoller.Remove(t.conn); err != nil {
			log.Printf("failed to remove %v", err)
		}
		t.conn.Close()
	}
	opsRate.Mark(1)
}
