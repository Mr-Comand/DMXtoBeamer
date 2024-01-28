import gloabals
from stupidArtnet import StupidArtnetServer
import webserver


class Client():
    dimmer, hueshift, animation, fx1, fx2, fx3, fx4 = 0, 0, 0, 0, 0, 0, 0
    listener_id = None
    _transform = {"rotateX": 0,
                 "rotateY": 0,
                 "rotateZ": 0,
                 "perspective": 0,
                 "scale": 1 }

    def __init__(self, id,  client_id=None, name=None, addresse=0, universe=None, disabled=False, flipY=False, flipX=False):
        self._id = id
        self._client_id = client_id
        self._name = name
        self._addresse = addresse
        self._universe = universe
        self._disabled = disabled
        self._flipY = flipY
        self._flipX = flipX

    def getAsJson(self):
        return f'("client_id":"{self.client_id}","name":"{self.name}","addresse":"{self.addresse}","universe":"{self.universe}", "disabled":"{self.disabled}", "flipX": "{self.flipX}", "flipY": "{self.flipY}")'.replace("(", "{").replace(")", "}")

    def __str__(self):
        return f"Client(ID={self.id}, client_id={self.client_id}, Name={self.name}, Addresse={self.addresse}, Universe={self.universe} , disabled= {self.disabled}, flipX= {self.flipX}, flipY={self.flipY})"
    def transform(self, property, value):
        self._transform[property]=value
        print(property, value,  self._transform)
        webserver.socketio.emit('chang_transform', {
                                'rotateX': self._transform["rotateX"], 'rotateY': self._transform["rotateY"],'rotateZ': self._transform["rotateZ"],'perspective': self._transform["perspective"],'scale': self._transform["scale"],}, room=self._client_id, namespace="/client")
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
        if value in gloabals.unassigned:
            gloabals.unassigned.remove(value)
        else:
            print(unassigned, "not in list")

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
        # dimmer, hueshift, animation, fx1, fx2
        dimmer, hueshift, animation, fx1, fx2, fx3, fx4 = data[self.addresse:self.addresse+7]
        if dimmer != self.dimmer or hueshift != self.hueshift or animation != self.animation or fx1 != self.fx1 or fx2 != self.fx2 or fx3 != self.fx3 or fx4 != self.fx4:
            self.dimmer = dimmer
            self.hueshift = hueshift
            self.animation = animation
            self.fx1 = fx1
            self.fx2 = fx2
            self.fx3 = fx3
            self.fx4 = fx4
            self.sendToDisplayClient()

    def sendToDisplayClient(self):
        webserver.socketio.emit('update_artnet', {
            'dimmer': self.dimmer,
            'hueshift': self.hueshift,
            'animation': self.animation,
            'fx1': self.fx1,
            'fx2': self.fx2,
            'fx3': self.fx3,
            'fx4': self.fx4
        }, room=self.client_id, namespace="/client")
        print("send update")
