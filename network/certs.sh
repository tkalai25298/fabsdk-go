export PATH=/home/kalai/Downloads/bin:$PATH
cryptogen generate --config=./crypto-config.yaml
configtxgen -profile OrdererGenesis -channelID cpu-sys-channel -outputBlock ./channel-artifacts/genesis.block -configPath .
configtxgen -profile ngpchannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID ngpchannel
configtxgen -profile ngpchannel -outputAnchorPeersUpdate ./channel-artifacts/ngpMSPanchors.tx -channelID ngpchannel -asOrg ngpMSP