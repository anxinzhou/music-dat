# coding: utf8

from flask import Flask, render_template, session, request, jsonify
from flask_socketio import SocketIO
from flask_session import Session
from configparser import ConfigParser
from model.core import *
from util.crypto import encrypt_file, decrypt_file
from util.image import resize_to_object
from util.contract_loader import make_connection
from werkzeug.utils import secure_filename

from web3 import Web3

config_file = ConfigParser()
config_file.read('app.conf')
assert 'basic' in config_file, 'app.conf file do not have basic section'
assert 'dev' in config_file, 'app.conf file do not have dev section'
assert 'prod' in config_file, 'app.conf file do not have prod section'
config_parameter_setting = config_file[config_file['basic']['runmode']]

# Create Tables
db.create_all()

# define upload file folders
app.config['UPLOAD_FOLDER'] = './upload'


# define a set of permitted file format
# ALLOWED_EXTENSION = {'jpg', 'jpeg', 'png', 'gif', 'mp3', 'bmp'}


@app.route('/')
def hello_world():
    return 'Hello World!'


@app.route('/admin', methods='POST')
def admin_handler():
    form = request.form  # the body of http post request

    email = form.get('email', '')
    password = form.get('password', '')

    assert email != '' and password != ''
    print(email)
    print(password)

    # TODO handle login later
    return '', 200


@app.route('/ws')
def websocket_handler():
    pass


@app.route('/file/<kind>', methods='POST')
def file_upload_handler(kind):
    f = request.files['file']
    file_name = secure_filename(f.filename)
    if kind == 'avatar':
        f = resize_to_object(f, file_name)
    encrypt_file(f, file_name)

    # create nft token on blockchain
    testAddress = Web3.toChecksumAddress("0x3c62aa7913bc303ee4b9c07df87b556b6770e3fc")
    tokenId = 3
    nftType = "1"
    nftName = "1"
    nftLdefIndex = "1"
    distIndex = "1"
    nftLifeIndex = 100
    nftPowerIndex = 100
    nftCharacId = "1"
    publicKey = Web3.toBytes(1)

    nft = make_connection()

    nft.functions.mint(
        testAddress,
        nftType,
        nftName,
        nftLdefIndex,
        distIndex,
        nftLifeIndex,
        nftPowerIndex,
        nftCharacId,
        publicKey
    ).transact()

    # insert into mysql db
    nft_info_record = NftInfoTable(
        NftLdefIndex=nftLdefIndex,
        nft_type=nftType,
        nft_name=nftName,
        dist_index=distIndex,
        nft_life_index=nftLifeIndex,
        nft_power_index=nftPowerIndex,
        nft_charac_id=nftCharacId,
        pubic_key=str(publicKey)
    )

    db.session.add(nft_info_record)

    type_id = '01'
    if kind != 'avatar':
        type_id = '02'

    nft_admin_id = 'sth'  # TODO
    nft_mapping_record = NftMappingTable(
        NftLdefIndex=nftLdefIndex,
        TypeId=type_id,
        FileName=file_name,
        Key='0x01',
        NftAdminId=nft_admin_id
    )
    db.session.add(nft_mapping_record)

    nft_market_record = NftMarketTable(
        NftLdefIndex=nftLdefIndex,
        MpId='sth',
        NftAdminId=nft_admin_id,
        Price=1,  # TODO fake price
        Qty=1,  # TODO fake qty
        NumSold=0,
        Active=True,
        ActiveTicker='sth'  # TODO
    )
    db.session.add(nft_market_record)

    nft_item_admin_record = NftItemAdmin(
        NftLdefIndex=nftLdefIndex,
        short='short descption',  # TODO
        long='long description',  # TODO
        num=1  # TODO qty
    )

    db.session.add(nft_item_admin_record)

    db.session.commit()


@app.route('/balance/<user>')
def nft_balance_handler(user: str):
    nft = make_connection()
    web3_user = Web3.toChecksumAddress(user)
    count = nft.functions.balanceOf(web3_user).call()

    return jsonify({'Count': count}), 200


@app.route('/nftList/<user>')
def nft_list_handler(user):
    web3_user = Web3.toChecksumAddress(user)
    nft = make_connection()
    nft_list = nft.functions.tokensOfUser(web3_user).call()

    pass


if __name__ == '__main__':
    app.run()
