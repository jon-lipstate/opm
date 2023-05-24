from django.db import models
from django.db.models import CharField


class Label(models.Model):
    name = models.CharField(max_length=200)
    # Could label development status (alpha, beta, pre-published, final etc.)
    # Could label various categories (Lib, Demo etc.)
    # Could label various types of licenses


class LabelGroup(models.Model):
    name = models.CharField(max_length=200)


class Package(models.Model):
    """ Represents an Odin Package """
    name = models.CharField(max_length=200)
    author = models.CharField(max_length=200)
    maintainers = models.CharField(max_length=200)
    created = models.DateTimeField("date created")
    # todo change versions so it points to version object(s) below
    versions = models.CharField(max_length=200)

    def __str__(self):
        return self.name


class Version(models.Model):
    """ Represents a single Odin Package Version """
    pkg_name = models.ForeignKey(Package, on_delete=models.RESTRICT, null=True)
    tag_name = models.CharField(max_length=200)
    tag_hash = models.CharField(max_length=200)
    content = models.CharField(max_length=200)
    author = models.CharField(max_length=200)
    created = models.DateTimeField("date created")

    def __str__(self):
        return f"{self.pkg_name}_{self.tag_name}"


class Org(models.Model):
    """ Represents an Organization (ex: Core Dev Team, Microsoft, Google etc.) """
    name: CharField = models.CharField(max_length=200)
    author = models.CharField(max_length=200)
    created = models.DateTimeField("date created")

    def __str__(self):
        return self.name
