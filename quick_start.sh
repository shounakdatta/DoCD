# Sample project build script

# Create a logs directory
mkdir -p logs

# Start React client in a background process
cd client
npm run start > ../logs/react-service.log 2>&1 &
REACT_PID=$!
cd ../

# Start Flask app in background
venv/Scripts/activate.bat
cd flask
export FLASK_APP=hello.py
flask run --host=0.0.0.0 > ../logs/flask-service.log 2>&1 &
FLASK_PID=$!

echo "PIDs:" $REACT_PID $FLASK_PID
