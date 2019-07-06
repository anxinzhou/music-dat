package transactionQueue

import (
	"encoding/hex"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/go-ethereum/common"
	common2 "github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/chainHelper"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/server/common/util"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"math/big"
	"time"
)

type NftPurchaseTransaction struct {
	Uuid string
	SellerUuid string
	NftLdefIndex string
	PurchaseId string
}

type UploadNftTransaction struct {
	Uuid string
	NftLdefIndex string
	NftType string
	NftName string
	DistIndex string
	NftLifeIndex *big.Int
	NftPowerIndex *big.Int
	NftCharacterId string
	PublicKey string
}

type RewardNftTransaction struct {
	Uuid string
	SellerUuid string
	NftLdefIndex string
	PurchaseId string
}

type TransferNftTransaction struct {
	Uuid string
	SellerUuid string
	NftLdefIndex string
}

type TransactionQueue struct {
	TransactionPool chan interface{}
	ChainHandler *chainHelper.ChainHandler
}

func NewTransactionQueue(transactionPool chan interface{}, chainHandler *chainHelper.ChainHandler) *TransactionQueue {
	return &TransactionQueue{
		TransactionPool:transactionPool,
		ChainHandler:  chainHandler,
	}
}

func (this *TransactionQueue) Start() {
	for transaction:= range this.TransactionPool {
		switch tx:=transaction.(type) {
		case *NftPurchaseTransaction:
			go this.SendNftPurchaseTransaction(tx)
		case *UploadNftTransaction:
			go this.SendUploadNftTransaction(tx)
		case *RewardNftTransaction:
			go this.SendRewardNftTransaction(tx)
		case *TransferNftTransaction:
			go this.SendTransferNftTransaction(tx)
		default:
			panic("no such kind of transaction")
		}
	}
}

func (this *TransactionQueue) Append(transaction interface{}) {
	this.TransactionPool<-transaction
}

func (this *TransactionQueue) Retry(transaction interface{}) {
	<-time.After(1*time.Second)
	this.Append(transaction)
}

func (this *TransactionQueue) SendNftPurchaseTransaction(nftPurchaseInfo *NftPurchaseTransaction) {
	// check if nft is active
	nftMarketInfo:= models.NftMarketInfo{
		NftLdefIndex: nftPurchaseInfo.NftLdefIndex,
	}
	o:= orm.NewOrm()
	err:=o.Read(&nftMarketInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Retry(nftPurchaseInfo)
		return
	}
	if nftMarketInfo.Active == false {
		err:=errors.New(nftPurchaseInfo.NftLdefIndex+" is not active now, retry after seconds")
		logs.Error(err.Error())
		this.Retry(nftPurchaseInfo)
		return
	}

	buyerMarketInfo:= models.UserMarketInfo{
		Uuid: nftPurchaseInfo.Uuid,
	}
	err = o.Read(&buyerMarketInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Retry(nftPurchaseInfo)
		return
	}
	sellerMarketInfo:= models.UserMarketInfo{
		Uuid: nftPurchaseInfo.SellerUuid,
	}
	o.Read(&sellerMarketInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Retry(nftPurchaseInfo)
		return
	}
	tokenId,err:= util.TokenIdFromNftLdefIndex(nftPurchaseInfo.NftLdefIndex)
	if err!=nil {
		panic(err.Error())
	}
	txErr := this.ChainHandler.ManagerAccount.SendFunction(
		this.ChainHandler.Contract,
		nil,
		nft.FuncDelegateTransfer,
		common.HexToAddress(sellerMarketInfo.Wallet),
		common.HexToAddress(buyerMarketInfo.Wallet),
		tokenId,
	)
	err = <-txErr
	if err != nil {
		logs.Error(err.Error())
		this.Retry(nftPurchaseInfo)
		return
	}
	// change status in nft purchaseInfo table
	userPurchaseInfo:= models.NftPurchaseInfo{
		PurchaseId: nftPurchaseInfo.PurchaseId,
		Status: common2.PURCHASE_CONFIRMED,
	}
	_,err=o.Update(&userPurchaseInfo,"status")   //TODO fault recovery process
	if err!=nil {
		logs.Emergency(err.Error())
	}
	logs.Info("done with a nft purchase transaction")
}

