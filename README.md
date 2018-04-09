# guotianchain
GuoTian Block Chain based Hyperledger Fabric.

# 区块链环境安装
1. C++
2. Golang
3. git
4. docker
5. nodejs


# 下载镜像以及工具包

curl -sSL https://goo.gl/6wtTN5 | bash -s 1.1.0


# 导出工具的路径到PATH中
export PATH=<path to download location>/bin:$PATH


# 项目结构设计
chain-network   --- 区块链网络拓扑结构
chain-node-api  --- 基于fabric-samples/balance-transfer例子修改而来
chaincode       --- 链码目录
docs            --- 文档说明  