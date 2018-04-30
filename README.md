# libp2p-chat-app
Chat app demo based on @whyrusleeping's gist.

# Add/change in tutorial

 * Do this in setup:  `go get -u github.com/ipfs/go-ipfs-addr`
 * Replace some `github.com` imports with `gx` equivalents as [described here](https://gist.github.com/whyrusleeping/169a28cffe1aedd4419d80aa62d361aa#gistcomment-2560605)
 * Create an empty `package.json` by running `npm init`
 * Import the `gx` packages you need:

        gx import QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp
        gx import QmSFihvoND3eDaAYRCeLgLPt62yCPgMZs1NSZmKFEtJQQw
        gx import QmY1y2M1aCcVhy8UuTbZJBvuFbegZm47f9cDAdgxiehQfx
        gx import QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH
        gx import QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY
        gx import QmXRKBQA4wXP7xWbFiZsR1GP4HV6wMDQ1aWFxZZ4uBcPX9
        gx import QmQViVWBHbU6HmYjXcdNq7tVASCNgdg64ZGcauuDkLCivW


 * Comment out this line in imports:  `"github.com/libp2p/go-libp2p-host"` (not used)
 * In the following lines, replace `:=` with `=`:

        c, _ := cid.NewPrefixV1(cid.Raw, multihash.SHA2_256).Sum([]byte("libp2p-demo-chat"))

        tctx, _ := context.WithTimeout(ctx, time.Second*10)

        peers, err := dht.FindProviders(tctx, c)

 * Set `TopicName` variable by adding this as the first line of the `main()` function:

        TopicName := "libp2p-demo-chat"


# Instructions

Create `demo.go`.  Follow @why's tutorial and put all code in `demo.go`.

Build:  `go build -o libp2p-demo main.go`

Run:  `./libp2p-demo`
