pragma solidity ^0.5.2;

contract DAT {

    // music dat
    struct NFT {
      string nftType;
      string nftName;
      string nftLdefIndex;
      string distIndex;
      uint256 nftLifeIndex;
      uint256 nftPowerIndex;
      string nftCharacterId;
      bytes publicKey;
    }

    mapping (uint256 => address) private _tokenOwner;
    mapping (uint256 => NFT) private _tokenInfo;
    mapping (address => uint256) private _ownedTokensCount;

    constructor () public {

    }

    function balanceOf(address owner) public view returns (uint256) {
        require(owner != address(0));
        return _ownedTokensCount[owner];
    }

    function ownerOf(uint256 tokenId) public view returns (address) {
        address owner = _tokenOwner[tokenId];
        require(owner != address(0));
        return owner;
    }

    function mint(
        address to,
        uint256 token_id,
        string memory nft_type,
        string memory nft_name,
        string memory nft_ldef_index,
        string memory dist_index,
        uint256 nft_life_index,
        uint256 nft_power_index,
        string memory nft_character_id,
        bytes memory public_key)  public {
        require(to != address(0));
        require(!_exists(token_id));
        NFT memory m = NFT(
            {
                nftType: nft_type,
                nftName: nft_name,
                nftLdefIndex: nft_ldef_index,
                nftLifeIndex: nft_life_index,
                distIndex:dist_index,
                nftPowerIndex: nft_power_index,
                nftCharacterId: nft_character_id,
                publicKey: public_key
            }
        );
        _tokenInfo[token_id] = m;
        _tokenOwner[token_id] = to;
        _ownedTokensCount[to] += 1 ;
    }

    function infoOfToken(uint256 token_id) external view returns (
        string memory nft_type,
        string memory nft_name,
        string memory nft_ldef_index,
        uint256 nft_life_index,
        uint256 nft_power_index,
        string memory nft_character_id,
        bytes  memory public_key){
        require(_exists(token_id));
        NFT memory nft = _tokenInfo[token_id];
        return (nft.nftType,nft.nftName,nft.nftLdefIndex,nft.nftLifeIndex,nft.nftPowerIndex,nft.nftCharacterId,nft.publicKey);
    }

    function ldefOfToken(uint256 token_id) external view returns (string memory nft_ldef_index) {
        require(_exists(token_id));
        return _tokenInfo[token_id].nftLdefIndex;
    }

    function transfer(address to, uint256 token_id) public {
        address from = msg.sender;
        require(ownerOf(token_id) == from);
        require(to != address(0));

        _ownedTokensCount[from]-=1;
        _ownedTokensCount[to]+=1;
        _tokenOwner[token_id] = to;
    }

    function _exists(uint256 tokenId) internal view returns (bool) {
        address owner = _tokenOwner[tokenId];
        return owner != address(0);
    }
}