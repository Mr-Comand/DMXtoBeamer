import webserver
import gloabals

def send_command_a():
    print("statup")
    while True:
        try:
            client= gloabals.clients[int(input("Enter client ID to sent artnet: "))]
            dimmer = input(f"Enter value dimmer for client {client.name}({client.id}): ")
            hueshift = input(f"Enter value hueshift for client {client.name}({client.id}): ")
            animation = input(f"Enter value animation for client {client.name}({client.id}): ")
            fx1 = input(f"Enter value fx1 for client {client.name}({client.id}): ")
            fx2 = input(f"Enter value fx2 for client {client.name}({client.id}): ")
            fx3 = input(f"Enter value fx3 for client {client.name}({client.id}): ")
            fx4 = input(f"Enter value fx4 for client {client.name}({client.id}): ")
            data = [0 for i in range(512)]
            data[client.addresse]=dimmer
            data[client.addresse+1]=hueshift
            data[client.addresse+2]=animation
            data[client.addresse+3]=fx1
            data[client.addresse+4]=fx2
            data[client.addresse+5]=fx3
            data[client.addresse+6]=fx4
            data[client.addresse+7]=127
            data[client.addresse+8]=127
            data[client.addresse+9]=127
            data[client.addresse+10]=64
            client.update_artnet(data)
        except:
            print("An error occurred. Retry from the beginning.")
        
    