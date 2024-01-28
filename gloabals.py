from pprint import pprint
from Client import Client
from stupidArtnet import StupidArtnetServer

sendConfigToAllCliets = None
ArtnetClient = StupidArtnetServer()





# unassigned ids
unassigned: list[str] = []

available: list[str] = []
# map of internal IDs to ther client_ids (1=>"SrPo2guNYcrdp00BAAAB")
clients: dict[int, Client] = {}


def printall():
    """Print all the known data"""
    pprint(unassigned, clients, names, addresses, universes)


def getClientByName(name):
    """Get a Client ID by its name. Returns None if not found."""
    for id in clients.keys():
        if clients[id].name == name:
            return clients[id]
    return None


def getClientByCleintId(client_id):
    """Get a Client ID by its name. Returns None if not found."""
    for id in clients.keys():
        if clients[id].client_id == client_id:
            return clients[id]
    return None


def getasJson():
    out = "{\"unassigned\":[ "
    for id in unassigned:
        out += "\""+id+"\","
    out = out[:-1]+"], \"clients\":{ "
    for key in clients.keys():
        out += f'"{key}": {clients[key].getAsJson()},'
    out = out[:-1]+"} }"
    pprint(out)
    return out
    sendConfigToAllCliets()


def addinternalClinet(id):
    """Adds an internal client number that is not yet assigned"""
    if id not in clients:
        clients[id] = Client(id)
    sendConfigToAllCliets()


def removeInternalClient(id):
    if clients[id].client_id != "" and clients[id].client_id != "None" and clients[id].client_id != None:
        unassigned.append(clients[id].client_id)
    clients.pop(id)
    sendConfigToAllCliets()


def disconnectCleint(client_id):
    if client_id in unassigned:
        unassigned.remove(client_id)
    else:
        getClientByCleintId(client_id)._client_id = None
    sendConfigToAllCliets()
