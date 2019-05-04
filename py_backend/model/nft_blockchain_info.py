# coding: utf8

from .core import db


class NftInfoTable(db.Model):
    NftLdefIndex = db.Column(db.String, primary_key=True)
    NftType = db.Column(db.String)
    NftName = db.Column(db.String)
    DistIndex = db.Column(db.String)
    NftLifeIndex = db.Column(db.Integer)
    NftPowerIndex = db.Column(db.Integer)
    NftCharacId = db.Column(db.String)
    PublicKey = db.Column(db.String)

    def __init__(self,
                 NftLdefIndex: str,
                 nft_type: str,
                 nft_name: str,
                 dist_index: str,
                 nft_life_index: int,
                 nft_power_index: int,
                 nft_charac_id: str,
                 pubic_key: str):
        self.NftLdefIndex = NftLdefIndex
        self.NftType = nft_type
        self.NftName = nft_name
        self.DistIndex = dist_index
        self.NftLifeIndex = nft_life_index
        self.NftPowerIndex = nft_power_index
        self.NftCharacId = nft_charac_id
        self.PublicKey = pubic_key

    def __repr__(self):
        return '<NftLdefIndex %r>' % self.NftLdefIndex
