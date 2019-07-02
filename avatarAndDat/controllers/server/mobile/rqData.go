package mobile
//
//const (
//	GuestId = iota
//	WeChatId
//	FBId
//	PhoneOrEmailId
//)
//
//// action status
//const (
//	ACTION_STATUS_FINISH = 1
//	ACTION_STATUS_PENDING = 2
//)
//
//// shopping cart operation
//const (
//	SHOPPING_CART_ADD = 0
//	SHOPPING_CART_DELETE = 1
//)
//
//// follow list operation
//const (
//	FOLLOW_LIST_ADD = 0
//	FOLLOW_LIST_DELETE = 1
//)
//
//type ConnectErrorResponse struct {
//	Reason string `json:"reason"`
//}
//
//type ErrorResponse struct {
//	RQBaseInfo
//	Reason string `json:"reason"`
//}
//
//type RQBaseInfo struct {
//	Event string `json:"event"`
//	Action string `json:"action"`
//	ActId string `json:"actId"`
//}
//
//type NftBaseInfo struct {
//	NftLdefIndex string `json:"nftLdefIndex"`
//	SupportedType string `json:"supportedType"`
//}
//
//// MP_LIST
//type MpListRequest struct {
//	RQBaseInfo
//	SupportedType string `json:"supportedType"`
//	Page int `json:"page"`
//	Offset int `json:"offset"`
//}
//
//type MpListNFTInfo struct {
//	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
//	NftName string `json:"nftName"`
//	NftValue int `json:"nftValue" orm:"column(price)"`
//	ActiveTicker string `json:"activeTicker"`
//	NftLifeIndex int64 `json:"nftLifeIndex"`
//	NftPowerIndex int64 `json:"nftPowerIndex"`
//	NftLdefIndex string `json:"nftLdefIndex"`
//	NftCharacId string `json:"nftCharacId"`
//	Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
//	Qty int `json:"qty"`
//}
//
//type MpListResponse struct {
//	RQBaseInfo
//	NftTranData []*MpListNFTInfo `json:"nftData"`
//}
//
//// PURCHASE_CONFIRM
//
//type PurchaseConfirmRequest struct {
//	RQBaseInfo
//	Uuid string `json:"uuid"`
//	NftTranData []*NftBaseInfo `json:"nftTranData"`
//}
//
//type NftPurchaseResponseInfo struct {
//	NftLdefIndex string `json:"nftLdefIndex"`
//	Status int `json:"status"`
//}
//
//type PurchaseConfirmResponse struct {
//	RQBaseInfo
//	NftTranData []*NftPurchaseResponseInfo `json:"nftTranData"`
//}
//
//// TOKEN_PURCHASE
//type AsUserPurchaseInfo struct {
//	AsId string `json:"asId"`
//	Type int `json:"type"`
//	AsWallet string `json:"asWallet"`
//	PackedKeys string `json:"packedKeys"`
//}
//
//type TokenPurchaseRequest struct {
//	RQBaseInfo
//	Nickname string `json:"nickname"`
//	AppTranId string `json:"appTranId"`
//	TransactionId string `json:"transactionId"`
//	AppId string `json:"appId"`
//	Amount int `json:"amount"`
//	ActionStatus int `json:"actionStatus"`
//}
//
//type TokenPurchaseResponse struct {
//	RQBaseInfo
//	ActionStatus int `json:"actionStatus"`
//	TransactionId string `json:"transactionId"`
//}
//
////NFT_SHOW
//type NftShowRequest struct {
//	RQBaseInfo
//	NftLdefIndex string `json:"nftLdefIndex"`
//	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
//}
//
//type NftShowResponse struct {
//	RQBaseInfo
//	NftLdefIndex string `json:"nftLdefIndex"`
//	DecSource string `json:"decSource"`
//}
//
//// Item Details
//
//type ItemDetailsRequestNftInfo struct {
//	NftLdefIndex string `json:"nftLdefIndex"`
//	SupportedType string `json:"supportedType"`
//}
//
//type ItemDetailsRequest struct {
//	RQBaseInfo
//	NftTranData []*ItemDetailsRequestNftInfo `json:"nftTranData"`
//}
//
//type ItemDetailsResponseNftInfo struct {
//	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
//	NftName string `json:"nftName"`
//	NftValue int `json:"nftValue" orm:"column(price)"`
//	ActiveTicker string `json:"activeTicker"`
//	NftLifeIndex int64 `json:"nftLifeIndex"`
//	NftPowerIndex int64 `json:"nftPowerIndex"`
//	NftLdefIndex string `json:"nftLdefIndex"`
//	NftCharacId string `json:"nftCharacId"`
//	ShortDesc string `json:"shortDesc" orm:"column(short_description)"`
//	LongDesc string `json:"longDesc" orm:"column(long_description)"`
//	Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
//	Qty int `json:"qty"`
//}
//
//type ItemDetailResponse struct {
//	RQBaseInfo
//	NftTranData []*ItemDetailsResponseNftInfo `json:"nftTranData"`
//}
//
//// TOTAL NFT
//
//type AsUserNFTRequest struct {
//	AsId string `json:"asId"`
//	AsWallet string `json:"asWallet"`
//}
//
//type TotalNFTRequest struct {
//	RQBaseInfo
//	AsUser *AsUserNFTRequest `json:"asUser"`
//}
//
//type TotalNFTResponse struct {
//	RQBaseInfo
//	Count int `json:"count"`
//}
//
//// LIST_NFT
//
//type AsUserListNFTRequest struct {
//	AsId string `json:"asId"`
//	AsWallet string `json:"asWallet"`
//}
//
//type ListNFTRequest struct {
//	RQBaseInfo
//	AsUser *AsUserListNFTRequest `json:""`
//}
//
//type NFTInfoListRes struct {
//	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
//	NftName string `json:"nftName"`
//	NftValue int `json:"nftValue" orm:"column(price)"`
//	ActiveTicker string `json:"activeTicker"`
//	NftLifeIndex int64 `json:"nftLifeIndex"`
//	NftPowerIndex int64 `json:"nftPowerIndex"`
//	NftLdefIndex string `json:"nftLdefIndex"`
//	NftCharacId string `json:"nftCharacId"`
//	Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
//	Qty int `json:"qty"`
//}
//
//type ListNFTResponse struct {
//	RQBaseInfo
//	NftTranData []*NFTInfoListRes `json:"nftTranData"`
//}
//
//
//// list of market user
//
//type MarketUserListRequest struct {
//	RQBaseInfo
//	Page int `json:"page"` //TODO
//	Offet int `json:"offset"`  //TODO
//	Nickname string `json:"nickname"`
//}
//
//type MarketUserWallet struct {
//	WalletId string `json:"walletId" orm:"column(wallet_id)"`
//	Nickname string `json:"nickname"`
//	Count int `json:"count"`
//	Intro string `json:"intro"`
//	Thumbnail string `json:"thumbnail" orm:"column(user_icon_url)"`
//	Followed bool `json:"followed"`
//}
//
//type MarketUserListResponse struct {
//	RQBaseInfo
//	WalletIdList []*MarketUserWallet `json:"walletIdList"`
//}
//
//// user market info
//
//type nftInfoListRes struct {
//	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
//	NftName string `json:"nftName"`
//	NftValue int `json:"nftValue" orm:"column(price)"`
//	ActiveTicker string `json:"activeTicker"`
//	NftLifeIndex int64 `json:"nftLifeIndex"`
//	NftPowerIndex int64 `json:"nftPowerIndex"`
//	NftLdefIndex string `json:"nftLdefIndex"`
//	NftCharacId string `json:"nftCharacId"`
//	ShortDesc string `json:"shortDesc" orm:"column(short_description)"`
//	LongDesc string `json:"longDesc" orm:"column(long_description)"`
//	Thumbnail string `json:"thumbnail" orm:"column(file_name)"`
//	DecSource string `json:"decSource"`
//	Qty int `json:"qty"`
//}
//
//type UserMarketInfoRequest struct {
//	RQBaseInfo
//	Nickname string `json:"nickname"`
//}
//
//type UserMarketInfoResponse struct {
//	RQBaseInfo
//	TotalNFT int `json:"totalNFT"`
//	NftTranData []*nftInfoListRes `json:"nftTranData"`
//}
//
//type NFTInfo struct {
//	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
//	NftName string `json:"nftName"`
//	NftValue int `json:"nftValue" orm:"column(price)"`
//	ActiveTicker string `json:"activeTicker"`
//	NftLifeIndex int64 `json:"nftLifeIndex"`
//	NftPowerIndex int64 `json:"nftPowerIndex"`
//	NftLdefIndex string `json:"nftLdefIndex"`
//	NftCharacId string `json:"nftCharacId"`
//	ShortDesc string `json:"shortDesc" orm:"column(short_description)"`
//	LongDesc string `json:"longDesc" orm:"column(long_description)"`
//	FileName string `json:"fileName" orm:"column(file_name)"`
//	IconFileName string `json:"iconFileName" orm:"column(icon_file_name)"`
//	Qty int `json:"qty"`
//}
//
//type NFTPurchaseHistoryRequest struct {
//	RQBaseInfo
//	Nickname string `json:"nickname"`
//	SupportedType string `json:"supportedType"`
//}
//
//type NFTPurchaseRecord struct {
//	PurchaseId string `json:"purchaseId"`
//	TransactionAddress string `json:"transactionAddress"`
//	WalletId string `json:"walletId"`
//	NftTranData *nftInfoListRes `json:"nftTranData"`
//	Timestamp string `json:"timeStamp"`
//	Status int `json:"status"`
//}
//
//type NFTPurchaseHistoryResponse struct {
//	RQBaseInfo
//	PurchaseList []*NFTPurchaseRecord `json:"purchaseList"`
//}
//
//type ShoppingCartChangeRequest struct {
//	RQBaseInfo
//	Nickname string `json:"nickname"`
//	Operation int `json:"operation"` // 0 for add 1 for delete
//	NFTList []string `json:"nftList"`
//}
//
//type ShoppingCartChangeResponse struct {
//	RQBaseInfo
//}
//
//type ShoppingCartListRequest struct {
//	RQBaseInfo
//	Nickname string `json:"Nickname"`
//}
//
//type ShoppingCartRecord struct {
//	NftTranData *nftInfoListRes `json:"nftTranData"`
//	Timestamp string `json:"timestamp"`
//}
//
//type ShoppingCartListResponse struct {
//	RQBaseInfo
//	NftList []*ShoppingCartRecord `json:"nftList"`
//}
//
//// todo
//// nft transfer
//type NftTransferRequest struct {
//	RQBaseInfo
//	SenderNickname string `json:"senderNickname"`
//	ReceiverNickname string `json:"receiverNickname"`
//	NftTranData *NftBaseInfo `json:"nftTranData"`
//}
//
//type NftTransferResponse struct {
//	RQBaseInfo
//	Status int `json:"status"`
//}
//
//// bind wallet
//type BindWalletRequest struct {
//	RQBaseInfo
//	Nickname string `json:"nickname"`
//	WalletId string `json:"walletId"`
//}
//
//type BindWalletResponse struct {
//	RQBaseInfo
//}
//
//// set nick name
//type SetNicknameRequest struct {
//	RQBaseInfo
//	Uuid string `json:"uuid"`
//	Nickname string `json:"nickname"`
//}
//
//type SetNicknameResponse struct {
//	RQBaseInfo
//}
//
//// is nick name duplicated
//type IsNicknameDuplicatedRequest struct {
//	RQBaseInfo
//	Nickname string `json:"nickname"`
//}
//
//type IsNicknameDuplicatedResponse struct {
//	RQBaseInfo
//	Duplicated bool `json:"duplicated"`
//}
//
//type FollowInfo struct {
//	Nickname string `json:"nickname" orm:"column(nickname)"`
//	Thumbnail string `json:"thumbnail" orm:"user_icon_url"`
//	Intro string `json:"intro"`
//}
//
//// follow list
//type FollowListRequest struct {
//	RQBaseInfo
//	Nickname string `json:"nickname" `
//}
//
//type FollowListResponse struct {
//	RQBaseInfo
//	FollowInfo []*FollowInfo `json:"followInfo"`
//}
//
//// follow list operation
//type FollowListOperationRequest struct {
//	RQBaseInfo
//	Nickname string `json:"nickname"`
//	FollowNickname string `json:"followNickname"`
//	Operation int `json:"operation"`
//}
//
//type FollowListOperationResponse struct {
//	RQBaseInfo
//	FollowNickname string `json:"followNickname"`
//	Operation int `json:"operation"`
//}
//
//// has user set nickname
//type IsNicknameSetRequest struct {
//	RQBaseInfo
//	Uuid string `json:"uuid"`
//}
//
//type IsNicknameSetResponse struct {
//	RQBaseInfo
//	Set bool `json:"set"`
//}