package pool

import (
	"telegram-for-cms/telegram_cms/messenger"
)

type ClientPool struct {
	clients map[string]*messenger.Client
}

// NewSessionPool Конструктор
func NewClientPool() *ClientPool {
	return &ClientPool{
		clients: make(map[string]*messenger.Client),
	}
}

func (sp *ClientPool) Add(name string, client *messenger.Client) {
	if _, exists := sp.clients[name]; !exists {
		sp.clients[name] = client
	}
}

func (sp *ClientPool) Remove(name string) {
	delete(sp.clients, name)
}

func (sp *ClientPool) Clear() {
	clear(sp.clients)
}

func (sp *ClientPool) Get(name string) (*messenger.Client, bool) {
	client, exists := sp.clients[name]
	return client, exists
}

func (sp *ClientPool) Len() int {
	return len(sp.clients)
}

func (sp *ClientPool) IsEmpty() bool {
	return sp.Len() == 0
}

func (sp *ClientPool) NonEmpty() bool {
	return sp.Len() > 0
}

func (sp *ClientPool) Snapshot() map[string]*messenger.Client {
	if sp.clients == nil {
		return nil
	}
	cp := make(map[string]*messenger.Client, len(sp.clients))
	for k, v := range sp.clients {
		cp[k] = v
	}
	return cp
}

func (sp *ClientPool) Close(name string) {
	if _, exists := sp.clients[name]; exists {
		sp.clients[name].Close()
		sp.Remove(name)
	}
}

func (sp *ClientPool) CloseAll() {
	names := make([]string, 0, sp.Len())

	for name := range sp.clients {
		names = append(names, name)
	}

	for _, name := range names {
		sp.Close(name)
	}
	sp.Clear()
}
