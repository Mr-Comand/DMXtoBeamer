
import threading
import webserver, artnet, gloabals, config, maunual


if __name__ == '__main__':
    # Create a thread for the send_command function
    command_thread = threading.Thread(target=artnet.send_command_a, daemon=True)
    command_thread.start()

    # Run the Flask application with SocketIO
    webserver.socketio.run(webserver.app, debug=True, host="0.0.0.0",ssl_context='adhoc', use_reloader=False, port=443)
    # Uncomment the line below if you want to use the default Flask server without SocketIO
    # app.run(debug=True)
