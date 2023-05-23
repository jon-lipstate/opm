from flask import Blueprint, redirect, url_for, session, jsonify
from flask_dance.contrib.github import github
from flask_dance.consumer import oauth_authorized
from os import getenv


auth_bp = Blueprint('login', __name__)
################################################################################################
@auth_bp.route("/github")
def index():
    session.clear()
    print("Redirect to ", url_for("github.login"))
    if not github.authorized:
        return redirect(url_for("github.login"))
    resp = github.get("/user")
    print("RESPUSER",resp)

    # assert resp.ok
    # return "You are @{login} on GitHub".format(login=resp.json()["login"])
    client_host = getenv("CLIENT_HOST")
    return redirect(client_host)

@auth_bp.route("/logout")
def logout():
    session.clear()
    client_host = getenv("CLIENT_HOST")
    return redirect(client_host)
################################################################################################

@auth_bp.route('/user', methods=['GET'])
def user():
    print("GITHUB",)
    if not github.authorized:
        return jsonify(user=None, isLoggedIn=False, isAdmin=False), 401  # Not Authorized

    resp = github.get("/user")
    print(">>>github user resp",resp.json())
    if not resp.ok:
        return jsonify(message="Could not fetch user info", isLoggedIn=False, isAdmin=False), 400

    user_data = resp.json()
    is_admin = check_admin(user_data)  # replace with your actual admin checking function
    
    return jsonify(user=user_data, isLoggedIn=True, isAdmin=is_admin), 200

