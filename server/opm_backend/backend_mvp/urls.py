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
    path("api/v0/get_jwt", views.get_jwt, name="get_jwt"),
    path("api/v0/packages", views.packages, name="api_packages"),
    path("api/v0/create_package", views.create_package, name="api_create_package"),
    path("api/v0/package/<int:package_id>/versions", views.versions, name="api_versions"),
    path("api/v0/package/<int:package_id>/version/<int:version_id>", views.version, name="api_single_version"),
]
