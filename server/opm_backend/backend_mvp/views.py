from django.shortcuts import render
from django.http import HttpResponse


def index(request):
    return HttpResponse("Hello, Odin. You're at the backend_mvp index.")


def list_api_routes(request):
    return HttpResponse("Hello, Odin. You're at the backend_mvp api list")


def list_packages(request):
    return HttpResponse("Should see some packages here ...")


def list_pkg_versions(request):
    return HttpResponse("Should see some package versions here ...")
