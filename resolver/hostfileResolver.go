package resolver

import (
	"log"
	"net"
	"sync"

	"github.com/lextoumbourou/goodhosts"
)

type hostfileResolver struct {
	mutex    sync.RWMutex
	hosts    map[string]*hostsEntry
	hostfile goodhosts.Hosts
	stopped  chan struct{}
}

// NewHostFileResolver is a factory method that returns a new instance of the hostfileResolver structure
func NewHostFileResolver() (*hostfileResolver, error) {
	hostfile, err := goodhosts.NewHosts()
	if err != nil {
		return nil, err
	}
	return &hostfileResolver{
		hosts:    make(map[string]*hostsEntry),
		hostfile: hostfile,
		stopped:  make(chan struct{}),
	}, nil
}

func (r *hostfileResolver) AddHost(id string, addr net.IP, name string, aliases ...string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	hosts := append([]string{name}, aliases...)
	r.hosts[id] = &hostsEntry{Address: addr, Names: hosts}
	err := r.hostfile.Add(addr.String(), hosts...)
	if err != nil {
		return err
	}
	return r.hostfile.Flush()
}

func (r *hostfileResolver) RemoveHost(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	hostEntry, ok := r.hosts[id]
	if ok {
		delete(r.hosts, id)
		err := r.hostfile.Remove(hostEntry.Address.String(), hostEntry.Names...)
		if err != nil {
			return err
		}
		return r.hostfile.Flush()
	}
	return nil
}

func (r *hostfileResolver) AddUpstream(id string, addr net.IP, port int, domains ...string) error {
	log.Printf("AddUpstream is not supported by HostFileResolver")
	return nil
}

func (r *hostfileResolver) RemoveUpstream(id string) error {
	log.Printf("RemoveUpstream is not supported by HostFileResolver")
	return nil
}

func (r *hostfileResolver) Listen() error {
	log.Printf("Listen is not supported by HostFileResolver")
	return nil
}

func (r *hostfileResolver) Close() {
	log.Printf("Close is not supported by HostFileResolver")
}

func (r *hostfileResolver) Wait() error {
	<-r.stopped
	return nil
}
