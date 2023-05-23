from flask import Blueprint, redirect, url_for, session, jsonify
from flask_dance.contrib.github import github
from flask_dance.consumer import oauth_authorized
from os import getenv


auth_bp = Blueprint('login', __name__)


@auth_bp.route("/logout")
def logout():
    session.clear()
    client_host = getenv("CLIENT_HOST")
    return redirect(client_host)


@auth_bp.route('/user', methods=['GET'])
def user():
    print()
    print("Begin auth.py, def user():")
    print(f"github.authorized = {github.authorized}")

    if not github.authorized:
        print("Returning (false branch) ... ")
        return jsonify(message="GitHub not authorized", user=None, isLoggedIn=False, isAdmin=False), 401  # Not Authorized

    resp = github.get("/user")
    if not resp.ok:
        return jsonify(message="Could not fetch user info", isLoggedIn=False, isAdmin=False), 400

    user_data = resp.json()
    print(user_data.keys())

    admin_list = ["Odin"]
    is_admin = user_data["login"] in admin_list
    print(f"is_admin = {is_admin}")

    print("Returning ... (true branch)")
    return jsonify(user=user_data, isLoggedIn=True, isAdmin=is_admin), 200

