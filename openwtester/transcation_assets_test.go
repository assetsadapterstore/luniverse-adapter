/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package openwtester

import (
	"github.com/blocktree/openwallet/v2/openw"
	"testing"

	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
)

func testGetAssetsAccountBalance(tm *openw.WalletManager, walletID, accountID string) {
	balance, err := tm.GetAssetsAccountBalance(testApp, walletID, accountID)
	if err != nil {
		log.Error("GetAssetsAccountBalance failed, unexpected error:", err)
		return
	}
	log.Info("balance:", balance)
}

func testGetAssetsAccountTokenBalance(tm *openw.WalletManager, walletID, accountID string, contract openwallet.SmartContract) {
	balance, err := tm.GetAssetsAccountTokenBalance(testApp, walletID, accountID, contract)
	if err != nil {
		log.Error("GetAssetsAccountTokenBalance failed, unexpected error:", err)
		return
	}
	log.Info("token balance:", balance.Balance)
}

func testCreateTransactionStep(tm *openw.WalletManager, walletID, accountID, to, amount, feeRate string, contract *openwallet.SmartContract) (*openwallet.RawTransaction, error) {

	//err := tm.RefreshAssetsAccountBalance(testApp, accountID)
	//if err != nil {
	//	log.Error("RefreshAssetsAccountBalance failed, unexpected error:", err)
	//	return nil, err
	//}

	rawTx, err := tm.CreateTransaction(testApp, walletID, accountID, amount, to, feeRate, "", contract, nil)

	if err != nil {
		log.Error("CreateTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTx, nil
}

func testCreateSummaryTransactionStep(
	tm *openw.WalletManager,
	walletID, accountID, summaryAddress, minTransfer, retainedBalance, feeRate string,
	start, limit int,
	contract *openwallet.SmartContract,
	feeSupportAccount *openwallet.FeesSupportAccount) ([]*openwallet.RawTransactionWithError, error) {

	rawTxArray, err := tm.CreateSummaryRawTransactionWithError(testApp, walletID, accountID, summaryAddress, minTransfer,
		retainedBalance, feeRate, start, limit, contract, feeSupportAccount)

	if err != nil {
		log.Error("CreateSummaryTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTxArray, nil
}

func testSignTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	_, err := tm.SignTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, "12345678", rawTx)
	if err != nil {
		log.Error("SignTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testVerifyTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	//log.Info("rawTx.Signatures:", rawTx.Signatures)

	_, err := tm.VerifyTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("VerifyTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testSubmitTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	tx, err := tm.SubmitTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("SubmitTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Std.Info("tx: %+v", tx)
	log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	return rawTx, nil
}

func TestTransfer_LUK(t *testing.T) {

	addrs := []string{
		"0x278c99037364a963e565b0c5cd17a69d6548c1ce",
		"0x859fb1ecf956b6e3c40544c2dbb3b0126383a0ec",
		"0xa41b5f136add3b19d1fe576a5157d03674f2b446",
		"0xb34fb7e104f99c0509b9d9e37127a5f9cdfebc66",
		"0xbc22aa0d50a9c4b3dd0e960c0286075d6c5619af",
		"0xd65772dc4d58fa901f8b883a45c15fd8f8d5c46a",

		//"0x8463accbf09e600434716677975a4b5bd8271d78", //手续费地址
	}

	tm := testInitWalletManager()
	walletID := "WLJPogGRRxso9qfA1k49XTXuJDATkNMnk4"
	accountID := "CYbnSBikZG6dZRKr8rvW6d8KzdNX9AjUJX98BcbXkthL"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	for _, to := range addrs {
		rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "2", "", nil)
		if err != nil {
			return
		}

		log.Std.Info("rawTx: %+v", rawTx)

		_, err = testSignTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTx)
		if err != nil {
			return
		}
	}
}


func TestTransfer_ERC20(t *testing.T) {

	addrs := []string{
		"0x278c99037364a963e565b0c5cd17a69d6548c1ce",
		"0x859fb1ecf956b6e3c40544c2dbb3b0126383a0ec",
		"0xa41b5f136add3b19d1fe576a5157d03674f2b446",
		"0xb34fb7e104f99c0509b9d9e37127a5f9cdfebc66",
		"0xbc22aa0d50a9c4b3dd0e960c0286075d6c5619af",
		"0xd65772dc4d58fa901f8b883a45c15fd8f8d5c46a",
	}

	tm := testInitWalletManager()
	walletID := "WLJPogGRRxso9qfA1k49XTXuJDATkNMnk4"
	accountID := "CYbnSBikZG6dZRKr8rvW6d8KzdNX9AjUJX98BcbXkthL"

	contract := openwallet.SmartContract{
		Address:  "0x35def49e4c26aadd2a0734f38f07dc6c7993f764",
		Symbol:   "LUK",
		Name:     "ZIK",
		Token:    "ZIK",
		Decimals: 18,
	}

	testGetAssetsAccountBalance(tm, walletID, accountID)

	testGetAssetsAccountTokenBalance(tm, walletID, accountID, contract)

	for _, to := range addrs {
		rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "1", "", &contract)
		if err != nil {
			return
		}

		log.Std.Info("rawTx: %+v", rawTx)

		_, err = testSignTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

	}

}

func TestSummary_LUK(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WLJPogGRRxso9qfA1k49XTXuJDATkNMnk4"
	accountID := "CoDfSV7i58GKJVSoKBFwYULDXuLpLXsoG9cmF34MRWm7"
	summaryAddress := "0x3935ed7c62e31da9cc8f6780ffd16194541f9c59"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "", "", "",
		0, 100, nil, nil)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return
	}

	//执行汇总交易
	for _, rawTxWithErr := range rawTxArray {

		if rawTxWithErr.Error != nil {
			log.Error(rawTxWithErr.Error.Error())
			continue
		}

		_, err = testSignTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}
	}

}


func TestSummary_ERC20(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WLJPogGRRxso9qfA1k49XTXuJDATkNMnk4"
	accountID := "CoDfSV7i58GKJVSoKBFwYULDXuLpLXsoG9cmF34MRWm7"
	summaryAddress := "0x3935ed7c62e31da9cc8f6780ffd16194541f9c59"

	feesSupport := openwallet.FeesSupportAccount{
		AccountID: "HyN6dFbjp95HfpdG9autzmZgkcnbLrCQ2P6GfJBqcQiG",
		//FixSupportAmount: "0.01",
		FeesSupportScale: "1.3",
	}

	contract := openwallet.SmartContract{
		Address:  "0x35def49e4c26aadd2a0734f38f07dc6c7993f764",
		Symbol:   "LUK",
		Name:     "ZIK",
		Token:    "ZIK",
		Decimals: 18,
	}

	testGetAssetsAccountBalance(tm, walletID, accountID)

	testGetAssetsAccountTokenBalance(tm, walletID, accountID, contract)

	rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "", "", "",
		0, 100, &contract, &feesSupport)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return
	}

	//执行汇总交易
	for _, rawTxWithErr := range rawTxArray {

		if rawTxWithErr.Error != nil {
			log.Error(rawTxWithErr.Error.Error())
			continue
		}

		_, err = testSignTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}
	}

}