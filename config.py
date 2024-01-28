from webserver import app, socketio
import webserver
import gloabals
from flask import render_template, request, jsonify

config_clients=[]
@app.route('/')
def index():
    return render_template('index.html')

@socketio.on('connect', namespace='/config')
def connect_config():
     print("Client connected config")
     # emit an event with data to the client
     global clients
     config_clients.append(request.sid)
     
def sendConfigToAllCliets():
    # sen a bradcas to all clients so thy can  update their configuration   
    for sessionId in config_clients:
        socketio.emit('update_config', "configs", room=sessionId, namespace='/config')    
gloabals.sendConfigToAllCliets = sendConfigToAllCliets
# handle events from the client 
@socketio.on('disconnect', namespace='/config')
def disconnect():
    print("config Client disconnected", request.sid)
    sendConfigToAllCliets()
    
@app.route('/api/config')
def getData():
    return gloabals.getasJson()
    
@app.route('/api/add/<int:client_id>' ,methods=["POST"])
def addclient(client_id):    
    gloabals.addinternalClinet(client_id)
    sendConfigToAllCliets()
    return gloabals.getasJson()

@app.route('/api/remove/<int:client_id>' ,methods=["DELETE"])
def removeclient(client_id):    
    gloabals.removeInternalClient(client_id)
    sendConfigToAllCliets()
    return gloabals.getasJson()


@app.route('/api/update/<int:id>/<string:field>/<string:value>', methods=['PUT'])
def updateConfig(id, field, value):
    match field:
     case 'name':
         print('name')
         gloabals.clients[id].name=value
     case 'address':
         print('address')
         gloabals.clients[id].addresse=value
     case 'universe':
         print('universe')
         gloabals.clients[id].universe=value
     case 'clientId':
         gloabals.clients[id].client_id=value
     case 'flipY':
         gloabals.clients[id].flipY = (value == 'true')
     case 'flipX':
         gloabals.clients[id].flipX = (value=='true')   
     case 'disabled':
         gloabals.clients[id].disabled = (value=='true')
    sendConfigToAllCliets()
    return gloabals.getasJson()

@app.route('/api/identify/<string:client_id>' ,methods=["POST"])
def identifyclient(client_id):    
    webserver.socketio.emit('run_function', {'function_name': "identify"}, room=client_id,namespace="/client")
    return gloabals.getasJson()