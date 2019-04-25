const solc = require('solc')
const fs = require('fs-extra')
const fsSync = require('fs')
const path = require('path')
var program = require('commander')
var Web3 = require('web3');
const {
    promisify
} = require('util')
const readFile = promisify(fs.readFile)
const readJson = promisify(fs.readJson)
const writeJson = promisify(fs.writeJson)
program
    .version('1.0.0')
// Related file should be in same directory
program
    .command('deploy <file> <contract>')
    .description('deploy contract in  file')
    .option('-r, --rpcport <port>', 'rpc path ')
    .option('-a, --arguments <arguments>', 'contract constructor arguments')
    .option('-k, --keystore <path>', 'keystore file to deploy contract')
    .option('-p, --password <password>', 'password to decrpt keystore')
    .option('-s, --save', 'kind to save contract Address 0 for public, else for private') // only used for demo project
    .action(async function(file, contractName, options) {

        // Exam parameter
        if (!options.keystore) {
            throw new Error('no keystore file specified')
        }
        if (!options.password) {
            throw new Error('no password specified')
        }

        if (!options.rpcport) {
            throw new Error("specify rpcport ")
        }
        // connect to api provider
        let web3 = new Web3(new Web3.providers.HttpProvider(options.rpcport))
        try {
            let ks = await readJson(options.keystore)
            let account = web3.eth.accounts.decrypt(ks, options.password)

            // compile contract
            let [bytecode, abi] = await compileContract(file, contractName)

            // pack transaction to deploy contract
            let contract = new web3.eth.Contract(abi)
            let dp = contract.deploy({
                data: '0x' + bytecode,
                arguments: !options.arguments ? [] : JSON.parse(options.arguments)
            })
            let tx = {
                data: dp.encodeABI(),
                from: account.address,
                gas: "3000000",
                gasPrice: '0',
                value: "0", // 1000 ether
            }

            // sign and send transaction
            let signedTx = await account.signTransaction(tx)
            let receipt = await web3.eth.sendSignedTransaction(signedTx.rawTransaction)
            if (receipt.status === '0x0') {
                throw new Error("contract deploy revert")
            }

            console.log("deploy contract at: " + receipt.contractAddress)

            // if save Save to specified location
            if (options.save != undefined) {
                let dst = '../etc/' + path.basename(file).split('.')[0] + '.json'
                let config
                if (fsSync.existsSync(dst)) {
                    config = await readJson(dst)
                } else {
                    config = {}
                }
                config.address = receipt.contractAddress
                await writeJson(dst, config)
            }
        } catch (err) {
            console.log(err)
        }
    })

program
    .command('genkey <password>')
    .description('use password to generate keystore')
    .option('-f, --file <path>', "path to store keystore")
    .action(function(password, options) {
        web3 = new Web3()
        account = web3.eth.accounts.create()
        ks = account.encrypt(password)
        console.log(account.address)
        if (options.file) {
            fs.writeJson(options.file, ks, () => console.log("write key store to", options.file))
        }
    })

program.parse(process.argv)

async function compileContract(file, contractName) {
    try {
        let code = await readFile(file)
        code = code.toString()
        let input = {
            language: 'Solidity',
            sources: {},
            settings: {
                optimizer: {
                    enabled: true
                },
                outputSelection: {
                    '*': {
                        '*': ['*']
                    }
                }
            }
        }

        let fileName = path.basename(file)
        input.sources[fileName] = {
            content: code
        }

        let findImports = package => {
            try {
                let baseName = path.dirname(file);
                package = path.join(baseName, package)
                let content = fsSync.readFileSync(package)
                return {
                    contents: content.toString()
                }
            } catch (err) {
                throw err
            }
        }


        let output = JSON.parse(solc.compile(JSON.stringify(input), findImports))
        if (output.errors) {
            throw output.errors
        }
        let compiledObject = output.contracts[fileName][contractName]
        let bytecode = compiledObject.evm.bytecode.object
        let abi = output.contracts[fileName][contractName].abi
        return [bytecode, abi]
    } catch (err) {
        throw err
    }
}