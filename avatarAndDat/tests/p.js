const WebSocket = require('ws');

// const ws = new WebSocket('ws://3.1.35.36:4000/ws');
const ws = new WebSocket('ws://localhost:4000/ws');
// const ws = new WebSocket("ws://18.136.134.84:4000");
// const now = require('performance-now');

let testNftLdefIndex = "M982980481";
let tokenId = "396285078";
let nickname = "AlphaBrain";
let testnickname = "testnickname2";
//done
var mkinfos = {
    action: "mp_list",
    event: "market",
    actId: "A1201212",
    supportedType: "721-04"
};

var purchaseConfirm = {
    "event": "nft_market", // string
    "action": "NFT_purchase_confirm", // string
    "actId": "APP0123456789", //string
    "nickname": "AlphaBrain", // string
    "nftTranData": [ //array of object
        {
            "supportedType": "721-02", // string "721-04" for music "721-02" for avatar
            "nftLdefIndex": "A99290519599110", // string
        },
    ]
};

var purchaseHistory = {
    "event": "nft_market", // string
    "action": "nft_purchase_history",
    "actId": "APP01234776789", //string
    "nickname": nickname,
    "supportedType": "721-02"
};

var tokenPurchasePending = {
    "event": "token_purchase", //string
    "action": "tokenbuy_paid", //string
    "actId": "APP01234776789", //string
    "nickname": nickname,
    "appTranId": "4fffc3b54bd8270df6fd85d1f17b8ec53adfa64d882",
    "appId": "XXX@appleid.com",
    "amount": 1, //int    berry number
    "actionStatus": 2
};

var tokenPurchaseFinish = {
    "event": "token_purchase", //string
    "action": "tokenbuy_paid", //string
    "actId": "APP01234776789", //string
    "nickname": nickname,
    "appTranId": "4fffc3b54bd8270df6fd85d1f17b8ec53adfa64d882",
    "transactionId": "bdacab81f70af17760819bd1cd0d27ba1cea165a2a7f2305257ae98cd2e01c18",
    "appId": "XXX@appleid.com",
    "amount": 1, //int    berry number
    "actionStatus": 1
};

var nftshow = {
    "event": "nft_show", //string
    "action": "NFT_display", //string
    "actId": "APP01234776789", //string
    "nftLdefIndex": "M160196874231178", //string  // nft to show
    "supportedType": "721-04", // string "721-04" for music "721-02" for avatar
};

var itemDetails = {
    "event": "nft_market", //string
    "action": "item_details", //string
    "actId": "APP0123456889", //string
    "nftTranData": [ // array of object
        {
            "nftLdefIndex": "M120398914028234",
            "supportedType": "721-04", // string "721-04" for music "721-02" for avatar
        },
    ]
};

var marketUserList = {
    "event": "nft_market", //string
    "action": "market_user_list", //string
    "actId": "APP0123456889", //string
    "nickname": testnickname,
};

var userMarketInfo = {
    "event": "nft_market", // string
    "action": "user_market_info", //string
    "actId": "APP0123456889", //string
    "nickname": nickname, //string
};

var shoppingCartChangeAdd = {
    "event": "nft_market", // string
    "action": "nft_shopping_cart_change",
    "actId": "APP01234776789", //string
    "operation": 0, // int 0 for add 1 for delete
    "nickname": nickname, // user account
    "nftList": [ // array of string
        "M90803573804276", //  index of nft
        "M60828873421037",
    ]
};

var shoppingCartList = {
    "event": "nft_market", //string
    "action": "nft_shopping_cart_list",
    "actId": "APP01234776789", //string
    "nickname": nickname,
};

var shoppingCartChangeDelete = {
    "event": "nft_market", // string
    "action": "nft_shopping_cart_change",
    "actId": "APP01234776789", //string
    "operation": 1, // int 0 for add 1 for delete
    "nickname": nickname, // user account
    "nftList": [ // array of string
        "M60828873421037", //  index of nft
        "M90803573804276",
    ]
};

var nftTransfer = {
    "event": "nft_market",  // string
    "action": "nft_transfer", // string
    "actId": "APP01234776789", // string
    "senderNickname": nickname, // string
    "receiverNickname": nickname, // string
    "nftTranData": {   // object of nft
        "supportedType": "721-04", // string "721-04" for music "721-02" for avatar,
        "nftLdefIndex": "M60828873421037", // string
    }
};

var BindWallet = {
    "event": "user_activity",  // string
    "action": "bind_wallet", // string
    "actId": "APP01234776789", // string
    "nickname": testnickname,
    "walletId": "0x3c62aa7913bc303ee4b9c07df87b556b6770e3fc",
};

var setNickname = {
    "event": "user_activity",  // string
    "action": "set_nickname", // string
    "actId": "APP01234776789", // string
    "uuid": "AKkjkHCHcbvMq6oLWPUEk2Wf",
    "nickname": testnickname
}

var isNicknameDuplicated = {
    "event": "user_activity", // string
    "action": "is_nickname_duplicated", // string
    "actId": "APP01234776789", // string
    "nickname": nickname,  //string
};

var followList = {
    "event": "user_activity",  // string
    "action": "follow_list", // string
    "actId": "APP01234776789", // string
    "nickname": testnickname,
};

var followOperationAdd = {
    "event": "user_activity",  // string
    "action": "follow_list_operation", // string
    "actId": "APP01234776789", // string
    "nickname": testnickname, // string
    "followNickname": nickname, // string
    "operation": 0,
};

var followOperationDelete = {
    "event": "user_activity",  // string
    "action": "follow_list_operation", // string
    "actId": "APP01234776789", // string
    "nickname": testnickname, // string
    "followNickname": nickname, // string
    "operation": 1,
};

var isNickNameSet = {
    "event": "user_activity",  // string
    "action": "is_nickname_set", // string
    "actId": "APP01234776789", // string
    "uuid": "AKkjkHCHcbvMq6oLWPUEk2Wf",
}

// var start
// var end

ws.on('open', async function open() {
    // start = now()
    // ws.send(JSON.stringify(mkinfos));
    // ws.send(JSON.stringify(purchaseConfirm));
    ws.send(JSON.stringify(purchaseHistory));
    // ws.send(JSON.stringify(tokenPurchasePending));
    // ws.send(JSON.stringify(tokenPurchaseFinish));
    // ws.send(JSON.stringify(nftshow))
    // ws.send(JSON.stringify(itemDetails))
    // ws.send(JSON.stringify(marketUserList));
    // ws.send(JSON.stringify(userMarketInfo));
    // ws.send(JSON.stringify(shoppingCartChangeAdd));
    // ws.send(JSON.stringify(shoppingCartList));
    // ws.send(JSON.stringify(shoppingCartChangeDelete));
    // ws.send(JSON.stringify(nftTransfer));
    // ws.send(JSON.stringify(BindWallet));
    // ws.send(JSON.stringify(setNickname));
    // ws.send(JSON.stringify(isNicknameDuplicated));
    // ws.send(JSON.stringify(followList));
    // ws.send(JSON.stringify(followOperationAdd));
    // ws.send(JSON.stringify(followOperationDelete));
    // ws.send(JSON.stringify(isNickNameSet));
});

ws.on('message', function incoming(data) {
    console.log("")
    console.log(data)
    // end = now()
    // console.log("consuming time: " + parseFloat((end - start).toFixed(4)) / 1000 + "s");
});