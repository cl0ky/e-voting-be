package election

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
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

func (r *RealBlockchainService) GetHash(electionId string) (string, error) {
	return r.GetStoredHash(electionId)
}

func (r *RealBlockchainService) StoreResultHash(electionId string, hash string) (string, error) {

	log.Printf("[BLOCKCHAIN] RPCUrl: %s", r.RPCUrl)
	log.Printf("[BLOCKCHAIN] PrivateKeyHex: %s", r.PrivateKeyHex)
	log.Printf("[BLOCKCHAIN] ContractAddr: %s", r.ContractAddr)
	log.Printf("[BLOCKCHAIN] ContractABI: %s", r.ContractABI)
	log.Printf("[BLOCKCHAIN] ElectionId: %s", electionId)
	log.Printf("[BLOCKCHAIN] Hash: %s", hash)

	// 1. Connect RPC
	client, err := ethclient.Dial(r.RPCUrl)
	if err != nil {
		return "", err
	}

	// 2. Private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(r.PrivateKeyHex, "0x"))
	if err != nil {
		return "", err
	}

	// 3. Derive sender address
	publicKey := privateKey.Public()
	pubKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("publicKey is not *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*pubKeyECDSA)
	log.Printf("[BLOCKCHAIN] From Address: %s", fromAddress.Hex())

	// 4. Nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	// 5. Gas params (EIP-1559)
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return "", err
	}
	gasFeeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	// 6. Chain ID → Ganache = 1337
	chainID := big.NewInt(1337)

	// 7. Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(r.ContractABI))
	if err != nil {
		return "", err
	}

	// 8. Encode function call
	input, err := parsedABI.Pack("storeHash", electionId, hash)
	if err != nil {
		return "", err
	}

	contractAddress := common.HexToAddress(r.ContractAddr)

	// 9. Create EIP-1559 DynamicFeeTx
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &contractAddress,
		Value:     big.NewInt(0),
		Gas:       uint64(300000),
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Data:      input,
	})

	// 10. Sign EIP-1559 transaction
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return "", err
	}

	// 11. Send
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	txHash := signedTx.Hash().Hex()
	log.Printf("[BLOCKCHAIN] TX Sent: %s", txHash)

	return txHash, nil
}

func (r *RealBlockchainService) GetStoredHash(electionId string) (string, error) {
	client, err := ethclient.Dial(r.RPCUrl)
	if err != nil {
		return "", err
	}

	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(r.ContractABI))
	if err != nil {
		return "", err
	}

	contractAddress := common.HexToAddress(r.ContractAddr)

	// Encode view call
	data, err := parsedABI.Pack("getHash", electionId)
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	// Execute view call
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	var storedHash string

	err = parsedABI.UnpackIntoInterface(&storedHash, "getHash", result)
	if err != nil {
		return "", err
	}

	return storedHash, nil
}
