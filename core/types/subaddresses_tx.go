// (c) 2019-2021, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//go:generate go run github.com/fjl/gencodec -type AccessTuple -out gen_access_tuple.go


// SubAddressesTx is the data of EIP-2930 access list transactions.
type SubAddressesTx struct {
	ChainID    *big.Int        // destination chain ID
	Nonce      uint64          // nonce of sender account
	GasPrice   *big.Int        // wei per gas
	Gas        uint64          // gas limit
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int        // wei amount
	Data       []byte          // contract invocation input data
	AccessList AccessList
	UseSubAddresses bool       // allow spend coins from subaddresses
	V, R, S    *big.Int        // signature values
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *SubAddressesTx) copy() TxData {
	cpy := &SubAddressesTx{
		Nonce: tx.Nonce,
		To:    copyAddressPtr(tx.To),
		Data:  common.CopyBytes(tx.Data),
		Gas:   tx.Gas,
		// These are copied below.
		AccessList: make(AccessList, len(tx.AccessList)),
		UseSubAddresses: false,
		Value:      new(big.Int),
		ChainID:    new(big.Int),
		GasPrice:   new(big.Int),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}

	copy(cpy.AccessList, tx.AccessList)
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasPrice != nil {
		cpy.GasPrice.Set(tx.GasPrice)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *SubAddressesTx) txType() byte           { return SubAddressesTxType }
func (tx *SubAddressesTx) chainID() *big.Int      { return tx.ChainID }
func (tx *SubAddressesTx) accessList() AccessList { return tx.AccessList }
func (tx *SubAddressesTx) data() []byte           { return tx.Data }
func (tx *SubAddressesTx) gas() uint64            { return tx.Gas }
func (tx *SubAddressesTx) gasPrice() *big.Int     { return tx.GasPrice }
func (tx *SubAddressesTx) gasTipCap() *big.Int    { return tx.GasPrice }
func (tx *SubAddressesTx) gasFeeCap() *big.Int    { return tx.GasPrice }
func (tx *SubAddressesTx) value() *big.Int        { return tx.Value }
func (tx *SubAddressesTx) nonce() uint64          { return tx.Nonce }
func (tx *SubAddressesTx) to() *common.Address    { return tx.To }
func (tx *SubAddressesTx) useSubAddresses() bool    { return tx.UseSubAddresses }

func (tx *SubAddressesTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *SubAddressesTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}
