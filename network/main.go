package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"

	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource/genesisconfig"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"gopkg.in/yaml.v2"
)

type Configtx struct {
	Organizations	genesisconfig.Organization
	Capabilities	Capabilities
	Application		genesisconfig.Application
	Orderer			genesisconfig.Orderer
	Channel			Channel
	Profiles		genesisconfig.Profile
}

type Capabilities struct {
	Channel map[string]bool
	Orderer map[string]bool
	Application map[string]bool
}

type Channel struct {
	Policies	[]map[string]genesisconfig.Policy
	Capabilities Capabilities
}

func main() {

	// base := flag.String("ccp", "ccp/cli-ccp.yaml", "The ccp path to use.")
	asLocalhost := flag.Bool("localhost", false, "To set weather we want to set localhost.")
	channelName := "ngpchannel"
	// configPath	:= "/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/channel.tx"
	flag.Parse()

	if *asLocalhost {
		// DISCOVERY_AS_LOCALHOST
		fmt.Println("Setting env for local")
		err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

		if err != nil {
			fmt.Println("Unable to set \"DISCOVERY_AS_LOCALHOST\" as \"true\"")
			os.Exit(1)
		}
	}

	// Create_Join_Channel(base,channelName,configPath)
	CreateArtifacts(channelName)
}

func CreateArtifacts(channelID string) {

	content, err := ioutil.ReadFile("configtx.yaml")
	if err != nil {
		fmt.Printf("reading configtx file err: %v",err)
		os.Exit(1)
	}
	
	configtx := Configtx{}
	err = yaml.Unmarshal([]byte(content),configtx)
	if err != nil {
		fmt.Printf("Unmarshalling configtx.yaml failed: %v",err)
		os.Exit(1)
	}

	fmt.Printf("successfully Unmarshalled configtx \n",configtx.Profiles.Orderer.Addresses)
}

func Create_Join_Channel(base *string,channelName string,configPath string) {
	new, err := fabsdk.New(config.FromFile(*base))

	if err != nil {
		fmt.Printf("The error while creating fab context : %v \n", err)
		os.Exit(1)
	}

	defer new.Close()

	fmt.Println("New Fabric context created")

	// chctx := new.ChannelContext("ngpchannel", fabsdk.WithOrg("ngpMSP"), fabsdk.WithUser("admin"))
	clctx := new.Context(fabsdk.WithOrg("ngpMSP"), fabsdk.WithUser("admin"))

	ch, err := resmgmt.New(clctx)
	if err != nil {
		fmt.Printf("The error while creating resource management client : %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("   - New Fabric resource mgmt-client context created \n")

	//creating reader for configtx file
	r, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("failed to open channel config: %s\n", err)
	}
	defer r.Close()

	_,err = ch.SaveChannel(
		resmgmt.SaveChannelRequest{
			ChannelID: channelName,
			ChannelConfig: r,
	})

	if err != nil {
		fmt.Printf("The error while creating channel : %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("	- Create channel successful \n")

	err = ch.JoinChannel(channelName)
	

	if err != nil {
		fmt.Printf("The error while joining channel : %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("	- Join channel successful \n")
}

