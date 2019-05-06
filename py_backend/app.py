# coding: utf8

from flask import Flask, render_template, session, request
from flask_socketio import SocketIO
from flask_session import Session
from configparser import ConfigParser
from py_backend.model.core import *

config_file = ConfigParser()
config_file.read('app.conf')
assert 'basic' in config_file, 'app.conf file do not have basic section'
assert 'dev' in config_file, 'app.conf file do not have dev section'
assert 'prod' in config_file, 'app.conf file do not have prod section'
config_parameter_setting = config_file[config_file['basic']['runmode']]

# Create Tables
db.create_all()


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
    if kind == 'avatar':
        pass
    elif kind == 'dat':
        pass
    else:
        assert False, 'unknown file type'


@app.route('/balance/<user>')
def nft_balance_handler(user: str):
    pass


@app.route('/nftList/<user>')
def nft_list_handler(user):
    pass


if __name__ == '__main__':
    app.run()
