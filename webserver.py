from flask import Flask, render_template, request, send_from_directory
from flask_socketio import SocketIO, emit
import gloabals
import os

app = Flask(__name__)
socketio = SocketIO(app)

# List to store client connections
connected_clients = []

# Route to serve the HTML page
@app.route('/client')
def client():
    return render_template('client.html')

# WebSocket event to notify the client to run a JavaScript function
@socketio.on('run_function', namespace='/client')
def run_function(data):
    function_name = data.get('function_name', '')
    socketio.emit('run_function', {'function_name': function_name})



# Event handler when a new client connects
@socketio.on('connect', namespace='/client')
def handle_connect():
    client_id = request.sid  # Change this line to use sid directly
    gloabals.unassigned.append(client_id)
    gloabals.sendConfigToAllCliets()
    socketio.emit('chang_settings', {'flipY': "false", 'flipX': "false", 'disabled': "true"}, room=client_id, namespace="/client")
    print("Client {} connected.".format(client_id))

# Event handler when a client disconnects
@socketio.on('disconnect', namespace='/client')
def handle_disconnect():
    client_id = request.sid  # Change this line to use sid directly
    gloabals.disconnectClient(client_id)
    print("Client {} disconnected.".format(client_id))


@app.route('/Monitor-Calibration.avif')
def testimage():
    #return image
    return send_from_directory(os.path.join(app.root_path),'Monitor-Calibration.avif')

@app.route('/animations/<int:id>')
def getanimations(id):
    if id<0 or id>255:
        return "Invalid animation ID"
    else:
        return send_from_directory(os.path.join(app.root_path, 'animations'), str(id)+".html")
    #return image
@app.route('/static/string:name>')
def staticfiles(name):
    #return image
    return send_from_directory(os.path.join(app.root_path,'static'),name)
@app.route('/animationassets/<string:id>/<string:name>')
def animationassets(id, name):
    #return image
    return send_from_directory(os.path.join(app.root_path,'animationassets/'+str(id)),name)