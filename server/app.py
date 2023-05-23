from dotenv import load_dotenv
import os
load_dotenv()
if os.getenv('FLASK_ENV') == 'development':
    os.environ['OAUTHLIB_INSECURE_TRANSPORT'] = '1'
from flask import Flask, session, redirect
from flask_cors import CORS
from flask_dance.contrib.github import make_github_blueprint, github
from flask_dance.consumer import oauth_authorized
from apis.search import search_bp
from apis.details import details_bp
from apis.auth import auth_bp

app = Flask(__name__)
CORS(app)
app.secret_key = os.getenv("SECRET_KEY")
app.config["GITHUB_OAUTH_CLIENT_ID"] = os.getenv("GITHUB_OAUTH_CLIENT_ID")
app.config["GITHUB_OAUTH_CLIENT_SECRET"] = os.getenv("GITHUB_OAUTH_CLIENT_SECRET")

github_blueprint = make_github_blueprint(
    client_id=os.getenv("GITHUB_OAUTH_CLIENT_ID"),
    client_secret=os.getenv("GITHUB_OAUTH_CLIENT_SECRET")
    )
app.register_blueprint(github_blueprint, url_prefix="/auth")
app.register_blueprint(search_bp)
app.register_blueprint(details_bp)
app.register_blueprint(auth_bp, url_prefix="/auth")



@oauth_authorized.connect_via(github_blueprint)
def github_logged_in(blueprint, token):
    resp = github.get("/user")
    assert resp.ok
    user = resp.json()
    session["login"] = user["login"]
    session["token"] = token
    print(session.keys())
    return redirect(os.getenv("CLIENT_HOST"))

# generate flask's SECRET_KEY
# import secrets
# print(secrets.token_hex(16))

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0')
