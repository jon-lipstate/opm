from flask import Blueprint, jsonify
from flask_dance.contrib.github import github

user_bp = Blueprint('user', __name__)

@user_bp.route('/user', methods=['GET'])
def user():
    # print("GITHUB",dir(github))
    if not github.authorized:
        return jsonify(user=None, isLoggedIn=False, isAdmin=False), 401  # Not Authorized

    resp = github.get("/user")
    print(">>>github user resp",resp.json())
    if not resp.ok:
        return jsonify(message="Could not fetch user info", isLoggedIn=False, isAdmin=False), 400

    user_data = resp.json()
    is_admin = check_admin(user_data)  # replace with your actual admin checking function
    
    return jsonify(user=user_data, isLoggedIn=True, isAdmin=is_admin), 200


def check_admin():
    return true