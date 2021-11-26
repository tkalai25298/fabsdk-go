package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource/genesisconfig"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func channelCapabilities() map[string]bool {
	return map[string]bool{
		"V2_0": true,
	}
}

func channelDefaults() (map[string]*genesisconfig.Policy, map[string]bool) {

	policies := map[string]*genesisconfig.Policy{
		"Admins": {
			Type: "ImplicitMeta",
			Rule: "ANY Admins",
		},
		"Readers": {
			Type: "ImplicitMeta",
			Rule: "ANY Readers",
		},
		"Writers": {
			Type: "ImplicitMeta",
			Rule: "ANY Writers",
		},
	}
	return policies, channelCapabilities()
}

func applicationDefaults() *genesisconfig.Application {

	_, capabilities := channelDefaults()

	return &genesisconfig.Application{
		Organizations: []*genesisconfig.Organization{},
		Policies: map[string]*genesisconfig.Policy{
			"LifecycleEndorsement": {
				Type: "ImplicitMeta",
				Rule: "ANY Endorsement",
			},
			"Endorsement": {
				Type: "ImplicitMeta",
				Rule: "ANY Endorsement",
			},
			"Readers": {
				Type: "ImplicitMeta",
				Rule: "ANY Readers",
			},
			"Writers": {
				Type: "ImplicitMeta",
				Rule: "ANY Writers",
			},
			"Admins": {
				Type: "ImplicitMeta",
				Rule: "ANY Admins",
			},
		},
		Capabilities: capabilities,
	}
}

func ngpPolicies() map[string]*genesisconfig.Policy {
	return map[string]*genesisconfig.Policy{
		"Readers": {
			Type: "Signature",
			Rule: "OR('ngpMSP.admin', 'ngpMSP.peer', 'ngpMSP.client')",
		},
		"Writers": {
			Type: "Signature",
			Rule: "OR('ngpMSP.admin', 'ngpMSP.client')",
		},
		"Admins": {
			Type: "Signature",
			Rule: "OR('ngpMSP.admin')",
		},
		"Endorsement": {
			Type: "Signature",
			Rule: "OR('ngpMSP.peer','ngpMSP.admin')",
		},
	}
}

func ngpOrg() *genesisconfig.Organization {
	return &genesisconfig.Organization{
		Name:          "ngpMSP",
		ID:            "ngpMSP",
		MSPDir:        filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer", "crypto-config/peerOrganizations/ngp.cpu-network.com/msp"),
		MSPType:       "bccsp",
		Policies:      ngpPolicies(),
		AnchorPeers: []*genesisconfig.AnchorPeer{
			{
				Host: "peer0.ngp.cpu-network.com",
				Port: 7051,
			},
		},
	}
}

func sampleSingleMSPChannel() *genesisconfig.Profile {

	policies, _ := channelDefaults()
	appDefaults := applicationDefaults()
	appDefaults.Organizations = []*genesisconfig.Organization{
		ngpOrg(),
	}

	return &genesisconfig.Profile{
		Policies:    policies,
		Consortium:  "cpuConsortium",
		Application: appDefaults,
	}
}

func main() {

	// base := flag.String("ccp", "ccp/cli-ccp.yaml", "The ccp path to use.")
	asLocalhost := flag.Bool("localhost", false, "To set weather we want to set localhost.")
	channelName := "ngpchannel"
	// configPath	:= "/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/channel.tx"
	configtxPath := "/opt/gopath/src/github.com/hyperledger/fabric/peer/configtx.yaml"
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
	
	// CreateArtifacts(channelName)
	Create_Artifacts(configtxPath,channelName)
	// Create_Join_Channel(base,channelName,configPath)

}

func CreateArtifacts(channelID string) {
	config := sampleSingleMSPChannel()
	fmt.Printf("config: %#v",config)

	configtx, err := resource.CreateChannelCreateTx(config, nil, channelID)

	if err != nil {
		fmt.Printf("The error while creating channel tx : %v \n", err)
		os.Exit(1)
	}

	f,err := os.Create("channel-artifacts/channel.tx")
	if err != nil {
		fmt.Printf("The error while creating channel tx file: %v \n", err)
		os.Exit(1)
	}

	_,err = f.Write(configtx)
	if err != nil {
		fmt.Printf("The error while writing channel tx file: %v \n", err)
		os.Exit(1)
	}
}

func Create_Artifacts(configtxfile string,channelID string) {
	configtx,err := ioutil.ReadFile(configtxfile)

	if err != nil {
		fmt.Printf("The error while reading configtx file : %v \n", err)
		os.Exit(1)
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(configtx))
	if err != nil {
		fmt.Printf("The error while reading config from configtx file : %v \n", err)
		os.Exit(1)
	}
	// fmt.Printf("profile config: %v",viper.Get("Profiles.ngpchannel"))

	profile := &genesisconfig.Profile{}

	err = mapstructure.Decode(viper.Get("Profiles.ngpchannel"),profile)

	if err != nil {
		fmt.Printf("The error while decoding map structure : %v \n", err)
		os.Exit(1)
	}

	policy := make(map[string]*genesisconfig.Policy)

	err = mapstructure.Decode(viper.Get("Profiles.ngpchannel.Application.Policies"),&policy)

	if err != nil {
		fmt.Printf("The error while decoding map structure : %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("policy: %#v",policy)

	config, err := resource.CreateChannelCreateTx(profile, nil, channelID)

	if err != nil {
		fmt.Printf("The error while creating channel tx : %v \n", err)
		os.Exit(1)
	}

	f,err := os.Create("channel-artifacts/channel.tx")
	if err != nil {
		fmt.Printf("The error while creating channel tx file: %v \n", err)
		os.Exit(1)
	}

	_,err = f.Write(config)
	if err != nil {
		fmt.Printf("The error while writing channel tx file: %v \n", err)
		os.Exit(1)
	}

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

