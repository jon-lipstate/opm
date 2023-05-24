from django.shortcuts import get_object_or_404, render
from django.http import HttpResponse
from .models import Package, Version


def index(request):
    return render(request, "backend_mvp/index.html")


def list_api_routes(request):
    return HttpResponse("Hello, Odin. You're at the backend_mvp api list")


def list_packages(request):
    pkg_list = Package.objects.all
    return render(request, "backend_mvp/packages.html", {"pkg_list": pkg_list})


def view_pkg_detail(request, package_id):
    pkg = get_object_or_404(Package, pk=package_id)
    return render(request, "backend_mvp/pkg_detail.html", {"pkg": pkg})


def view_ver_detail(request, version_id):
    ver = get_object_or_404(Version, pk=version_id)
    return render(request, "backend_mvp/ver_detail.html", {"ver": ver})


def api_make_package(request):
    # todo
    # Get JSON, validate, write to db
    return HttpResponse("todo")


def api_make_version(request):
    ver = get_object_or_404(Version, pk=version_id)
    return HttpResponse("todo")


def api_make_organization(request):
    ver = get_object_or_404(Version, pk=version_id)
    return HttpResponse("todo")