from django.urls import path

from . import views

app_name = "backend_mvp"
urlpatterns = [
    # Interactive
    path("", views.index, name="index"),
    path("view-api-routes", views.view_api_routes, name="view_api_routes"),
    path("view-packages", views.view_packages, name="view_all_packages"),
    path("package/<int:package_id>/", views.view_pkg_detail, name="view_pkg_detail"),
    path("version/<int:version_id>", views.view_ver_detail, name="view_ver_detail"),
    # API
    path("api/v0/", views.api_v0, name="api_v0"),
    # API Auth
    path("api/v0/get_jwt", views.get_jwt, name="get_jwt"),
    # API Pkg Ops
    path("api/v0/packages", views.packages, name="api_packages"),
    path("api/v0/package/<int:package_id>/versions", views.versions_for_pkg, name="api_versions_for_pkg"),
    path("api/v0/create_package", views.create_package, name="api_create_package"),
    # API Version Ops
    path("api/v0/versions/all", views.versions_all, name="api_versions_all"),
    path("api/v0/version/<int:version_id>", views.version, name="api_version_single"),
    path("api/v0/create_version", views.create_version, name="api_create_version"),
]
