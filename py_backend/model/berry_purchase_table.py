# coding: utf8

from .core import db
from datetime import datetime


class BerryPurchaseTable(db.Model):
    TransactionId = db.Column(db.String, primary_key=True)
    RefillAsId = db.Column(db.String)
    NumPurchased = db.Column(db.Integer)
    Timestamp = db.Column(db.DateTime)
    AppTranId = db.Column(db.String)
    AppId = db.Column(db.String)
    Status = db.Column(db.Integer)

    def __init__(self,
                 transaction_id: str,
                 refill_as_id: str,
                 num_purchased: int,
                 timestamp: datetime,
                 app_tran_id: str,
                 app_id: str,
                 status: int):
        self.TransactionId = transaction_id
        self.RefillAsId = refill_as_id
        self.NumPurchased = num_purchased
        self.Timestamp = timestamp
        self.AppTranId = app_tran_id
        self.AppId = app_id
        self.Status = status


    def __repr__(self):
        return '<Transaction Id %r>' % self.TransactionId
