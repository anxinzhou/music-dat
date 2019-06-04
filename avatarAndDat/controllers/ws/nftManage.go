package ws

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xxRanger/blockchainUtil/contract/nft"
	"github.com/xxRanger/music-dat/avatarAndDat/controllers/client"
	"github.com/xxRanger/music-dat/avatarAndDat/models"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path"
)

func (m *Manager) ItemDetailsHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req ItemDetailsRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	nftTranData := req.NftTranData
	nftResponseTranData := make([]*ItemDetailsResponseNftInfo, len(nftTranData))
	itemDetailRes := ItemDetailResponse{
		RQBaseInfo:  *bq,
		NftTranData: nftResponseTranData,
	}
	o := orm.NewOrm()
	for i, itemDetailsRequestNftInfo := range nftTranData {
		nftLdefIndex := itemDetailsRequestNftInfo.NftLdefIndex
		nftType := itemDetailsRequestNftInfo.SupportedType
		r := o.Raw(`
		select ni.nft_type, ni.nft_name,
		mk.price,mk.active_ticker, ni.nft_life_index, ni.nft_power_index, ni.nft_ldef_index,
		ni.nft_charac_id,  na.short_description, na.long_description ,mp.file_name, mp.icon_file_name,mk.qty from  
		nft_market_table as mk, 
		nft_mapping_table as mp,
		nft_info_table as ni,
		nft_item_admin as na 
		where mk.nft_ldef_index = mp.nft_ldef_index and mk.nft_ldef_index = ni.nft_ldef_index and  mp.nft_admin_id = na.nft_admin_id and  ni.nft_ldef_index = ? `, nftLdefIndex)
		var nftInfo NFTInfo
		err = r.QueryRow(&nftInfo)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		nftResInfo := &ItemDetailsResponseNftInfo{
			SupportedType: nftInfo.SupportedType,
			NftName:       nftInfo.NftName,
			NftValue:      nftInfo.NftValue,
			ActiveTicker:  nftInfo.ActiveTicker,
			NftLifeIndex:  nftInfo.NftLifeIndex,
			NftPowerIndex: nftInfo.NftPowerIndex,
			NftLdefIndex:  nftInfo.NftLdefIndex,
			NftCharacId:   nftInfo.NftCharacId,
			ShortDesc:     nftInfo.ShortDesc,
			LongDesc:      nftInfo.LongDesc,
			Thumbnail:     nftInfo.FileName,
			Qty:           nftInfo.Qty,
		}
		prefix := PathPrefixOfNFT(nftType, PATH_KIND_MARKET)
		if nftType == TYPE_NFT_AVATAR || nftType == TYPE_NFT_OTHER {
			nftResInfo.Thumbnail = prefix + nftInfo.FileName
		} else if nftType == TYPE_NFT_MUSIC {
			nftResInfo.Thumbnail = prefix + nftInfo.IconFileName
		} else {
			err := errors.New("unexpected type")
			logs.Emergency(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		nftResponseTranData[i] = nftResInfo
	}
	m.wrapperAndSend(c, bq, itemDetailRes)
}

// nft show
func (m *Manager) NFTDisplayHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req NftShowRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	// check nft type
	if err := validSupportedType(req.SupportedType); err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	nftLdefIndex := req.NftLdefIndex
	// check nft ldefindex
	if err := validNftLdefIndex(req.NftLdefIndex); err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	mp := models.NftMappingTable{
		NftLdefIndex: nftLdefIndex,
	}
	o := orm.NewOrm()
	err = o.Read(&mp)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	fileName := mp.FileName
	//TODO user symmetric key from client to decrypt file
	var encryptedFilePath string
	var decryptedFilePath string
	logs.Debug("nft type from request,", req.SupportedType)
	if req.SupportedType == TYPE_NFT_AVATAR {
		encryptedFilePath = path.Join(ENCRYPTION_FILE_PATH, NAME_NFT_AVATAR, fileName)
		decryptedFilePath = path.Join(DECRYPTION_FILE_PATH, NAME_NFT_AVATAR, fileName)
	} else if req.SupportedType == TYPE_NFT_MUSIC {
		encryptedFilePath = path.Join(ENCRYPTION_FILE_PATH, NAME_NFT_MUSIC, fileName)
		decryptedFilePath = path.Join(DECRYPTION_FILE_PATH, NAME_NFT_MUSIC, fileName)
	} else if req.SupportedType == TYPE_NFT_OTHER {
		encryptedFilePath = path.Join(ENCRYPTION_FILE_PATH, NAME_NFT_OTHER, fileName)
		decryptedFilePath = path.Join(DECRYPTION_FILE_PATH, NAME_NFT_OTHER, fileName)
	}
	cipherText, err := ioutil.ReadFile(encryptedFilePath)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	nonce, ct := cipherText[:aesgcm.NonceSize()], cipherText[aesgcm.NonceSize():]
	originalData, err := aesgcm.Open(nil, nonce, ct, nil)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	logs.Debug("length of original data", len(originalData))
	if req.SupportedType == TYPE_NFT_AVATAR || req.SupportedType == TYPE_NFT_OTHER {
		out, err := os.Create(decryptedFilePath)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		defer out.Close()
		originalImage, _, err := image.Decode(bytes.NewBuffer(originalData))
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
		err = jpeg.Encode(out, originalImage, nil)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
	} else if req.SupportedType == TYPE_NFT_MUSIC {
		err := ioutil.WriteFile(decryptedFilePath, originalData, 0777)
		if err != nil {
			logs.Error(err.Error())
			m.errorHandler(c, bq, err)
			return
		}
	}

	m.wrapperAndSend(c, bq, NftShowResponse{
		RQBaseInfo:   *bq,
		NftLdefIndex: nftLdefIndex,
		DecSource: beego.AppConfig.String("prefix") + beego.AppConfig.String("hostaddr") + ":" +
			beego.AppConfig.String("fileport") + "/" + decryptedFilePath,
	})
}

