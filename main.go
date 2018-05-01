package main

import (
    "bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/multiformats/go-multihash"

	// TODO:  Would be nice to use github imports instead of gx since
	// a libp2p demo should not rely on any higher level abstractions 
	// like IPFS
	//"github.com/libp2p/go-floodsub"
	//"github.com/libp2p/go-libp2p"
	//"github.com/libp2p/go-libp2p-host"
	//"github.com/libp2p/go-libp2p-kad-dht"
	//"github.com/libp2p/go-libp2p-peerstore"
	//"github.com/ipfs/go-cid"
	//"github.com/ipfs/go-datastore"
	//"github.com/ipfs/go-ipfs-addr"

	// For now we use gx instead:
	"gx/ipfs/QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp/go-libp2p"
	"gx/ipfs/QmSFihvoND3eDaAYRCeLgLPt62yCPgMZs1NSZmKFEtJQQw/go-libp2p-floodsub"
	"gx/ipfs/QmY1y2M1aCcVhy8UuTbZJBvuFbegZm47f9cDAdgxiehQfx/go-libp2p-kad-dht"
	"gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"gx/ipfs/QmXRKBQA4wXP7xWbFiZsR1GP4HV6wMDQ1aWFxZZ4uBcPX9/go-datastore"
	"gx/ipfs/QmQViVWBHbU6HmYjXcdNq7tVASCNgdg64ZGcauuDkLCivW/go-ipfs-addr"
)

func main() {
	TopicName := "libp2p-demo-chat"
	ctx := context.Background()

	// Trap Ctrl+C and cleanup before exit
	sigIntChan := make(chan os.Signal, 1)
	signal.Notify(sigIntChan, os.Interrupt)
	go func(){
		for sig := range sigIntChan {
			// sig is a ^C, handle it
			fmt.Println("Cleaning up... (Ctrl+C again to abort)")

			// unregister the signal handler
			signal.Stop(sigIntChan)

			// cleanup
			// TODO:  not sure what to do here to remove this peer from DHT...

			// terminate
			os.Exit(0)
		}
	}()

	// Set up a libp2p host.
	host, err := libp2p.New(ctx, libp2p.Defaults)
	if err != nil {
		panic(err)
	}

	// Construct ourselves a pubsub instance using that libp2p host.
	fsub, err := floodsub.NewFloodSub(ctx, host)
	if err != nil {
		panic(err)
	}

	// Using a DHT for discovery.
	dht := dht.NewDHTClient(ctx, host, datastore.NewMapDatastore())
	if err != nil {
		panic(err)
	}

	bootstrapPeers := []string{
		"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
		"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
		"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
		"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
		"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
	}

	fmt.Println("bootstrapping...")
	for _, addr := range bootstrapPeers {
		iaddr, _ := ipfsaddr.ParseString(addr)

		pinfo, _ := peerstore.InfoFromP2pAddr(iaddr.Multiaddr())

		if err := host.Connect(ctx, *pinfo); err != nil {
			fmt.Println("bootstrapping to peer failed: ", err)
		}
	}

	// Using the sha256 of our "topic" as our rendezvous value
	c, _ := cid.NewPrefixV1(cid.Raw, multihash.SHA2_256).Sum([]byte("libp2p-demo-chat"))

	// First, announce ourselves as participating in this topic
	fmt.Println("announcing ourselves...")
	tctx, _ := context.WithTimeout(ctx, time.Second*10)
	if err := dht.Provide(tctx, c, true); err != nil {
		panic(err)
	}

	// Now, look for others who have announced
	fmt.Println("searching for other peers...")
	tctx, _ = context.WithTimeout(ctx, time.Second*10)
	peers, err := dht.FindProviders(tctx, c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d peers!\n", len(peers))

	// Now connect to them!
	for _, p := range peers {
		if p.ID == host.ID() {
			// No sense connecting to ourselves
			continue
		}

		tctx, _ := context.WithTimeout(ctx, time.Second*5)
		if err := host.Connect(tctx, p); err != nil {
			fmt.Println("failed to connect to peer: ", err)
		}
	}

	fmt.Println("bootstrapping and discovery complete!")

	sub, err := fsub.Subscribe(TopicName)
	if err != nil {
		panic(err)
	}

	// Go and listen for messages from them, and print them to the screen
	go func() {
		for {
			msg, err := sub.Next(ctx)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%s: %s\n", msg.GetFrom(), string(msg.GetData()))
		}
	}()

	// Now, wait for input from the user, and send that out!
	fmt.Println("Type something and hit enter to send:")
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		if err := fsub.Publish(TopicName, scan.Bytes()); err != nil {
			panic(err)
		}
	}











}

