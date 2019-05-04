# coding: utf8

from .core import db


class NftMappingTable(db.Model):
    NftLdefIndex = db.Column(db.String, primary_key=True)
    TypeId = db.Column(db.String)
    FileName = db.Column(db.String)
    Key = db.Column(db.String)
    NftAdminId = db.Column(db.String)

    def __init__(self,
                 NftLdefIndex: str,
                 TypeId: str,
                 FileName: str,
                 Key: str,
                 NftAdminId: str):
        self.NftLdefIndex = NftLdefIndex
        self.TypeId = TypeId
        self.FileName = FileName
        self.Key = Key
        self.NftAdminId = NftAdminId

    def __repr__(self):
        return '<NftLdefIndex %r>' % self.NftLdefIndex
