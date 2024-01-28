import webserver
import gloabals

def send_command_a():
    print("statup")
    while True:
        client_id = input("Enter client ID to run JavaScript function (or 'exit' to quit): ")
        if client_id.lower() == 'exit':
            break
        dimmer = input("Enter velue dimmer run for client {}: ".format(client_id))
        hueshift = input("Enter velue hueshift run for client {}: ".format(client_id))
        animation = input("Enter velue animation run for client {}: ".format(client_id))
        fx1 = input("Enter velue fx1 run for client {}: ".format(client_id))
        fx2 = input("Enter velue fx2 run for client {}: ".format(client_id))
        webserver.socketio.emit('update_artnet',{
        'dimmer': dimmer,
        'hueshift': hueshift,
        'animation': animation,
        'fx1': fx1,
        'fx2': fx2
    }, room=client_id,namespace="/client")
    