from flask import Blueprint, redirect, url_for, session
from flask_dance.contrib.github import github
from flask_dance.consumer import oauth_authorized
from os import getenv


login_bp = Blueprint('login', __name__)

@login_bp.route("/")
def index():
    if not github.authorized:
        print("Redirect to ", url_for("github.login"))
        return redirect(url_for("github.login"))
    print("AFTERRR",github.authorized)
    print("AFTERRR",github.access_token)
    resp = github.get("/user")
    print("RESPUSER",resp)

    # assert resp.ok
    # return "You are @{login} on GitHub".format(login=resp.json()["login"])
    client_host = getenv("CLIENT_HOST")
    return redirect(client_host)

@oauth_authorized.connect
def logged_in(blueprint, token):
    if not token:
        flash("Failed to log in.", category="error")
        return False

    resp = blueprint.session.get("/user")
    if not resp.ok:
        msg = "Failed to fetch user info."
        flash(msg, category="error")
        return False

    github_info = resp.json()

@login_bp.route("/logout")
def logout():
    session.clear()
    client_host = getenv("CLIENT_HOST")
    return redirect(client_host)