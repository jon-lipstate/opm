from django.shortcuts import get_object_or_404, render
from django.http import HttpResponse, JsonResponse
from .models import Package, Version


def index(request):
    return render(request, "backend_mvp/index.html")


def view_api_routes(request):
    return render(request, "backend_mvp/view_api_routes.html")


def view_packages(request):
    pkg_list = Package.objects.all
    return render(request, "backend_mvp/view_all_packages.html", {"pkg_list": pkg_list})


def view_pkg_detail(request, package_id):
    pkg = get_object_or_404(Package, pk=package_id)
    return render(request, "backend_mvp/view_pkg_detail.html", {"pkg": pkg})


def view_ver_detail(request, package_id, version_id):
    pkg = get_object_or_404(Package, pk=package_id)
    ver = get_object_or_404(Version, pk=version_id)
    return render(request, "backend_mvp/view_ver_detail.html", {"pkg": pkg, "ver": ver})


def packages(request):
    """ Return a list of all packages in the db, JSON format """
    resp_data = dict()
    pkg_list = Package.objects.all()
    pkg_list = pkg_list.values()

    print(len(pkg_list))

    for i in range(1, len(pkg_list) + 1):
        print(i)
        resp_data[i - 1] = pkg_list[i - 1]

    return JsonResponse(resp_data)


def create_package(request):
    # todo
    # Get JSON, validate, write to db
    return HttpResponse("todo")


def mod_package(request):
    # todo
    # Get JSON, validate, write to db
    return HttpResponse("todo")


def del_package(request):
    # todo
    # Get JSON, validate, write to db
    return HttpResponse("todo")


def create_version(request):
    ver = get_object_or_404(Version, pk=version_id)
    return HttpResponse("todo")


def organizations(request):
    ver = get_object_or_404(Version, pk=version_id)
    return HttpResponse("todo")


def create_organization(request):
    ver = get_object_or_404(Version, pk=version_id)
    return HttpResponse("todo")