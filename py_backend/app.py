# coding: utf8

from flask import Flask, render_template, session
from flask_socketio import SocketIO
from flask.ext.session import Session

from py_backend.model.core import *

db.create_all()


@app.route('/')
def hello_world():
    return 'Hello World!'


if __name__ == '__main__':
    app.run()
