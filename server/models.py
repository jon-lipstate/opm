"""Data models."""
from . import db

class packages(db.Model):
    """Data model for Package manager"""

    __tablename__ = 'packages'
    __table_args__ = {'schema': 'package_manager'}
    packageID = db.Column(
        db.Double,
        primary_key=True,
        nullable=False,
        unique=True
    )

    description = db.Column(
        db.Double,
        index=False,
        nullable=True
    )


    def __repr__(self):
        return '<User {}>'.format(self.packageID)



# This example is more readeable but very verbose to write out
# example from https://hackersandslackers.com/flask-sqlalchemy-database-models/
class User(db.Model):
    """Data model for user accounts."""

    __tablename__ = 'user'
    __table_args__ = {'schema': 'user_manager'}
    id = db.Column(
        db.Integer,
        primary_key=True,
        nullable=False
    )
    username = db.Column(
        db.String(64),
        index=True,
        unique=True,
        nullable=False
    )
    email = db.Column(
        db.String(255),
        index=True,
        unique=True,
        nullable=False
    )
    created = db.Column(
        db.DateTime,
        index=False,
        unique=False,
        nullable=False
    )
    bio = db.Column(
        db.Text,
        index=False,
        unique=False,
        nullable=True
    )
    admin = db.Column(
        db.Boolean,
        index=False,
        unique=False,
        nullable=False
    )

    # Can also add functions here
    def __init__(self, username):
        self.username = username

    def to_json(self):
        return dict(name=self.username, email=self.email, id=self.id)

    def __repr__(self):
        return '<User {}>'.format(self.username)