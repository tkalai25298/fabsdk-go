#Creating channel
echo "######### Creating channel #########"
peer channel create -o orderer.cpu-network.com:7050 -c ngpchannel -f ./channel-artifacts/channel.tx 
CORE_PEER_ADDRESS=peer0.ngp.cpu-network.com:7051
#joining peer
echo "######### Joining peer to ngpchannel #########"
peer channel join -b ngpchannel.block
peer channel update -o orderer.cpu-network.com:7050 -c ngpchannel -f ./channel-artifacts/ngpMSPanchors.tx
#chaincode installation
# echo "######### cloning the chainode from git #########"
#  git clone git@github.com:powerofn/enegry-consumption-chaincode.git
 peer lifecycle chaincode package cpu.tar.gz --path ./enegry-consumption-chaincode/ --lang golang --label cpu

echo "######### Installing chaincode #########"
 peer lifecycle chaincode install cpu.tar.gz
 peer lifecycle chaincode queryinstalled

 echo "######### Getting package id #########"
 export CC_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep Package | cut -d \  -f 3 | cut -d , -f 1)

 echo "######### Approving org #########"
 peer lifecycle chaincode approveformyorg -o orderer.cpu-network.com:7050 --channelID ngpchannel --name cpu --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1
 peer lifecycle chaincode checkcommitreadiness --channelID ngpchannel --name cpu --version 1.0 --sequence 1 --output json
 
 echo "######### chaincode commit #########"
 peer lifecycle chaincode commit -o orderer.cpu-network.com:7050 --channelID ngpchannel --name cpu --version 1.0 --sequence 1
 peer lifecycle chaincode querycommitted --channelID ngpchannel --name cpu

#  echo "#########Invoking Chaincode - create #########"
#  peer chaincode invoke -o orderer.cpu-network.com:7050 -C ngpchannel -n cpu -c '{"Args":["AddCPU","testasset","mpan"]}'

# peer chaincode invoke -o orderer.cpu-network.com:7050 -C ngpchannel -n cpu -c '{"Args":["AddCPU","testasset-tx","mpan"]}'

#  echo "#########Invoking Chaincode - list #########"
#  peer chaincode invoke -o orderer.cpu-network.com:7050 -C ngpchannel -n cpu -c '{"function":"GetHistory","Args":["testasset"]}'

# peer chaincode invoke -o orderer.cpu-network.com:7050 -C ngpchannel -n cpu -c '{"Args":["AddUsage","testasset","mpan","mac_id","time","[{\"phaseID\":1,\"kwh\":2892.27}","{\"phaseID\":2,\"kwh\":2892.27}","{\"phaseID\":3,\"kwh\":2892.27}]"]}'

peer chaincode invoke -o orderer.cpu-network.com:7050 -C ngpchannel -n cpu -c '{"Args":["AddCPU","testasset-tx-1","mpan"]}'
# sleep 2
# peer chaincode invoke -o orderer.cpu-network.com:7050 -C ngpchannel -n cpu -c '{"function":"AddUsage","Args":["testasset-tx-1","mpan","macid","0001-01-01T00:00:00","{\"phaseID\":1,\"kwh\":2892.27}","{\"phaseID\":2,\"kwh\":2892.27}","{\"phaseID\":3,\"kwh\":2892.3}"]}'

# peer chaincode invoke -o orderer.cpu-network.com:7050 -C ngpchannel -n cpu -c '{"Args":["GetHistory","testasset-tx-1","mpan"]}'
