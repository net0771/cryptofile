# 1. channel
1. export CHANNEL_NAME=mychannel  
2. ../bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
