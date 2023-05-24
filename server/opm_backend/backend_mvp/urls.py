from django.urls import path

from . import views

app_name = "backend_mvp"
urlpatterns = [
    path("", views.index, name="index"),
    path("api/v0/endpoints", views.list_api_routes),
    path("api/v0/packages", views.list_packages),
    path("api/v0/package/<int:package_id>/", views.view_pkg_detail),
    path("api/v0/package/<int:version_id>/", views.view_ver_detail),
]
