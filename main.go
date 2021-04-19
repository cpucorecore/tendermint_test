package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	chainID                = "abc"
	privValidatorKeyFile   = "./test_data/priv_validator_key.json"
	privValidatorStateFile = "./test_data/priv_validator_state.json"

	blockHashStr      = "0x426D04A90896F55308E1B708ADBC4AFAA36762B7640705C7CC1F45AFF991A29C"
	blockPartsHashStr = "0x753B1095804F9A6E88BB443BB1C7A41A5E78D265D78BF291299BD342ACBB4CD4"

	blockSigStr = "WoeNkFpbKOkz5fh+1dpnBMQ4SuHpzuPOFB0hE52t35EY7w6Ti9k16v5O8B8dwMHW1Tkogk8riGXHogFjMctuDg=="
)

func main() {
	filePV := privval.LoadFilePV(privValidatorKeyFile, privValidatorStateFile)

	t, _ := time.Parse("2006-01-02T15:04:05.000000Z", "2021-04-14T09:25:27.209219Z")

	blockHash, err := hexutil.Decode(blockHashStr)
	if err != nil {
		fmt.Errorf("%s", err)
	}

	blockPartsHash, err := hexutil.Decode(blockPartsHashStr)
	if err != nil {
		fmt.Errorf("%s", err)
	}

	fmt.Printf("blockHash: [%x]\n", blockHash)
	fmt.Printf("blockPartsHash: [%x]\n", blockPartsHash)

	var partsHash common.HexBytes
	err = (&partsHash).Unmarshal(blockPartsHash)
	if err != nil {
		fmt.Errorf("%s\n", err)
	}

	addr := filePV.GetPubKey().Address()
	vote := types.Vote{
		ValidatorAddress: addr,
		ValidatorIndex:   0,
		Height:           140188,
		Round:            0,
		Timestamp:        t,
		Type:             types.PrecommitType,
		BlockID:          types.BlockID{Hash: blockHash, PartsHeader: types.PartSetHeader{Total: 1, Hash: blockPartsHash}},
	}

	fmt.Printf("vote: %s\n", vote.String())
	fmt.Printf("signBytes: [%x]\n", vote.SignBytes(chainID))
	err = filePV.SignVote(chainID, &vote)
	if err != nil {
		fmt.Errorf("%s\n", err)
	}

	fmt.Printf("sig: [%x]\n", vote.Signature)

	blockSig, err := base64.StdEncoding.DecodeString(blockSigStr)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	fmt.Printf("actualBlockSig: [%x]\n", blockSig)
}
