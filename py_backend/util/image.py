# encoding: utf-8

import io
import cv2
import numpy as np
from PIL import Image


def resize_to_file(file, filename, a=200, b=200):
    img = Image.open(file)
    img = np.array(img)
    img = cv2.resize(img, (a, b))
    cv2.imwrite(filename, img)


def resize_to_object(file, filename, a=200, b=200):
    img = Image.open(file)
    img = np.array(img)
    img = cv2.resize(img, (a, b))
    is_success, buffer = cv2.imencode('.jpg', img)
    io_buf = io.BytesIO(buffer)
    return io_buf
