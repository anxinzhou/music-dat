package common

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

// login type
const (
	LOGIN_TYPE_USERNAME = 3
)
