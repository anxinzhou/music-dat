package ws

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

type MpListRequest struct {
	RQBaseInfo
	SupportedType string `json:"supportedType"`
	Page int `json:"page"`
	Offset int `json:"offset"`
}

type NFTInfo struct {
	SupportedType string `json:"supportedType" orm:"column(nft_type)"`
	NftName string `json:"nftName"`
	NftValue float64 `json:"nftValue" orm:"column(price)"`
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
	NftData []*NFTInfo `json:"nftData"`
}