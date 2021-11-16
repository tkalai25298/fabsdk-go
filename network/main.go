package main

import (
	// "encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// "reflect"
	"regexp"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspprov "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite/bccsp/sw"
	mspimpl "github.com/hyperledger/fabric-sdk-go/pkg/msp"

	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"encoding/json"
)

func main() {

	base := flag.String("ccp", "ccp/cli-ccp.yaml", "The ccp path to use.")
	asLocalhost := flag.Bool("localhost", false, "To set weather we want to set localhost.")
	// certPath := flag.String("cert_path","/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/ngp.cpu-network.com/users/Admin@ngp.cpu-network.com/msp/signcerts/cert.pem","The path to the MSP identity cert.")
	// keyPath := flag.String("key_path","/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/ngp.cpu-network.com/users/Admin@ngp.cpu-network.com/msp/keystore/priv_sk","The path to the MSP identity key.")
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

	Base_Using_Users_In_CCP(base)
	// Base_Using_Created_Signing_Identity(base,certPath,keyPath)

	// *base = "ccp/cli-msp-ccp.yml"

	// Custom_User_KVStore(base,certPath,keyPath)

	// Identity_Config(base)
}

func Base_Using_Users_In_CCP(base *string) {
	new, err := fabsdk.New(config.FromFile(*base))

	if err != nil {
		fmt.Printf("The error while creating fab context : %v \n", err)
		os.Exit(1)
	}

	defer new.Close()

	fmt.Println("New Fabric context created")

	chctx := new.ChannelContext("ngpchannel", fabsdk.WithOrg("ngpMSP"), fabsdk.WithUser("admin"))
	clctx := new.Context(fabsdk.WithOrg("ngpMSP"))

	_, err = msp.New(clctx)
	if err != nil {
		fmt.Printf("The error while creating ledger context : %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("   - New Fabric msp-client context created \n")

	ld, err := ledger.New(chctx)
	if err != nil {
		fmt.Printf("The error while creating ledger context : %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("  - New ledger context created \n")

	info, err := ld.QueryInfo()
	if err != nil {
		fmt.Printf("Unable to get info: %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("The info is %+v\n", info)

	block, err := ld.QueryBlock(5)
	ld.QueryInfo()
	if err != nil {
		fmt.Printf("failed to query block: %s\n", err)
	}
	fmt.Printf("     - New block\n")

	blockData := block.Data
	if blockData != nil {
		data := blockData.String()
		// fmt.Printf("Retrieved block data #1 %+v \n", data)

		// re := regexp.MustCompile(`te.*?set`)
		re := regexp.MustCompile(`{.*?}}`)
		fmt.Printf("pattern %v\n", re.String())
		fmt.Println("match :",re.MatchString(data))

		fetchstring := re.FindAllString(data, -1)
		// fmt.Printf("fetch string %s\n", fetchstring)

		
		s := string(fetchstring[0])

		str := strings.ReplaceAll(s, "\\", "")

		fmt.Printf("json: %s\n",str)

		sec := map[string]interface{}{}
		if err := json.Unmarshal([]byte(str), &sec); err != nil {
			println("err:",err.Error())
		}
		fmt.Println(sec)

	}
}

func Base_Using_Created_Signing_Identity(base, certPath, keyPath *string) {

	new, err := fabsdk.New(config.FromFile(*base))

	if err != nil {
		fmt.Printf("The error while creating fab context : %v \n", err)
		os.Exit(1)
	}

	defer new.Close()

	fmt.Println("New Fabric context created")

	clctx := new.Context(fabsdk.WithOrg("ngpMSP"))

	newMSP, err := msp.New(clctx)
	if err != nil {
		fmt.Printf("The error while creating ledger context : %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("   - New Fabric msp-client context created \n")

	fmt.Println("New Fabric with signing identity")

	cert, err := ioutil.ReadFile(*certPath)
	if err != nil {
		fmt.Printf("The error while reading cert : %v \n", err)
		os.Exit(1)
	}
	key, err := ioutil.ReadFile(*keyPath)
	if err != nil {
		fmt.Printf("The error while reading key : %v \n", err)
		os.Exit(1)
	}

	identity, err := newMSP.CreateSigningIdentity(mspprov.WithCert(cert), mspprov.WithPrivateKey(key))

	if err != nil {
		fmt.Printf("The error while creating identity : %v \n", err)
		os.Exit(1)
	}

	hctx := new.ChannelContext("ngpchannel", fabsdk.WithIdentity(identity))

	ldnew, err := ledger.New(hctx)
	if err != nil {
		fmt.Printf("The error while creating ledger context : %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("  - New ledger context created \n")

	info, err := ldnew.QueryInfo()
	if err != nil {
		fmt.Printf("Unable to get info: %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("The info new is == %+v\n", info)
}

func Custom_User_KVStore(base, certPath, keyPath *string) {
	fmt.Printf("Using the following ccp %s", *base)

	_, err := sw.GetSuite(256, "sha2", mspimpl.NewMemoryKeyStore([]byte("hi")))

	new, err := fabsdk.New(config.FromFile(*base))

	if err != nil {
		fmt.Printf("The error while creating fab context : %v \n", err)
		os.Exit(1)
	}

	defer new.Close()

	fmt.Println("New Fabric context created")

	clctx := new.Context(fabsdk.WithOrg("ngpMSP"))

	newMSP, err := msp.New(clctx)
	if err != nil {
		fmt.Printf("The error while creating ledger context : %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("   - New Fabric msp-client context created \n")

	fmt.Println("New Fabric with signing identity")

	cert, err := ioutil.ReadFile(*certPath)
	if err != nil {
		fmt.Printf("The error while reading cert : %v \n", err)
		os.Exit(1)
	}
	key, err := ioutil.ReadFile(*keyPath)
	if err != nil {
		fmt.Printf("The error while reading key : %v \n", err)
		os.Exit(1)
	}

	identity, err := newMSP.CreateSigningIdentity(mspprov.WithCert(cert), mspprov.WithPrivateKey(key))

	if err != nil {
		fmt.Printf("The error while creating identity : %v \n", err)
		os.Exit(1)
	}

	hctx := new.ChannelContext("ngpchannel", fabsdk.WithIdentity(identity))

	ldnew, err := ledger.New(hctx)
	if err != nil {
		fmt.Printf("The error while creating ledger context : %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("  - New ledger context created \n")

	info, err := ldnew.QueryInfo()
	if err != nil {
		fmt.Printf("Unable to get info: %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("The info new is == %+v\n", info)
}

func Identity_Config(base *string) {
	configBackend, err := config.FromFile(*base)()

	if err != nil {
		fmt.Println(err)
	}

	identityCfg, err := mspimpl.ConfigFromBackend(configBackend...)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("The identity config %+v", identityCfg)
}