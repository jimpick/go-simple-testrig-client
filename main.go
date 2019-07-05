package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-http"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// newHost illustrates how to build a libp2p host with secio using
// a randomly generated key-pair
func newHost(listen multiaddr.Multiaddr) host.Host {
	h, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(listen),
	)
	if err != nil {
		panic(err)
	}
	return h
}

func main() {
	if len(os.Args) < 2 {
		panic("Need multiaddr")
	}
	addr, err := multiaddr.NewMultiaddr(os.Args[1])
	if err != nil {
		panic(err)
	}
	peer, err := peerstore.InfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connecting to:", os.Args[1])

	m2, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	clientHost := newHost(m2)
	defer clientHost.Close()

	clientHost.Peerstore().AddAddrs(peer.ID, peer.Addrs, peerstore.PermanentAddrTTL)

	// print the node's PeerInfo in multiaddr format
	peerInfo := &peerstore.PeerInfo{
		ID:    clientHost.ID(),
		Addrs: clientHost.Addrs(),
	}
	addrs, err := peerstore.InfoToP2pAddrs(peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address:", addrs[0])

	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))
	client := &http.Client{Transport: tr}
	res, err := client.Get("libp2p://" + peer.ID.String())
	if err != nil {
		panic(err)
	}
	text, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("res text: %s\n", text)
}
