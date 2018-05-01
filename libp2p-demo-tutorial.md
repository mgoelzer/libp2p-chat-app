# Building a P2P app with go-libp2p

Libp2p is a peer-to-peer networking library that allows developers to easily add p2p connectivity between their users. Setting it up and using it is simple and easy! Let's build a demo libp2p app that lets us type messages to other nodes running the app, and see messages from other nodes. A very simple chat app.

To start, make sure you have Go installed and set up. Then install libp2p and some other deps we need with:

```shell
go get -u github.com/libp2p/go-libp2p
go get -u github.com/libp2p/go-floodsub
go get -u github.com/libp2p/go-libp2p-kad-dht
go get -u github.com/ipfs/go-ipfs-addr
```

Now for some code, We will start with a few imports. These imports include go-libp2p itself, our pubsub library "floodsub", the IPFS DHT, and a few other helper packages to tie things together. 

```go
package main

import (
    "bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/libp2p/go-floodsub"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs-addr"

	"github.com/multiformats/go-multihash"
)
```

Next up, lets start constructing the pieces! First, we will set up our libp2p host. The host is the main abstraction that users of go-libp2p will deal with, It lets you connect to other peers, open new streams, and register protocol stream handlers.

```go
func main() {
	ctx := context.Background()

	// Set up a libp2p host.
	host, err := libp2p.New(ctx, libp2p.Defaults)
	if err != nil {
		panic(err)
	}
 
    // ... everything else goes here ...
}
```

Next, we set up our libp2p "floodsub" pubsub instance. This is how we will communicate with other users of our app. It gives us a simple many to many communication primitive to play with.

```go
// Construct ourselves a pubsub instance using that libp2p host.
fsub, err := floodsub.NewFloodSub(ctx, host)
if err != nil {
	panic(err)
}
```

And finally, we need a way to discover other peers. Future versions of pubsub will likely have discovery and rendezvous built into the protocol, but for now we have to do it ourselves. We will use the DHT for this since its pretty straightforward, but it does tend to be slow for what we want. Bootstrapping an entire DHT and filling its routing tables take a little while.

```go
// Using a DHT for discovery.
dht := dht.NewDHTClient(ctx, host, datastore.NewMapDatastore())
if err != nil {
	panic(err)
}
```

Now for that discovery I mentioned, we need to do a few things. First, we need to connect to some initial bootstrap peers. We can use some of the IPFS bootstrap peers for this, even though we aren't using IPFS for our app, we are running libp2p and using the same DHT so it all works out.

```go
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
```

Next up the rendezvous. We need a way for users of our app to automatically find eachother. One way of doing this with the DHT is to tell it that you are providing a certain unique value, and then to search for others in the DHT claiming to also be providing that value. This way we can use that value's location in the DHT as a rendezvous point to meet other peers at.

```go
// Using the sha256 of our "topic" as our rendezvous value
TopicName := "libp2p-demo-chat"
c, _ = cid.NewPrefixV1(cid.Raw, multihash.SHA2_256).Sum([]byte(TopicName))

// First, announce ourselves as participating in this topic
fmt.Println("announcing ourselves...")
tctx, _ := context.WithTimeout(ctx, time.Second*10)
if err := dht.Provide(tctx, c, true); err != nil {
	panic(err)
}

// Now, look for others who have announced
fmt.Println("searching for other peers...")
tctx, _ = context.WithTimeout(ctx, time.Second*10)
peers, err = dht.FindProviders(tctx, c)
if err != nil {
	panic(err)
}
fmt.Printf("Found %d peers!\n", len(peers))
```

This process might take a while, As I said earlier, booting up a DHT node from scratch takes a bit of time. Once it completes, we will need to connect to those new peers.

```go
// Now connect to them!
for _, p := range peers {
	if p.ID == host.ID() {
		// No sense connecting to ourselves
		continue
	}

	tctx, _ = context.WithTimeout(ctx, time.Second*5)
	if err := host.Connect(tctx, p); err != nil {
		fmt.Println("failed to connect to peer: ", err)
	}
}

fmt.Println("bootstrapping and discovery complete!")
```

At this point, and with any luck, we should be all connected up to other users of our app. Now to do something with it. Let's subscribe to a pubsub channel, listen for messages on it, and then send anything the user types to stdin out as a message.

```go
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
```

And with that, you have a simple chat app! Build it with:

```shell
go build -o libp2p-demo main.go
```

And then run it:

```shell
./libp2p-demo
```

