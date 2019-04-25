const W = require('ethereumjs-wallet');
import {AbiCoder} from 'web3-eth-abi';
import {Accounts} from 'web3-eth-accounts';
const coder = new AbiCoder();
const accounts = new Accounts('');

const publicChainConfig = require('../config/publicBridge.json');
const privateChainConfig = require('../config/privateBridge.json');

function encodeFunction(abi,funcName,...args) {
    let jsonInterface = searchForJsonInterface(abi,funcName);
    return coder.encodeFunctionCall(jsonInterface,args)
}

function encodePublicChainFunction(funcName,...args) {
    return encodeFunction(publicChainConfig.abi,funcName,...args)
}

function encodePrivateChainFunction(funcName,...args) {
    return encodeFunction(privateChainConfig.abi,funcName,...args)
}

function searchForJsonInterface(abi,funcName) {
    for(let i =0;i<abi.length;++i) {
        if (abi[i].name === funcName && abi[i].type === 'function') {
            return abi[i];
        }
    }
}

function getPublicChainBaseTxObject() {
    return {
        to: publicChainConfig.address,
        gasPrice: 10000000000,   // 10Gwei
        gas:"800000",
        value:"0",
    }
}

function getPrivateChainBaseTxObject() {
    return {
        to: privateChainConfig.address,
        gasPrice: 0,
        gas:"3000000",
        value:"0",
    }
}

function signTx(tx,privateKey) {
    return accounts.signTransaction(tx,privateKey).then(signedTx=>{
        return signedTx.rawTransaction;
    })
}

const createAccount = ()=> {
    let wallet =  W.generate();
    return {
        address: wallet.getAddress().toString('hex'),
        privateKey: wallet.getPrivateKey().toString('hex')
    }
};

export { createAccount, signTx, encodePublicChainFunction,encodePrivateChainFunction, getPublicChainBaseTxObject,getPrivateChainBaseTxObject};