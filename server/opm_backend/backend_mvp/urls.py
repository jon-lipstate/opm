from django.urls import path

from . import views

urlpatterns = [
    path("", views.index, name="index"),
    path("api/v0/list", views.api_list),
]