package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"time"

	process "gx/ipfs/QmSF8fPo3jgVBAy8fpdjjYqgG87dkJgUprRBHRd2tmfgpP/goprocess"

	"github.com/andlabs/ui"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corehttp"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/labstack/gommon/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	err := ui.Main(func() {
		name := ui.NewEntry()
		title := ui.NewLabel(" Referèndum 2017")
		status := ui.NewLabel(" IPFS is starting...")
		box := ui.NewVerticalBox()
		box.Disable()
		box.SetPadded(true)
		box.Append(name, false)
		box.Append(title, false)
		box.Append(status, false)
		window := ui.NewWindow("Referèndum 2017", 260, 200, false)
		window.SetChild(box)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
		ui.OnShouldQuit(func() bool {
			return true
		})
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}
		go func() {
			node := GetIPFSNode(path.Join(home, ".ipfs-simple"))
			if node.OnlineMode() {
				status.SetText("Running")
			} else {
				status.SetText("Error")
			}
			go func() {
				interval, err := time.ParseDuration("1s")
				if err != nil {
					panic(err)
				}
				tick := time.Tick(interval)
				for _ = range tick {
					peers := node.Peerstore.Peers()
					amount := strconv.Itoa(len(peers))
					fmt.Println(amount + " connected peers")
					log.Debug(amount + " connected peers")
					text := " IPFS is currently running and you\n are connected to " + amount + " other peers"
					status.SetText(text)
				}
			}()
			open.Run("http://localhost:9080/ipns/QmZxWEBJBVkGDGaKdYPQUXX4KC5TCWbvuR4iYZrTML8XCR/")
		}()
	})
	if err != nil {
		panic(err)
	}
}

func GetIPFSNode(repositoryPath string) *core.IpfsNode {
	hasRepo := fsrepo.IsInitialized(repositoryPath)
	if !hasRepo {
		ipfsConfig, err := config.Init(ioutil.Discard, 2048)
		if err != nil {
			panic(err)
		}
		ipfsConfig.Addresses = config.Addresses{
			Swarm: []string{
				"/ip4/0.0.0.0/tcp/0",
				"/ip4/0.0.0.0/tcp/0/ws",
				"/ip6/::/tcp/0",
			},
			API:     "",
			Gateway: "/ip4/127.0.0.1/tcp/9080",
		}
		// ipfsConfig.Bootstrap = []string{}
		err = fsrepo.Init(repositoryPath, ipfsConfig)
		if err != nil {
			panic(err)
		}
	}
	r, err := fsrepo.Open(repositoryPath)
	if err != nil {
		panic(err)
	}

	cfg := &core.BuildCfg{
		Repo:      r,
		Online:    true,
		Permament: true,
	}

	ctx := context.Background()
	nd, err := core.NewNode(ctx, cfg)
	nd.SetLocal(false)

	var opts = []corehttp.ServeOption{
		corehttp.GatewayOption(true, "/ipfs", "/ipns"),
		corehttp.WebUIOption,
	}
	proc := process.WithParent(process.Background())
	proc.Go(func(p process.Process) {
		if err := corehttp.ListenAndServe(nd, "/ip4/127.0.0.1/tcp/9080", opts...); err != nil {
			return
		}
	})

	if err != nil {
		panic(err)
	}
	return nd
}
