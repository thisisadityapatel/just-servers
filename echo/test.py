import socket

HOST = '0.0.0.0'
PORT = 10000

# Create a socket object
with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((HOST, PORT))
     
    message = b'Hello, TCP Server!'
    s.sendall(message)
    print(f"Sent: {message.decode()}")

    data = s.recv(1024)
    print(f"Received: {data.decode()}")
