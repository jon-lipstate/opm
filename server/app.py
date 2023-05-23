from dotenv import load_dotenv
import os
load_dotenv()
if os.getenv('FLASK_ENV') == 'development':
    os.environ['OAUTHLIB_INSECURE_TRANSPORT'] = '1'
from flask import Flask, session, redirect, url_for
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


@app.route("/")
def index():
    link1 = f"<a href=\"{url_for('list_all_urls')}\">{str(url_for('list_all_urls'))}</a>"
    link2 = f"<a href=\"{url_for('github.login')}\">{str(url_for('github.login'))}</a>"

    content = f"Welcome to the Flask backend.  <br>"
    content += f"github.authorized = {github.authorized} <br>"
    content += f"See {link1} for a list of available routes <br>"
    content += f"See {link2} for GitHub OAuth Login <br>"

    return content


@app.route("/list-all-urls")
def list_all_urls():
    content = ""
    for r in app.url_map.iter_rules():
        content += f"{str(r.methods)} {str(r)} <br>"

    return content


@app.route("/inspect-session")
def inspect_session():
    content = "Contents of Flask session: <br>"
    for key, val in session.items():
        content += f"  {key}: {val} <br>"
    return content


@app.route("/login-github")
def login():
    if not github.authorized:
        print(f"Redirect to: {url_for('github.login')}")
        return redirect(url_for('github.login'))

    return redirect(url_for('index'))

# generate flask's SECRET_KEY
# import secrets
# print(secrets.token_hex(16))


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0')
