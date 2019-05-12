# encoding: utf-8

import web3
import json

from web3 import Web3
from web3.contract import ConciseContract
from web3.middleware import geth_poa_middleware


def make_connection():
    w3 = Web3(Web3.HTTPProvider('http://3.1.35.36:8540'))
    w3.middleware_stack.inject(geth_poa_middleware, layer=0)
    w3.eth.defaultAccount = w3.eth.accounts[0]

    with open('abi.txt', 'r') as f:
        abi = json.load(f)

    nft = w3.eth.contract(abi=abi, address=Web3.toChecksumAddress('0x1ee1d5178c8cc69796293976df78b36789f12f75'))

    return nft
