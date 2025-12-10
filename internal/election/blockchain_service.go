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

	log.Printf("[BLOCKCHAIN] RPCUrl: %s", r.RPCUrl)
	log.Printf("[BLOCKCHAIN] PrivateKeyHex: %s", r.PrivateKeyHex)
	log.Printf("[BLOCKCHAIN] ContractAddr: %s", r.ContractAddr)
	log.Printf("[BLOCKCHAIN] ContractABI: %s", r.ContractABI)
	log.Printf("[BLOCKCHAIN] Hash: %s", hash)

	// 1. Connect to Ganache
	client, err := ethclient.Dial(r.RPCUrl)
	if err != nil {
		log.Printf("[BLOCKCHAIN] ethclient.Dial error: %v", err)
		return "", err
	}

	// 2. Load private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(r.PrivateKeyHex, "0x"))
	if err != nil {
		log.Printf("[BLOCKCHAIN] HexToECDSA error: %v", err)
		return "", err
	}

	// 3. Get sender address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Printf("[BLOCKCHAIN] publicKey type assertion failed")
		return "", errors.New("cannot assert type: publicKey is not *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Printf("[BLOCKCHAIN] fromAddress: %s", fromAddress.Hex())

	// 4. Nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Printf("[BLOCKCHAIN] PendingNonceAt error: %v", err)
		return "", err
	}

	// 5. Gas parameters
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Printf("[BLOCKCHAIN] SuggestGasTipCap error: %v", err)
		return "", err
	}

	gasFeeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("[BLOCKCHAIN] SuggestGasPrice error: %v", err)
		return "", err
	}

	// 6. Chain ID
	chainID := big.NewInt(1337) // Ganache chainID
	log.Printf("[BLOCKCHAIN] chainID: %s", chainID.String())

	// 7. Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(r.ContractABI))
	if err != nil {
		log.Printf("[BLOCKCHAIN] abi.JSON error: %v", err)
		return "", err
	}

	// 8. Prepare contract call input
	contractAddress := common.HexToAddress(r.ContractAddr)
	input, err := parsedABI.Pack("storeHash", hash)
	if err != nil {
		log.Printf("[BLOCKCHAIN] abi.Pack error: %v", err)
		return "", err
	}

	// 9. Create EIP-1559 dynamic fee transaction
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

	// 10. Sign transaction (EIP-1559 compatible signer)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		log.Printf("[BLOCKCHAIN] SignTx error: %v", err)
		return "", err
	}

	// 11. Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Printf("[BLOCKCHAIN] SendTransaction error: %v", err)
		return "", err
	}

	txHash := signedTx.Hash().Hex()
	log.Printf("[BLOCKCHAIN] TX sent: %s", txHash)

	return txHash, nil
}
