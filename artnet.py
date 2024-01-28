import webserver
import gloabals

def send_command_a():
    print("statup")
    while True:
        client= gloabals.clients[int(input("Enter client ID to sent artnet: "))]
        dimmer = input(f"Enter velue dimmer run for client {client.name}({client.id}): ")
        hueshift = input(f"Enter velue hueshift run for client {client.name}({client.id}): ")
        animation = input(f"Enter velue animation run for client {client.name}({client.id}): ")
        fx1 = input(f"Enter velue fx1 run for client {client.name}({client.id}): ")
        fx2 = input(f"Enter velue fx2 run for client {client.name}({client.id}): ")
        fx3 = input(f"Enter velue fx3 run for client {client.name}({client.id}): ")
        fx4 = input(f"Enter velue fx4 run for client {client.name}({client.id}): ")
        data = [0 for i in range(512)]
        data[client.addresse]=dimmer
        data[client.addresse+1]=hueshift
        data[client.addresse+2]=animation
        data[client.addresse+3]=fx1
        data[client.addresse+4]=fx2
        data[client.addresse+5]=fx3
        data[client.addresse+6]=fx4
        client.update_artnet(data)
        
    