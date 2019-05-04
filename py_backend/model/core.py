# coding: utf-8

from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from datetime import datetime

app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'mysql+pymysql://test:test123@localhost:3306/music_and_dat'
db = SQLAlchemy(app)


class BerryPurchaseTable(db.Model):
    TransactionId = db.Column(db.String(80), primary_key=True)
    RefillAsId = db.Column(db.String(80))
    NumPurchased = db.Column(db.Integer)
    Timestamp = db.Column(db.DateTime)
    AppTranId = db.Column(db.String(80))
    AppId = db.Column(db.String(80))
    Status = db.Column(db.Integer)

    def __init__(self,
                 transaction_id: str,
                 refill_as_id: str,
                 num_purchase: int,
                 timestamp: datetime,
                 app_tran_id: str,
                 app_id: str,
                 status: int):
        self.TransactionId = transaction_id
        self.RefillAsId = refill_as_id
        self.NumPurchased = num_purchase
        self.Timestamp = timestamp
        self.AppTranId = app_tran_id
        self.AppId = app_id
        self.Status = status

    def __repr__(self):
        return '<Transaction Id {}>'.format(self.TransactionId)


class NftInfoTable(db.Model):
    NftLdefIndex = db.Column(db.String(80), primary_key=True)
    NftType = db.Column(db.String(80))
    NftName = db.Column(db.String(80))
    DistIndex = db.Column(db.String(80))
    NftLifeIndex = db.Column(db.Integer)
    NftPowerIndex = db.Column(db.Integer)
    NftCharacId = db.Column(db.String(80))
    PublicKey = db.Column(db.String(80))

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
        return '<NftLdefIndex {}>'.format(self.NftLdefIndex)


class NftItemAdmin(db.Model):
    NftLdefIndex = db.Column(db.String(80), primary_key=True)
    ShortDescription = db.Column(db.String(80))
    LongDescription = db.Column(db.String(80))
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
        return '<NftLdefIndex {}>'.format(self.NftLdefIndex)


class NftMappingTable(db.Model):
    NftLdefIndex = db.Column(db.String(80), primary_key=True)
    TypeId = db.Column(db.String(80))
    FileName = db.Column(db.String(80))
    Key = db.Column(db.String(80))
    NftAdminId = db.Column(db.String(80))

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
        return '<NftLdefIndex {}>'.format(self.NftLdefIndex)


class NftMarketTable(db.Model):
    NftLdefIndex = db.Column(db.String(80), primary_key=True)
    MpId = db.Column(db.String(80))
    NftAdminId = db.Column(db.String(80))
    Price = db.Column(db.Integer)
    Qty = db.Column(db.Integer)
    NumSold = db.Column(db.Integer)
    Active = db.Column(db.Boolean)
    ActiveTicker = db.Column(db.String(80))

    def __repr__(self):
        return '<NftLdefIndex {}>'.format(self.NftLdefIndex)

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


class StorePurchaseHistory(db.Model):
    PurchaseId = db.Column(db.String(80), primary_key=True)
    AsId = db.Column(db.String(80))
    TransactionAddress = db.Column(db.String(80))
    NftName = db.Column(db.String(80))
    TotalPaid = db.Column(db.String(80))
    NftLdefIndex = db.Column(db.String(80))
    Timestamp = db.Column(db.DateTime)
    Status = db.Column(db.Integer)

    def __init__(self,
                 purchase_id: str,
                 asid: str,
                 transaction_address: str,
                 nft_name: str,
                 total_paid: int,
                 nftl_def_index: str,
                 time_stamp: datetime,
                 status: int):
        self.PurchaseId = purchase_id
        self.AsId = asid
        self.TransactionAddress = transaction_address
        self.NftName = nft_name
        self.TotalPaid = total_paid
        self.NftLdefIndex = nftl_def_index
        self.Timestamp = time_stamp
        self.Status = status

    def __repr__(self):
        return '<PurchaseId {}>'.format(self.PurchaseId)
