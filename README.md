# libp2p-chat-app
Chat app demo based on @whyrusleeping's gist.

# Add/change in tutorial

 * Do this in setup:  `go get -u github.com/ipfs/go-ipfs-addr`
 * In the rare case you have rewritten your go-libp2p repo's paths to gx paths, run this command: `cd $GOPATH/src/github.com/libp2p/go-libp2p && git reset --hard && git pull origin master`
 * Comment out this line in imports:  `"github.com/libp2p/go-libp2p-host"` (not used)
 * In the following lines, replace `:=` with `=`:

        c, _ := cid.NewPrefixV1(cid.Raw, multihash.SHA2_256).Sum([]byte("libp2p-demo-chat"))

        tctx, _ := context.WithTimeout(ctx, time.Second*10)

 * Set `TopicName` variable by adding this as the first line of the `main()` function:

        TopicName := "libp2p-demo-chat"


# Instructions

Create `demo.go`.  Follow @why's tutorial and put all code in `demo.go`.

Build:  `go build -o libp2p-demo main.go`

Run:  `./libp2p-demo`
