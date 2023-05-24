""" Views for Backend MVP """

# imports
# base
import json
# open-source
#   Django Standard
from django.core.validators import ValidationError
from django.http import HttpResponse, JsonResponse
from django.http import HttpResponseBadRequest   # 400
from django.http import HttpResponseForbidden    # 403
from django.http import HttpResponseNotFound     # 404
from django.http import HttpResponseNotAllowed   # 405
from django.http import HttpResponseServerError  # 500
from django.shortcuts import get_object_or_404, render
from django.utils import timezone
from django.views.decorators.csrf import csrf_exempt
# local
from .models import Package, Version


# Error messages
error_msg_400 = """
<h1>HTTP ERROR 400</h1>
<p>You're a baaaaaaaaaad request ... do you even wash your hands?</p>
"""

error_msg_405 = """
<h1>HTTP ERROR 405</h1>
<p>Oh fiddle sticks, you used a disallowed method ...</p>
"""

error_msg_500 = """
<h1>HTTP ERROR 500</h1>
<p>You crashed the server! +100 raven points!</p>
"""


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


def api_v0(request):
    return HttpResponse("Welcome to API, v0")


def get_jwt(request):
    # Expect GitHub Token in response body
    if not "token is present":
        return HttpResponse("Nope, not authorized")
    else:
        return HttpResponse("Welcome authorized user")


def packages(request):
    """ API, Return a list of all packages in the db, JSON format """
    resp_data = dict()
    pkg_list = Package.objects.all()
    print(f"Query Set : pkg_list got {len(pkg_list)} packages!")
    pkg_list = pkg_list.values()
    for i in range(1, len(pkg_list) + 1):
        resp_data[i - 1] = pkg_list[i - 1]

    return JsonResponse(resp_data)


@csrf_exempt
def create_package(request):
    allowed_methods = ["POST"]
    if request.method not in allowed_methods:
        return HttpResponseNotAllowed(allowed_methods, error_msg_405)
    else:
        try:
            json_data = json.loads(request.body)
        except json.JSONDecodeError as e:
            print("Error: JSON Decode Error")
            print(e)
            return HttpResponseServerError(error_msg_500)

        try:
            new_pkg = Package()
            # Get these from request body
            for key, val in json_data.items():
                setattr(new_pkg, key, val)
            # Set these ourselves
            new_pkg.created = timezone.now()
            new_pkg.versions = "None"
            # Validate
            new_pkg.full_clean()
        except ValidationError as e:
            print("Error: Model validation error")
            print(e)
            return HttpResponseBadRequest(error_msg_400 + str(e))

        new_pkg.save()

        return HttpResponse("200 Success")


def mod_package(request):
    # todo
    # Get JSON, validate, write to db
    return HttpResponse("todo")


def del_package(request):
    # todo
    # Get JSON, validate, write to db
    return HttpResponse("todo")


def versions(request, package_id):
    """ API, For one package, return a list of all assoc. Versions in the db, JSON format """

    pkg = get_object_or_404(Package, pk=package_id)
    vers = pkg.version_set.all()
    print(f"Query Set : vers got {len(vers)} versions!")
    vers = vers.values()
    resp_data = dict()
    for i in range(1, len(vers) + 1):
        resp_data[i - 1] = vers[i - 1]
        print(resp_data[i - 1])

    return JsonResponse(resp_data)


def create_version(request):
    pass
    return HttpResponse("todo")


def orgs(request):
    """ API, Return a list of all orgs in the db, JSON format """
    pass
    return HttpResponse("todo")


def create_org(request):
    pass
    return HttpResponse("todo")
