package common

import "time"

const FILE_SAVING_PATH = "./resource/"
const ENCRYPTION_FILE_PATH = FILE_SAVING_PATH+ "encryption/"
const DECRYPTION_FILE_PATH = FILE_SAVING_PATH+"public/"
const MARKET_PATH = FILE_SAVING_PATH+ "market/"
const USER_ICON_PATH = FILE_SAVING_PATH+"userIcon/"

// NFT TYPE
const (
	TYPE_NFT_AVATAR = "721-02"
	TYPE_NFT_MUSIC = "721-04"
	TYPE_NFT_OTHER = "721-05"
)

// NFT NAME
const (
	NAME_NFT_AVATAR = "avatar"
	NAME_NFT_MUSIC = "dat"
	NAME_NFT_OTHER = "other"
)

// ACTIVE_TICKER
const (
	ACTIVE_TICKER = "berry"
)

// base file path
const (
	BASE_FILE_PATH = "resource"
)

// purchase nft status
const (
	PURCHASE_CONFIRMED = 1
	PURCHASE_PENDING = 2
)
// path kind
const (
	PATH_KIND_MARKET = "market"
	PATH_KIND_PUBLIC = "public"
	PATH_KIND_ENCRYPT = "encrypt"
	PATH_KIND_DEFAULT = "default"
	PATH_KIND_USER_ICON = "userIcon"
)

// NFT transfer status
const (
	NFT_TRANSFER_SUCCESS = 0
	NFT_TRANSFER_PENDING = 1
)

const (
	MARKETPLACE_ID = "musicHotpot"
)

// const response status
const (
	RESPONSE_STATUS_SUCCESS = 0
	RESPONSE_STATUS_FAIL = 1
)

// berry purchase action status
const (
	BERRY_PURCHASE_PENDING = 2
	BERRY_PURCHASE_FINISH = 1
)

// shopping cart change
const (
	SHOPPING_CART_ADD = 0
	SHOPPING_CART_DELETE = 1
)

// follow list operation
const (
	FOLLOW_LIST_ADD = 0
	FOLLOW_LIST_DELETE = 1
)

// mysql extension error
const (
	DUPLICATE_ENTRY = "Error 1062"
)

type NftInfo struct {
	NftLdefIndex string `json:"nftLdefIndex"`
	NftType string	`json:"nftType"`
	NftName string	`json:"nftName"`
	ShortDescription string	`json:"shortDesc"`
	LongDescription string	`json:"longDesc"`
	FileName string	`json:"fileName"`
	NftParentLdef string	`json:"nftParentLdef"`
}

type AvatarNftInfo struct {
	NftInfo
	NftLifeIndex int	`json:"nftLifeIndex"`
	NftPowerIndex int	`json:"nftPowerIndex"`
}

type DatNftInfo struct {
	NftInfo
	MusicFileName string `json:"musicFileName"`
}

type OtherNftInfo struct {
	NftInfo
}

type NftMarketInfo struct {
	SellerWallet string	`json:"sellerWallet"`
	SellerUuid string	`json:"sellerUuid"`
	Price int	`json:"price"`
	Qty int	`json:"qty"`
	NumSold int	`json:"numSold"`
	Timestamp time.Time `json:"timestamp"`
}

type DatNftMarketInfo struct {
	DatNftInfo
	NftMarketInfo
	AllowAirdrop bool `json:"allowAirdrop"`
	CreatorPercent float64	`json:"creatorPercent"`
	LyricsWriterPercent float64	`json:"lyricsWriterPercent"`
	SongComposerPercent float64	`json:"songComposerPercent"`
	PublisherPercent float64	`json:"publisherPercent"`
	UserPercent float64	`json:"userPercent"`
}

type AvatarNftMarketInfo struct {
	AvatarNftInfo
	NftMarketInfo
}

type OtherNftMarketInfo struct {
	OtherNftInfo
	NftMarketInfo
}

type MarketPlaceInfo struct {
	MpId string `json:"mpId"`
	Active bool `json:"active"`
	ActiveTicker string `json:"active_ticker"`
}