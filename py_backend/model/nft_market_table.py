# coding: utf8

from .core import db


class NftMarketTable(db.Model):
    NftLdefIndex = db.Column(db.String, primary_key=True)
    MpId = db.Column(db.String)
    NftAdminId = db.Column(db.String)
    Price = db.Column(db.Integer)
    Qty = db.Column(db.Integer)
    NumSold = db.Column(db.Integer)
    Active = db.Column(db.Boolean)
    ActiveTicker = db.Column(db.String)

    def __repr__(self):
        return '<NftLdefIndex %r>' % self.NftLdefIndex

    def __init__(self,
                 NftLdefIndex: str,
                 MpId: str,
                 NftAdminId: str,
                 Price: int,
                 Qty: int,
                 NumSold: int,
                 Active: bool,
                 ActiveTicker: str):
        self.NftLdefIndex = NftLdefIndex
        self.MpId = MpId
        self.NftAdminId = NftAdminId
        self.Price = Price
        self.Qty = Qty
        self.NumSold = NumSold
        self.Active = Active
        self.ActiveTicker = ActiveTicker
