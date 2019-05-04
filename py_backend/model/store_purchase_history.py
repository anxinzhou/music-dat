# coding: utf-8

from .core import db
from datetime import datetime


class StorePurchaseHistory(db.Model):
    PurchaseId = db.Column(db.String, primary_key=True)
    AsId = db.Column(db.String)
    TransactionAddress = db.Column(db.String)
    NftName = db.Column(db.String)
    TotalPaid = db.Column(db.String)
    NftLdefIndex = db.Column(db.String)
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
