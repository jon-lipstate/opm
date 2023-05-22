from django.db import models
from django.db.models import CharField


class Package(models.Model):
    """ Represents an Odin Package """
    name = models.CharField(max_length=200)
    author = models.CharField(max_length=200)
    maintainers = models.CharField(max_length=200)
    created = models.DateTimeField("date created")
    # todo change versions so it points to version object(s) below
    versions = models.CharField(max_length=200)
    # Type could be Lib, Demo etc.
    type = models.CharField(max_length=200)
    license = models.CharField(max_length=200)


class Version(models.Model):
    """ Represents a single Odin Package Version """
    version_number = models.CharField(max_length=200)
    content = models.CharField(max_length=200)
    author = models.CharField(max_length=200)
    created = models.DateTimeField("date created")
    # State could be alpha, beta, pre-published, final etc.
    state = models.CharField(max_length=200)


class Org(models.Model):
    """ Represents an Organization (ex: Core Dev Team, Microsoft, Google etc.) """
    name: CharField = models.CharField(max_length=200)
    author = models.CharField(max_length=200)
    created = models.DateTimeField("date created")
