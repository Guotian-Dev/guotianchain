# chain-node-api接口说明

1. 注册和enroll新的user在org1
curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Barry&orgName=org1'


{"success":true,"secret":"","message":"Barry enrolled Successfully","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw"}

2. 创建channel
curl -s -X POST \
  http://localhost:4000/channels \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json" \
  -d '{
    "channelName":"mychannel",
    "channelConfigPath":"../artifacts/channel/mychannel.tx"
}'


{"success":true,"message":"Channel 'mychannel' created Successfully"}
注意:authorization Bearer后跟之前的enroll的token

3. org1加入channel
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json" \
  -d '{
    "peers": ["peer1","peer2"]
}'

{"success":true,"message":"Successfully joined peers in organization org1 to the channel 'mychannel'"}


4. org1 安装channelcode
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json" \
  -d '{
    "peers": ["peer1", "peer2"],
    "chaincodeName":"mycc",
    "chaincodePath":"github.com/example_cc",
    "chaincodeVersion":"v0"
}'

Successfully Installed chaincode on organization org1

5. org1 的peer1 实例化channelcode
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json" \
  -d '{
    "chaincodeName":"mycc",
    "chaincodeVersion":"v0",
    "args":["a","100","b","200"]
}'

Chaincode Instantiation is SUCCESS[

6. org1的peer1 invoke channelcode
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json" \
  -d '{
    "fcn":"move",
    "args":["a","b","10"]
}'


1f0d78257600c8d22496cf44bfc258d0636df4701b908b74bd5fd4c7fee5271f

7. 查询chaincode
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/mycc?peer=peer1&fcn=query&args=%5B%22a%22%5D" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json"

a now has 90 after the move

8. 根据区块号码查询
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/blocks/1?peer=peer1" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json"

9. 根据transactionID查询交易
transactions/后的字符为invoke返回的transactionID

curl -s -X GET http://localhost:4000/channels/mychannel/transactions/1f0d78257600c8d22496cf44bfc258d0636df4701b908b74bd5fd4c7fee5271f?peer=peer1 \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json"


10. 获取ChainInfo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel?peer=peer1" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json"

11. 获取已安装的chaincode
curl -s -X GET \
  "http://localhost:4000/chaincodes?peer=peer1&type=installed" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json"

["name: mycc, version: v0, path: github.com/example_cc"]

12. 获取已实例化的chaincode
curl -s -X GET \
  "http://localhost:4000/chaincodes?peer=peer1&type=instantiated" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json"

["name: mycc, version: v0, path: github.com/example_cc"]

13. 获取channels
curl -s -X GET \
  "http://localhost:4000/channels?peer=peer1" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTc0ODY2ODAsInVzZXJuYW1lIjoiQmFycnkiLCJvcmdOYW1lIjoib3JnMSIsImlhdCI6MTUxNzQ1MDY4MH0.Ha6l8k5chOhBmwkiIpYVw--fW4ny-KUhH4cG14-kZLw" \
  -H "content-type: application/json"

{"channels":[{"channel_id":"mychannel"}]}


























