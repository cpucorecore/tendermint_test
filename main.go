package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	chainID                = "abc"
	privValidatorKeyFile   = "./test_data/priv_validator_key.json"
	privValidatorStateFile = "./test_data/priv_validator_state.json"

	blockHashStr      = "0x5924ED59C4184241253AB480FB722797D44DB86C843F4C0953513771C63F74E7"
	blockPartsHashStr = "0x17529FD19FE6DE9B92A64C1CCF9D6E367C0FE9FAF982A850D858F2C59C207CB8"

	blockSigStr = "A7IMN3B47bB0MmCTh/sv7gpAvtfTOKvIkpkEICFsC2xG7PoRLT5RRuqKAI4vehTg/r5da1+265fZlnBh/ERkAg=="
)

func MustNoErr(err error) {
	if err != nil {
		fmt.Println(fmt.Errorf("err: [%s]", err.Error()))
		os.Exit(-1)
	}
}

func main() {
	filePV := privval.LoadFilePV(privValidatorKeyFile, privValidatorStateFile)

	t, err := time.Parse("2006-01-02T15:04:05.000000Z", "2021-04-14T09:25:27.209219Z")
	MustNoErr(err)

	blockHash, err := hexutil.Decode(blockHashStr)
	MustNoErr(err)

	blockPartsHash, err := hexutil.Decode(blockPartsHashStr)
	MustNoErr(err)

	addr := filePV.GetPubKey().Address()
	vote := types.Vote{
		ValidatorAddress: addr,
		ValidatorIndex:   0,
		Height:           16235,
		Round:            0,
		Timestamp:        t,
		Type:             types.PrecommitType,
		BlockID:          types.BlockID{Hash: blockHash, PartsHeader: types.PartSetHeader{Total: 1, Hash: blockPartsHash}},
	}

	fmt.Printf("signBytes: [%x]\n", vote.SignBytes(chainID))
	err = filePV.SignVote(chainID, &vote)
	MustNoErr(err)

	fmt.Printf("sig: [%x]\n", vote.Signature)

	blockSig, err := base64.StdEncoding.DecodeString(blockSigStr)
	MustNoErr(err)
	fmt.Printf("actualBlockSig: [%x]\n", blockSig)
}
