from django.urls import path

from . import views

app_name = "backend_mvp"
urlpatterns = [
    # Interactive
    path("", views.index, name="index"),
    path("view-api-routes", views.view_api_routes, name="view_api_routes"),
    path("view-packages", views.view_packages, name="view_packages"),
    path("package/<int:package_id>/", views.view_pkg_detail, name="view_pkg_detail"),
    path("package/<int:package_id>/version/<int:version_id>", views.view_ver_detail, name="view_ver_detail"),
    # API
    path("api/v0/packages", views.packages),
]
