from django.urls import path

from . import views

urlpatterns = [
    path("", views.index, name="index"),
    path("api/v0/endpoints", views.list_api_routes),
    path("api/v0/packages", views.list_packages),
    path("api/v0/1/versions", views.list_pkg_versions),
]
