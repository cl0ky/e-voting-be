package election

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RealBlockchainService struct {
	RPCUrl        string
	PrivateKeyHex string
	ContractAddr  string
	ContractABI   string
}

func (r *RealBlockchainService) StoreResultHash(hash string) (string, error) {
	client, err := ethclient.Dial(r.RPCUrl)
	if err != nil {
		return "", err
	}
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(r.PrivateKeyHex, "0x"))
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	parsedABI, err := abi.JSON(strings.NewReader(r.ContractABI))
	if err != nil {
		return "", err
	}
	contractAddress := common.HexToAddress(r.ContractAddr)
	input, err := parsedABI.Pack("storeHash", hash)
	if err != nil {
		return "", err
	}
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), 300000, gasPrice, input)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	log.Printf("TX sent: %s", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}
