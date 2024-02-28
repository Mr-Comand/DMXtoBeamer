import gloabals
from stupidArtnet import StupidArtnetServer
import webserver


class Client():
    
    

    def __init__(self, id,  client_id=None, name=None, addresse=0, universe=None, disabled=False, flipY=False, flipX=False):
        self._id = id
        self._client_id = client_id
        self._name = name
        self._addresse = addresse
        self._universe = universe
        self._disabled = disabled
        self._flipY = flipY
        self._flipX = flipX
        self._transform = {
                 "rotateX": 0,
                 "rotateY": 0,
                 "rotateZ": 0,
                 "perspective": 500,
                 "scale": 1,
                 "translateX":  0,
                 "translateY":  0,
                 "scaleFull": 1
                 }
        self.dimmer, self.hueshift, self.animation, self.fx1, self.fx2, self.fx3, self.fx4, self.pan, self.tilt, self.rotate, self.zoom = 0, 0, 0, 0, 0, 0, 0, 127, 127, 127, 64
        self.listener_id = None
    
    _calibratemode = False
    def getAsJson(self):
        return f'("client_id":"{self.client_id}","name":"{self.name}","addresse":"{self.addresse}","universe":"{self.universe}", "disabled":"{self.disabled}", "flipX": "{self.flipX}", "flipY": "{self.flipY}")'.replace("(", "{").replace(")", "}")

    def __str__(self):
        return f"Client(ID={self.id}, client_id={self.client_id}, Name={self.name}, Addresse={self.addresse}, Universe={self.universe} , disabled= {self.disabled}, flipX= {self.flipX}, flipY={self.flipY})"
    def transform(self, property, value):
        self._transform[property]=value
        print(property, value,  self._transform)
        webserver.socketio.emit('chang_transform', {
                                'rotateX': self._transform["rotateX"], 'rotateY': self._transform["rotateY"],'rotateZ': self._transform["rotateZ"],'perspective': self._transform["perspective"],'scale': self._transform["scale"], 'scaleFull': self._transform["scaleFull"], "translateX": self._transform["translateX"], "translateY": self._transform["translateY"], "calibratemode":self.calibratemode}, room=self._client_id, namespace="/client")
    
    @property
    def calibratemode(self):
        return self._calibratemode

    @calibratemode.setter
    def calibratemode(self, value):
        self._calibratemode = value
        webserver.socketio.emit('chang_transform', {
                                'rotateX': self._transform["rotateX"], 'rotateY': self._transform["rotateY"],'rotateZ': self._transform["rotateZ"],'perspective': self._transform["perspective"],'scale': self._transform["scale"], 'scaleFull': self._transform["scaleFull"],  "translateX": self._transform["translateX"], "translateY": self._transform["translateY"], "calibratemode":self.calibratemode}, room=self._client_id, namespace="/client")
        
    @property
    def id(self):
        return self._id

    @property
    def client_id(self):
        return self._client_id

    @client_id.setter
    def client_id(self, value):
        if self._client_id != "" and self._client_id != "None" and self._client_id != None:
            gloabals.unassigned.append(self._client_id)
            webserver.socketio.emit('chang_settings', {
                                    'flipY': "false", 'flipX': "false", 'disabled': "true"}, room=self._client_id, namespace="/client")
            
        self._client_id = value
        webserver.socketio.emit('chang_settings', {
                                'flipY': self._flipY, 'flipX': self._flipX, 'disabled': self._disabled}, room=self._client_id, namespace="/client")
        webserver.socketio.emit('chang_transform', {
                                'rotateX': self._transform["rotateX"], 'rotateY': self._transform["rotateY"],'rotateZ': self._transform["rotateZ"],'perspective': self._transform["perspective"],'scale': self._transform["scale"], 'scaleFull': self._transform["scaleFull"],  "translateX": self._transform["translateX"], "translateY": self._transform["translateY"], "calibratemode":self.calibratemode}, room=self._client_id, namespace="/client")
        if value in gloabals.unassigned:
            gloabals.unassigned.remove(value)
        else:
            print(globals.unassigned, "not in list")

        self.sendToDisplayClient()

    @property
    def name(self):
        return self._name

    @name.setter
    def name(self, value):
        self._name = value

    @property
    def addresse(self):
        return self._addresse

    @addresse.setter
    def addresse(self, value):
        if value == None:
            self._addresse = 0
            return
        self._addresse = max(min(int(value), 506), 0)

    @property
    def universe(self):
        return self._universe

    @universe.setter
    def universe(self, value):
        if value == self.universe:
            return
        if value == None or int(value) < 0:
            if self.listener_id == None:
                None
            else:
                gloabals.ArtnetClient.delete_listener(self.listener_id)
                self._universe = None
        else:
            value = int(value)
            if self.listener_id == None:
                None
            else:
                gloabals.ArtnetClient.delete_listener(self.listener_id)
            self._universe = value
            self.listener_id = gloabals.ArtnetClient.register_listener(
                universe=self._universe, callback_function=self.update_artnet)

        self._universe = value

    @property
    def disabled(self):
        return self._disabled

    @disabled.setter
    def disabled(self, value):
        self._disabled = value
        webserver.socketio.emit('chang_settings', {
                                'flipY': self._flipY, 'flipX': self._flipX, 'disabled': self._disabled}, room=self._client_id, namespace="/client")

    @property
    def flipY(self):
        return self._flipY

    @flipY.setter
    def flipY(self, value):
        self._flipY = value
        webserver.socketio.emit('chang_settings', {
                                'flipY': self._flipY, 'flipX': self._flipX, 'disabled': self._disabled}, room=self._client_id, namespace="/client")

    @property
    def flipX(self):
        return self._flipX

    @flipX.setter
    def flipX(self, value):
        self._flipX = value
        webserver.socketio.emit('chang_settings', {
                                'flipY': self._flipY, 'flipX': self._flipX, 'disabled': self._disabled}, room=self._client_id, namespace="/client")

    def update_artnet(self, data):
        # self.dimmer, self.hueshift, self.animation, self.fx1, self.fx2, self.fx3, self.fx4, self.pan, self.tilt, self.rotate, self.zoom = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
        dimmer, hueshift, animation, fx1, fx2, fx3, fx4, pan, tilt, rotate, zoom = data[self.addresse:self.addresse+11]
        if dimmer != self.dimmer or hueshift != self.hueshift or animation != self.animation or fx1 != self.fx1 or fx2 != self.fx2 or fx3 != self.fx3 or fx4 != self.fx4 or pan != self.pan or tilt != self.tilt or rotate != self.rotate or zoom != self.zoom:
            self.dimmer = dimmer
            self.hueshift = hueshift
            self.animation = animation
            self.fx1 = fx1
            self.fx2 = fx2
            self.fx3 = fx3
            self.fx4 = fx4
            self.pan = pan
            self.tilt = tilt
            self.rotate = rotate
            self.zoom = zoom
            self.sendToDisplayClient()

    def sendToDisplayClient(self):
        webserver.socketio.emit('update_artnet', {
            'dimmer': self.dimmer,
            'hueshift': self.hueshift,
            'animation': self.animation,
            'fx1': self.fx1,
            'fx2': self.fx2,
            'fx3': self.fx3,
            'fx4': self.fx4,
            'pan': self.pan,
            'tilt': self.tilt,
            'rotate': self.rotate,
            'zoom': self.zoom
        }, room=self.client_id, namespace="/client")
        print("send update")
