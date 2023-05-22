from django.shortcuts import render
from django.http import HttpResponse


def index(request):
    return HttpResponse("Hello, odin. You're at the backend_mvp index.")


def api_list(request):
    return HttpResponse("Hello, odin. You're at the backend_mvp api list")