func (this *TransactionQueue) SendUploadNftTransaction(uploadNftInfo *UploadNftTransaction) {
	//logs.Info("upload nft, nftLdefindex",uploadNftInfo.NftLdefIndex,"uuid",uploadNftInfo.Uuid)
	userMarkerInfo:= models.UserMarketInfo{
		Uuid: uploadNftInfo.Uuid,
	}
	o:=orm.NewOrm()
	err:=o.Read(&userMarkerInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Append(uploadNftInfo)
		return
	}

	// send transaction
	// prepare parameters for transaction
	tokenId,err:= util.TokenIdFromNftLdefIndex(uploadNftInfo.NftLdefIndex)
	if err!=nil {
		panic(err)
	}
	wallet:= common.HexToAddress(userMarkerInfo.Wallet)
	nftType:= uploadNftInfo.NftType
	nftLdefIndex:= uploadNftInfo.NftLdefIndex
	distIndex:= uploadNftInfo.DistIndex
	nftCharacterId:= uploadNftInfo.NftCharacterId
	publicKey,err:= hex.DecodeString(uploadNftInfo.PublicKey)
	if err!=nil {
		panic(err)
	}
	// follow
	var nftName string
	nftLifeIndex:= uploadNftInfo.NftLifeIndex
	nftPowerIndex:= uploadNftInfo.NftPowerIndex
	switch nftType {
	case common2.TYPE_NFT_AVATAR:
		var avatarInfo models.AvatarNftInfo
		err:=o.QueryTable("avatar_nft_info").RelatedSel("NftInfo").One(&avatarInfo)
		if err!=nil {
			logs.Error(err.Error())
			this.Retry(uploadNftInfo)
			return
		}
		nftName = avatarInfo.NftInfo.NftName

	case common2.TYPE_NFT_OTHER:
		var otherInfo  models.OtherNftInfo
		err:=o.QueryTable("other_nft_info").RelatedSel("NftInfo").One(&otherInfo)
		if err!=nil {
			logs.Error(err.Error())
			this.Retry(uploadNftInfo)
			return
		}
		nftName = otherInfo.NftInfo.NftName
	case common2.TYPE_NFT_MUSIC:
		var datInfo models.DatNftInfo
		err:=o.QueryTable("dat_nft_info").RelatedSel("NftInfo").One(&datInfo)
		if err!=nil {
			logs.Error(err.Error())
			this.Retry(uploadNftInfo)
			return
		}
		nftName = datInfo.NftInfo.NftName
	}
	txErr:=this.ChainHandler.ManagerAccount.SendFunction(
		this.ChainHandler.Contract,
		nil,
		nft.FuncMint,
		wallet,
		tokenId,
		nftType,
		nftName,
		nftLdefIndex,
		distIndex,
		nftLifeIndex,
		nftPowerIndex,
		nftCharacterId,
		publicKey,
		)
	err = <-txErr
	if err!=nil {
		logs.Error(err.Error())
		this.Retry(uploadNftInfo)
		return
	}
	// send transaction success, modify tag in marketplace table
	nftMarketInfo:= models.NftMarketInfo{
		NftLdefIndex: nftLdefIndex,
		Active: true,
	}
	_,err=o.Update(&nftMarketInfo,"active")
	if err!=nil {
		logs.Emergency(err.Error())  //TODO fault recovery process
	}
	logs.Info("done with a upload transaction")
}

func (this *TransactionQueue) SendRewardNftTransaction(RewardNftInfo *RewardNftTransaction) {
	// check if nft is active
	nftMarketInfo:= models.NftMarketInfo{
		NftLdefIndex: RewardNftInfo.NftLdefIndex,
	}
	o:= orm.NewOrm()
	err:=o.Read(&nftMarketInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Retry(RewardNftInfo)
		return
	}
	if nftMarketInfo.Active == false {
		err:=errors.New(RewardNftInfo.NftLdefIndex+" is not active now, retry after seconds")
		logs.Error(err.Error())
		this.Retry(RewardNftInfo)
		return
	}

	buyerMarketInfo:= models.UserMarketInfo{
		Uuid: RewardNftInfo.Uuid,
	}
	err = o.Read(&buyerMarketInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Retry(RewardNftInfo)
		return
	}
	sellerMarketInfo:= models.UserMarketInfo{
		Uuid: RewardNftInfo.SellerUuid,
	}
	o.Read(&sellerMarketInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Retry(RewardNftInfo)
		return
	}
	tokenId,err:= util.TokenIdFromNftLdefIndex(RewardNftInfo.NftLdefIndex)
	if err!=nil {
		panic(err)
	}
	txErr := this.ChainHandler.ManagerAccount.SendFunction(
		this.ChainHandler.Contract,
		nil,
		nft.FuncDelegateTransfer,
		common.HexToAddress(sellerMarketInfo.Wallet),
		common.HexToAddress(buyerMarketInfo.Wallet),
		tokenId,
	)
	err = <-txErr
	if err != nil {
		logs.Error(err.Error())
		this.Retry(RewardNftInfo)
		return
	}
	// change status in nft purchaseInfo table
	userPurchaseInfo:= models.NftPurchaseInfo{
		PurchaseId: RewardNftInfo.PurchaseId,
		Status: common2.PURCHASE_CONFIRMED,
	}
	_,err=o.Update(&userPurchaseInfo,"status")   //TODO fault recovery process
	if err!=nil {
		logs.Emergency(err.Error())
	}
	logs.Info("done with a reward nft transaction")
}

func (this *TransactionQueue) SendTransferNftTransaction(transferNftInfo *TransferNftTransaction) {
	// check if nft is active
	nftMarketInfo:= models.NftMarketInfo{
		NftLdefIndex: transferNftInfo.NftLdefIndex,
	}
	o:= orm.NewOrm()
	err:=o.Read(&nftMarketInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Retry(transferNftInfo)
		return
	}
	if nftMarketInfo.Active == false {
		err:=errors.New(transferNftInfo.NftLdefIndex+" is not active now, retry after seconds")
		logs.Error(err.Error())
		this.Retry(transferNftInfo)
		return
	}
	buyerMarketInfo:= models.UserMarketInfo{
		Uuid: transferNftInfo.Uuid,
	}
	err = o.Read(&buyerMarketInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Retry(transferNftInfo)
		return
	}
	sellerMarketInfo:= models.UserMarketInfo{
		Uuid: transferNftInfo.SellerUuid,
	}
	o.Read(&sellerMarketInfo)
	if err!=nil {
		logs.Error(err.Error())
		this.Retry(transferNftInfo)
		return
	}

	tokenId,err:= util.TokenIdFromNftLdefIndex(transferNftInfo.NftLdefIndex)
	if err!=nil {
		panic(err)
	}
	txErr:= this.ChainHandler.ManagerAccount.SendFunction(
		this.ChainHandler.Contract,
		nil,
		nft.FuncDelegateTransfer,
		common.HexToAddress(sellerMarketInfo.Wallet),
		common.HexToAddress(buyerMarketInfo.Wallet),
		tokenId,
	)
	err = <-txErr
	if err != nil {
		logs.Error(err.Error())
		this.Retry(transferNftInfo)
		return
	}
	logs.Info("done with a transfer nft transaction")
}

