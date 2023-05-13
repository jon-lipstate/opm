from flask import Flask
from flask_sqlalchemy import SQLAlchemy

# Create connection to mySQL database
db = SQLAlchemy()

def create_app(test_config=None):
    # create and configure the app

    # Tells the app that configuration files are relative to the instance folder. The instance
    # folder is located outside the flaskr package and can hold local data that shouldn’t be
    # committed to version control, such as configuration secrets and the database file.
    app = Flask(__name__, instance_relative_config=False)

    if test_config == None:
        # Using a production configuration
        app.config.from_object('config.ProdConfig')
    else:
        # Using a development configuration
        app.config.from_object('config.DevConfig')


    # initialize the app with the extension
    db.init_app(app)


    # Uncomment if using routes.py and models.py.
    with app.app_context():
        from . import routes  # Import routes
        db.create_all()  # Create sql tables for our data models

        return app

    # # a simple page that says hello
    # @app.route('/hello')
    # def hello():
    #     return 'Hello, World!'
    #
    # return app