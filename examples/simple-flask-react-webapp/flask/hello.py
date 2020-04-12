from flask import Flask, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app)
cors = CORS(app, resorces={
    r"/*": {
        "origins": "*"
    }
})


@app.route('/')
def hello_world():
    return jsonify({"text": "This is Server! I hear you loud and clear."})