func (m *Manager) NFTTransferHandler(c *client.Client, bq *RQBaseInfo, data []byte) {
	var req NftTransferRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	senderNickname := req.SenderNickname
	receiverNickname := req.ReceiverNickname
	nftInfo := req.NftTranData
	nftLdefIndex := nftInfo.NftLdefIndex
	supportedType := nftInfo.SupportedType
	if err := validNftLdefIndex(nftLdefIndex); err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	if err := validSupportedType(supportedType); err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}

	senderWalletId, err := models.WalletIdOfNickname(senderNickname)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	receiverWalletId, err := models.WalletIdOfNickname(receiverNickname)
	if err != nil {
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	o:=orm.NewOrm()
	o.Begin()
	// add count for sender
	_,err = o.QueryTable("market_user_table").Filter("nickname",senderNickname).Update(orm.Params{
		"count": orm.ColValue(orm.ColAdd,1),
	})
	if err!=nil {
		o.Rollback()
		logs.Emergency("can not add count for nickname:", senderNickname)
		m.errorHandler(c, bq, err)
		return
	}
	logs.Warn("add count in market table for",senderNickname)
	// reduce count for receiver
	_,err = o.QueryTable("market_user_table").Filter("nickname",receiverNickname).Update(orm.Params{
		"count": orm.ColValue(orm.ColMinus,1),
	})
	if err!=nil {
		o.Rollback()
		logs.Emergency("can not add count for nickname:", receiverNickname)
		m.errorHandler(c, bq, err)
		return
	}
	logs.Warn("reduce count in market table for",receiverNickname)


	tokenId := TokenIdFromNftLdefIndex(nftLdefIndex)
	txErr := m.chainHandler.ManagerAccount.SendFunction(m.chainHandler.Contract,
		nil,
		nft.FuncDelegateTransfer,
		common.HexToAddress(senderWalletId),
		common.HexToAddress(receiverWalletId),
		tokenId)
	err = <-txErr
	if err != nil {
		o.Rollback()
		logs.Error(err.Error())
		m.errorHandler(c, bq, err)
		return
	}
	o.Commit()
	logs.Debug("send", nftLdefIndex, "from", senderWalletId, "to", receiverWalletId)
	m.wrapperAndSend(c, bq, &NftTransferResponse{
		RQBaseInfo: *bq,
		Status:     NFT_TRANSFER_SUCCESS,
	})
}
