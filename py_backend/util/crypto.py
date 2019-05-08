# encoding: utf-8

from Crypto.Cipher import AES
from Crypto import Random
# import typing lib for type hint usage, do not really need to import it.
import typing

KEY = b'\x05\xc0\xed\x93D\x1d\x8ef\xf44%$\x81\xc79K'
PAD = b'\0'
ENCRYPT_PATH = '../res/encrypted'
PLAINTEXT_PATH = '../res/plaintext'


def padding(s):
    return s + PAD * (AES.block_size - len(s) % AES.block_size)


def encrypt(message, key):
    message = padding(message)
    iv = Random.new().read(AES.block_size)
    cipher = AES.new(key, AES.MODE_CBC, iv)
    return iv + cipher.encrypt(message)


def decrypt(ciphertext, key):
    iv = ciphertext[:AES.block_size]
    cipher = AES.new(key, AES.MODE_CBC, iv)
    plaintext = cipher.decrypt(ciphertext[AES.block_size:])
    return plaintext.rstrip(PAD)


def encrypt_file(in_file: typing.BinaryIO, filename: str, key: bytes = KEY):
    plaintext = in_file.read()
    enc = encrypt(plaintext, key)
    with open(ENCRYPT_PATH + filename + '.enc', 'wb') as out_file:
        out_file.write(enc)
    return True


def decrypt_file(filename, key: bytes = KEY):
    with open(ENCRYPT_PATH + filename, 'rb') as in_file:
        cipher_text = in_file.read()
        dec = decrypt(cipher_text, key)
        with open(PLAINTEXT_PATH + filename[:-4], 'wb') as out_file:
            out_file.write(dec)
