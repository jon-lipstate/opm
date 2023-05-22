from dotenv import load_dotenv
import os
load_dotenv()
if os.getenv('FLASK_ENV') == 'development':
    os.environ['OAUTHLIB_INSECURE_TRANSPORT'] = '1'

from flask import Flask
from flask_cors import CORS
from flask_dance.contrib.github import make_github_blueprint
from apis.search import search_bp
from apis.details import details_bp
from apis.login import login_bp
from apis.user import user_bp

app = Flask(__name__)
CORS(app)
app.secret_key = os.getenv("SECRET_KEY")
app.config["GITHUB_OAUTH_CLIENT_ID"] = os.getenv("GITHUB_OAUTH_CLIENT_ID")
app.config["GITHUB_OAUTH_CLIENT_SECRET"] = os.getenv("GITHUB_OAUTH_CLIENT_SECRET")

github_blueprint = make_github_blueprint()
app.register_blueprint(github_blueprint, url_prefix="/github-callback")
app.register_blueprint(search_bp)
app.register_blueprint(details_bp)
app.register_blueprint(user_bp)
app.register_blueprint(login_bp, url_prefix="/login")


# generate flask's SECRET_KEY
# import secrets
# print(secrets.token_hex(16))


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0')
