"""
URL configuration for opm_backend project.
Ref : https://docs.djangoproject.com/en/4.2/topics/http/urls/
"""
from django.contrib import admin
from django.urls import include, path

urlpatterns = [
    path("backend_mvp/", include("backend_mvp.urls")),
    path('admin/', admin.site.urls),
]
