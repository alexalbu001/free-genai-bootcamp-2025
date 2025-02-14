from flask import Flask, g
from flask_cors import CORS

from lib.db import Db

import routes.words
import routes.groups
import routes.study_sessions
import routes.dashboard
import routes.study_activities

def create_app(test_config=None):
    app = Flask(__name__)
    
    if test_config is None:
        app.config.from_mapping(
            DATABASE='words.db',
            CORS_ORIGINS=["http://localhost:8080", "http://localhost:5173"]
        )
    else:
        app.config.update(test_config)
    
    # Initialize database
    app.db = Db(database=app.config['DATABASE'])
    
    # Configure CORS using config
    CORS(app, resources={
        r"/*": {
            "origins": app.config['CORS_ORIGINS'],
            "methods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
            "allow_headers": ["Content-Type"]
        }
    })

    # Close database connection
    @app.teardown_appcontext
    def close_db(exception):
        app.db.close()

    # Register routes
    routes.words.load(app)
    routes.groups.load(app)
    routes.study_sessions.load(app)
    routes.dashboard.load(app)
    routes.study_activities.load(app)
    
    return app

if __name__ == '__main__':
    app = create_app()
    print("Starting Flask server...")
    print("Database path:", app.config['DATABASE'])
    print("CORS settings:", app.config['CORS_ORIGINS'])
    
    try:
        app.run(debug=True, port=3000, host='127.0.0.1')
        print("Server started successfully")
    except Exception as e:
        print(f"Failed to start server: {str(e)}")