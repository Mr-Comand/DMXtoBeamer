import json
from webserver import app, socketio
import webserver
import gloabals
from flask import render_template, request, jsonify

# API endpoint to get client IDs


@app.route('/api/client-ids', methods=['GET'])
def get_client_ids():
    return jsonify(list(gloabals.clients.keys()))

# Function to load animations from the file


def load_animations():
    with open('animations.json') as f:
        return json.load(f)

# API endpoint to get animations list


@app.route('/api/animations/list', methods=['GET'])
def get_animations():
    # Load animations from the file on each request
    animations_data = load_animations()

    return jsonify(animations_data)

# API endpoint to update manual control status for a given animation


@app.route('/api/<int:client_id>/animation/<int:animation_id>', methods=['PUT'])
def update_animation_control(client_id, animation_id):
    client = gloabals.clients.get(client_id)
    if not client:
        return jsonify({'error': 'Client not found'}), 404

    client.animation = animation_id
    client.sendToDisplayClient()
    # Update the manual control status in your database or data storage here
    return jsonify({'message': 'Manual control updated successfully'})

# API endpoint to update properties for a given client ID


@app.route('/api/<int:client_id>/properties', methods=['PUT'])
def update_properties(client_id):
    client = gloabals.clients.get(client_id)
    if not client:
        return jsonify({'error': 'Client not found'}), 404

    # Get the JSON data from the request
    data = request.json

    # Set default values or update properties from the request
    client.dimmer = data.get('dimmer', client.dimmer)
    client.hueshift = data.get('hueshift', client.hueshift)
    client.fx1 = data.get('fx1', client.fx1)
    client.fx2 = data.get('fx2', client.fx2)
    client.fx3 = data.get('fx3', client.fx3)
    client.fx4 = data.get('fx4', client.fx4)
    client.pan = data.get('pan', client.pan)
    client.tilt = data.get('tilt', client.tilt)
    client.rotate = data.get('rotate', client.rotate)
    client.zoom = data.get('zoom', client.zoom)

    # Update properties in your database or data storage here
    # For example: db.session.commit()
    client.sendToDisplayClient()
    return jsonify({'message': 'Properties updated successfully'})


# Endpoint to get properties of a specific client
@app.route('/api/<int:client_id>/properties', methods=['GET'])
def get_client_properties(client_id):
    client = gloabals.clients.get(client_id)
    if not client:
        return jsonify({'error': 'Client not found'}), 404

    # Return all client properties
    return jsonify({
        'id': client.id,
        'animation': client.animation,
        'dimmer': client.dimmer,
        'hueshift': client.hueshift,
        'fx1': client.fx1,
        'fx2': client.fx2,
        'fx3': client.fx3,
        'fx4': client.fx4,
        'pan': client.pan,
        'tilt': client.tilt,
        'rotate': client.rotate,
        'zoom': client.zoom
    })

# API endpoint to get detailed client information


@app.route('/api/clients', methods=['GET'])
def get_clients():
    clients_info = []
    for client_id, client in gloabals.clients.items():
        current_animation = next(
            (anim for anim in animations if anim['id']
             == client.animation), None
        )
        clients_info.append({
            'id': client_id,
            'currentAnimation': client.animation,
            'thumbnailUrl': current_animation['thumbnailUrl'] if current_animation else None
        })
    return jsonify(clients_info)


presets = [
    {
        'id': 1,
        'name': 'Bright Scene',
        'animation': 1,
        'config': {
            'dimmer': 255,
            'hueshift': 50,
            'fx1': 100,
            'fx2': 150,
            'fx3': 200,
            'fx4': 50,
            'pan': 90,
            'tilt': 120,
            'rotate': 30,
            'zoom': 10,
        }
    }
]


@app.route('/api/presets', methods=['GET'])
def get_presets():
    return jsonify(presets)


@app.route('/api/presets', methods=['POST'])
def create_preset():
    data = request.json
    new_preset = {
        'id': len(presets) + 1,
        'name': data['name'],
        'animation': data['animation'],
        'config': data['config'],
    }
    presets.append(new_preset)
    return jsonify({'message': 'Preset created successfully', 'preset': new_preset}), 201


@app.route('/api/<int:client_id>/apply-preset/<int:preset_id>', methods=['PUT'])
def apply_preset(client_id, preset_id):
    client = gloabals.clients.get(client_id)
    if not client:
        return jsonify({'error': 'Client not found'}), 404

    preset = next((p for p in presets if p['id'] == preset_id), None)
    if not preset:
        return jsonify({'error': 'Preset not found'}), 404

    client.animation = preset['animation']
    for key, value in preset['config'].items():
        setattr(client, key, value)

    return jsonify({'message': 'Preset applied successfully'})


@app.route('/api/<int:client_id>/save-as-preset', methods=['POST'])
def save_client_as_preset(client_id):
    client = gloabals.clients.get(client_id)
    if not client:
        return jsonify({'error': 'Client not found'}), 404

    # Get the data from the request (client's current configuration)
    data = request.json
    preset_name = data.get('name', f"Preset for {client_id}")

    # Create a new preset with the client's current configuration
    new_preset = {
        'id': len(presets) + 1,
        'name': preset_name,
        'animation': client.animation,  # Include the current animation
        'config': {
            'dimmer': client.dimmer,
            'hueshift': client.hueshift,
            'fx1': client.fx1,
            'fx2': client.fx2,
            'fx3': client.fx3,
            'fx4': client.fx4,
            'pan': client.pan,
            'tilt': client.tilt,
            'rotate': client.rotate,
            'zoom': client.zoom,
        }
    }

    presets.append(new_preset)
    return jsonify({'message': 'Preset saved successfully', 'preset': new_preset}), 201
