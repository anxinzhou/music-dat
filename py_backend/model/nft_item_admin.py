# coding: utf8

from .core import db


class NftItemAdmin(db.Model):
    NftLdefIndex = db.Column(db.String, primary_key=True)
    ShortDescription = db.Column(db.String)
    LongDescription = db.Column(db.String)
    NumDistribution = db.Column(db.Integer)

    def __init__(self,
                 NftLdefIndex: str,
                 short: str,
                 long: str,
                 num: int):
        self.NftLdefIndex = NftLdefIndex
        self.ShortDescription = short
        self.LongDescription = long
        self.NumDistribution = num

    def __repr__(self):
        return '<NftLdefIndex %r>' % self.NftLdefIndex
