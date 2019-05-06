package ws

const (
	GuestId = iota
	WeChatId
	FBId
	PhoneOrEmailId
)

// action status
const (
	ACTION_STATUS_FINISH = 1
	ACTION_STATUS_PENDING = 2
)

type ConnectErrorResponse struct {
	Reason string `json:"reason"`
}

type ErrorResponse struct {
	RQBaseInfo
	Reason string `json:"reason"`
}

type RQBaseInfo struct {
	Event string `json:"event"`
	Action string `json:"action"`
	ActId string `json:"actId"`
}

// MP_LIST
type MpListRequest struct {
	RQBaseInfo
	SupportedType string `json:"supportedType"`
	Page int `json:"page"`
	Offset int `json:"offset"`
}

type MpListNFTInfo struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName string `json:"nftName"`
	NftValue int `json:"nftValue" orm:"column(price)"`
	ActiveTicker string `json:"activeTicker"`
	NftLifeIndex int64 `json:"nftLifeIndex"`
	NftPowerIndex int64 `json:"nftPowerIndex"`
	NftLdefIndex string `json:"nftLdefIndex"`
	NftCharacId string `json:"nftCharacId"`
	Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
	Qty int `json:"qty"`
}

type MpListResponse struct {
	RQBaseInfo
	NftTranData []*MpListNFTInfo `json:"nftData"`
}

// PURCHASE_CONFIRM
type AsUserPurchaseConfirmInfo struct {
	AsId string `json:"asId"`
	AsWallet string `json:"asWallet"`
	Type int `json:"type"`
}

type PurchaseNftInfo struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName string `json:"nftName"`
	NftValue int `json:"nftValue" orm:"column(price)"`
	ActiveTicker string `json:"activeTicker"`
	NftLifeIndex int64 `json:"nftLifeIndex"`
	NftPowerIndex int64 `json:"nftPowerIndex"`
	NftLdefIndex string `json:"nftLdefIndex"`
	NftCharacId string `json:"nftCharacId"`
	Qty int `json:"qty"`
}

type PurchaseConfirmRequest struct {
	RQBaseInfo
	AsUser *AsUserPurchaseConfirmInfo `json:"asUser"`
	PackedKeys string `json:"packedKeys"`
	NftTranData []*PurchaseNftInfo `json:"nftTranData"`
}

type NftPurchaseResponseInfo struct {
	NftLdefIndex string `json:"nftLdefIndex"`
	Status int `json:"status"`
}

type PurchaseConfirmResponse struct {
	RQBaseInfo
	NftTranData []*NftPurchaseResponseInfo `json:"nftTranData"`
}

// TOKEN_PURCHASE
type AsUserPurchaseInfo struct {
	AsId string `json:"asId"`
	Type int `json:"type"`
	AsWallet string `json:"asWallet"`
	PackedKeys string `json:"packedKeys"`
}

type TokenPurchaseRequest struct {
	RQBaseInfo
	AsUser *AsUserPurchaseInfo `json:"asUser"`
	AppTranId string `json:"appTranId"`
	AppId string `json:"appId"`
	Amount int `json:"amount"`
	ActionStatus int `json:"actionStatus"`
}

type TokenPurchaseResponse struct {
	RQBaseInfo
	ActionStatus int `json:"actionStatus"`
}

//NFT_SHOW
type NftShowAsUserInfo struct {
	AsId string `json:"asId"`
	AsWallet string `json:"asWallet"`
	PackedKeys string `json:"packedKeys"`
}

type NftShowRequest struct {
	RQBaseInfo
	AsUser *NftShowAsUserInfo `json:"asUser"`
	TranAddress string `json:"tranAddress"`
	NftLdefIndex string `json:"nftLdefIndex"`
	SupportedType string `json:"supportedType"`
}

type NftShowResponse struct {
	RQBaseInfo
	NftLdefIndex string `json:"nftLdefIndex"`
	DecSource string `json:"decSource"`
}

// Item Details

type ItemDetailsRequestNftInfo struct {
	NftLdefIndex string `json:"nftLdefIndex"`
	SupportedType string `json:"supportedType"`
}

type ItemDetailsRequest struct {
	RQBaseInfo
	NftTranData []*ItemDetailsRequestNftInfo `json:"nftTranData"`
}

type ItemDetailsResponseNftInfo struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName string `json:"nftName"`
	NftValue int `json:"nftValue" orm:"column(price)"`
	ActiveTicker string `json:"activeTicker"`
	NftLifeIndex int64 `json:"nftLifeIndex"`
	NftPowerIndex int64 `json:"nftPowerIndex"`
	NftLdefIndex string `json:"nftLdefIndex"`
	NftCharacId string `json:"nftCharacId"`
	ShortDesc string `json:"shortDesc" orm:"column(short_description)"`
	LongDesc string `json:"longDesc" orm:"column(long_description)"`
	Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
	Qty int `json:"qty"`
}

type ItemDetailResponse struct {
	RQBaseInfo
	NftTranData []*ItemDetailsResponseNftInfo `json:"nftTranData"`
}

// TOTAL NFT

type AsUserNFTRequest struct {
	AsId string `json:"asId"`
	AsWallet string `json:"asWallet"`
}

type TotalNFTRequest struct {
	RQBaseInfo
	AsUser *AsUserNFTRequest `json:"asUser"`
}

type TotalNFTResponse struct {
	RQBaseInfo
	Count int `json:"count"`
}

// LIST_NFT

type AsUserListNFTRequest struct {
	AsId string `json:"asId"`
	AsWallet string `json:"asWallet"`
}

type ListNFTRequest struct {
	RQBaseInfo
	AsUser *AsUserListNFTRequest `json:""`
}

type NFTInfoListRes struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName string `json:"nftName"`
	NftValue int `json:"nftValue" orm:"column(price)"`
	ActiveTicker string `json:"activeTicker"`
	NftLifeIndex int64 `json:"nftLifeIndex"`
	NftPowerIndex int64 `json:"nftPowerIndex"`
	NftLdefIndex string `json:"nftLdefIndex"`
	NftCharacId string `json:"nftCharacId"`
	Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
	Qty int `json:"qty"`
}

type ListNFTResponse struct {
	RQBaseInfo
	NftTranData []*NFTInfoListRes `json:"nftTranData"`
}


// list of market user

type MarketUserListRequest struct {
	RQBaseInfo
}

type MarketUserWallet struct {
	WalletId string `json:"walletId" orm:"column(wallet_id)"`
}

type MarketUserListResponse struct {
	RQBaseInfo
	WalletIdList []*MarketUserWallet `json:"walletIdList"`
}

// user market info

type nftInfoListRes struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName string `json:"nftName"`
	NftValue int `json:"nftValue" orm:"column(price)"`
	ActiveTicker string `json:"activeTicker"`
	NftLifeIndex int64 `json:"nftLifeIndex"`
	NftPowerIndex int64 `json:"nftPowerIndex"`
	NftLdefIndex string `json:"nftLdefIndex"`
	NftCharacId string `json:"nftCharacId"`
	ShortDesc string `json:"shortDesc" orm:"column(short_description)"`
	LongDesc string `json:"longDesc" orm:"column(long_description)"`
	Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
	Qty int `json:"qty"`
}

type UserMarketInfoRequest struct {
	RQBaseInfo
	WalletId string `json:"walletId"`
}

type UserMarketInfoResponse struct {
	RQBaseInfo
	TotalNFT int `json:"totalNFT"`
	NftTranData []*nftInfoListRes `json:"nftTranData"`
}