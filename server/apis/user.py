from flask import Blueprint, jsonify
from flask_dance.contrib.github import github

user_bp = Blueprint('user', __name__)

@user_bp.route('/user', methods=['GET'])
def user():
    if not github.authorized:
        return jsonify(user=None), 401  # Not Authorized
    resp = github.get("/user")
    if not resp.ok:
        return jsonify(message="Could not fetch user info"), 400
    return jsonify(user=resp.json()), 200