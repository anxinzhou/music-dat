const ethers = require('ethers');
function createWallet(mnemonic) {
    let wallet = ethers.Wallet.fromMnemonic(mnemonic);
    return {
        address: wallet.address,
        privateKey: wallet.privateKey,
    }
}

export {createWallet}
