package chainHelper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xxRanger/blockchainUtil/chain"
	"github.com/xxRanger/blockchainUtil/contract"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/blockchainUtil/sender"
)

const (
	DEFAULT_TRANSFER_VALUE = "10000000000000000" // 0.01 ether
)


type ChainHandler struct {
	ManagerAccount      *sender.User
	Client              *chain.EthClient
	Contract contract.Contract
}

type AccountConfig struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

type ChainConfig struct {
	Account            AccountConfig `json:"account"`
	Port               string        `json:"port"`
	ContractAddress string        `json:"contractAddress"`
}

func NewChainHandler(config *ChainConfig) (*ChainHandler, error) {
	// eth client
	client, err := chain.NewEthClient(config.Port)
	if err != nil {
		panic(err)
	}
	// manager account
	pk, err := crypto.HexToECDSA(config.Account.PrivateKey)
	if err != nil {
		panic(err)
	}
	address := common.HexToAddress(config.Account.Address)
	managerAccount := sender.NewUser(address, pk)
	managerAccount.BindEthClient(client,sender.CHAIN_KIND_PRIVATE)

	// eth contract
	smartContract := nft.NewNFT(common.HexToAddress(config.ContractAddress)) //TODO use interface to general init contract
	smartContract.BindClient(client)

	handler := &ChainHandler{
		ManagerAccount:      managerAccount,
		Contract: smartContract,
		Client:              client,
	}
	return handler, nil
}
