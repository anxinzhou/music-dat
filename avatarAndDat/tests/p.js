const WebSocket = require('ws');

// const ws = new WebSocket('ws://3.1.35.36:4000/ws');
const ws = new WebSocket('ws://localhost:4000/ws');
// const ws = new WebSocket("ws://18.136.134.84:4000");
// const now = require('performance-now');

let testNftLdefIndex = "M982980481";
let tokenId = "396285078";
//done
var mkinfos = {
    action: "mp_list",
    event: "market",
    actId: "A1201212",
    supportedType: "721-04"
}

var tokenPurchase = {
    "event": "token_purchase", //string
    "action": "tokenbuy_paid", //string
    "actId": "APP01234776789", //string
    "asUser": {
        "asId": "107255250381177", // string
        "type": 1,
        "asWallet": "4fffc3b55bd8270df6fd85d1f17b8ec53adf13d4552c2079de0f0bceea64d299", //string
        "packedKeys": "U1074972YHBFGHJHDOJDJKDKDKKDLLD5462", //string
    },
    "appTranId": "4fffc3b54bd8270df6fd85d1f17b8ec53adfa64d882",
    "appId": "XXX@appleid.com",
    "amount": 1, //int    berry number
    "actionStatus": 1
}

var nftshow = {
    "event": "nft_show", //string
    "action": "NFT_display", //string
    "actId": "APP01234776789", //string
    "asUser": {
        "asId": "AS4789", //string
        "asWallet": "4fffc3b55bd8270df6fd85d1f17b8ec53adf13d4552c2079de0f0bceea64d299", //string
        "packedKeys": "U1074972YHBFGHJHDOJDJKDKDKKDLLD5462", //string
    },
    "tranAddress": "4fffc3b55bd8270df6fd85d1f17b8ec53adf13d4552c2079de0f0bceea64d882", //string
    "nftLdefIndex": testNftLdefIndex + "1", //string  // nft to show
    "supportedType": "721-02", // string "721-04" for music "721-02" for avatar
}

var musicshow = {
    "event": "nft_show", //string
    "action": "NFT_display", //string
    "actId": "APP01234776789", //string
    "asUser": {
        "asId": "AS4789", //string
        "asWallet": "4fffc3b55bd8270df6fd85d1f17b8ec53adf13d4552c2079de0f0bceea64d299", //string
        "packedKeys": "U1074972YHBFGHJHDOJDJKDKDKKDLLD5462", //string
    },
    "tranAddress": "4fffc3b55bd8270df6fd85d1f17b8ec53adf13d4552c2079de0f0bceea64d882", //string
    "nftLdefIndex": testNftLdefIndex, //string  // nft to show
    "supportedType": "721-04", // string "721-04" for music "721-02" for avatar
}

var itemDetails = {
    "event": "nft_market", //string
    "action": "item_details", //string
    "actId": "APP0123456889", //string
    "nftTranData": [ // array of object
        {
            "nftLdefIndex": testNftLdefIndex,
            "supportedType": "721-04", // string "721-04" for music "721-02" for avatar
        },
    ]
}

var purchaseConfirm = {
    "event": "nft_market", // string
    "action": "NFT_purchase_confirm", // string
    "actId": "APP0123456789", //string
    "asUser": { // object
        "asId": "Simonchen", // string
        "type": 3, // int  1,2 for login by wechat or facebook  3 for phone or email
        "asWallet": "0x3c62aa7913bc303ee4b9c07df87b556b6770e3fc", //string
    },
    "packedKeys": "U1074972YHBFGHJHDOJDJKDKDKKDLLD5462",
    "nftTranData": [ //array of object
        {
            "supportedType": "721-04", // string "721-04" for music "721-02" for avatar
            "nftName": "heart", // string
            "nftValue": 1, // int
            "activeTicker": "berry", //string
            "nftLifeIndex": 7, // int
            "nftPowerIndex": 23, //int
            "nftLdefIndex": "M229680162372527", // string
            "nftCharacId": "def", //string
            "qty": 100, // int
        },
    ]
};

var userMarketInfo = {
    "event": "nft_market", // string
    "action": "user_market_info", //string
    "actId": "APP0123456889", //string
    "walletId": "0xaC39b311DCEb2A4b2f5d8461c1cdaF756F4F7Ae9", //string
};

var shoppingCartList = {
    "event": "nft_market", //string
    "action": "nft_shopping_cart_list",
    "actId": "APP01234776789", //string
    "username": "Simonchen",
}

var shoppingCartChangeAdd = {
    "event": "nft_market", // string
    "action": "nft_shopping_cart_change",
    "actId": "APP01234776789", //string
    "operation": 0, // int 0 for add 1 for delete
    "username": "Simonchen", // user account
    "nftList": [ // array of string
        "M45923035625389", //  index of nft
    ]
};

var shoppingCartChangeDelete = {
    "event": "nft_market", // string
    "action": "nft_shopping_cart_change",
    "actId": "APP01234776789", //string
    "operation": 1, // int 0 for add 1 for delete
    "username": "Simonchen", // user account
    "nftList": [ // array of string
        "M45923035625389", //  index of nft
    ]
};

var purchaseHistory = {
    "event": "nft_market", // string
    "action": "nft_purchase_history",
    "actId": "APP01234776789", //string
    "username": "Simonchen",
};

// var start
// var end

ws.on('open', async function open() {
    // start = now()
    // ws.send(JSON.stringify(mkinfos))
    // ws.send(JSON.stringify(tokenPurchase))
    // ws.send(JSON.stringify(nftshow))
    // // ws.send(JSON.stringify(musicshow))
    // ws.send(JSON.stringify(itemDetails))
    // ws.send(JSON.stringify(purchaseConfirm))
    // ws.send(JSON.stringify(userMarketInfo));
    // ws.send(JSON.stringify(shoppingCartChangeAdd));
    // ws.send(JSON.stringify(shoppingCartChangeDelete));
    // ws.send(JSON.stringify(purchaseHistory));
    ws.send(JSON.stringify(shoppingCartList));
});

ws.on('message', function incoming(data) {
    console.log("")
    console.log(data)
    // end = now()
    // console.log("consuming time: " + parseFloat((end - start).toFixed(4)) / 1000 + "s");
});