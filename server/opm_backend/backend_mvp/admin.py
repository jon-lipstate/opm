from django.contrib import admin
from .models import Package
from .models import Version
from .models import Org

# Register your models here.

admin.site.register(Package)
admin.site.register(Version)
admin.site.register(Org)