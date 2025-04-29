# Uniclipboard-share - Secure Network Clipboard Sharing Tool

Uniclipboard-share is a Go-based application that enables secure clipboard sharing between computers on the same local network. It supports both text and image content sharing in real-time, with secret key authentication for enhanced security.

## Features

- Real-time clipboard synchronization across devices
- Secure communication using secret key authentication
- Support for both text and image content
- Automatic server discovery on local network
- WebSocket-based communication for efficient data transfer
- Web interface for secret key management
- Fallback to local server if no existing server is found

## Prerequisites

- Go 1.23.0 or higher
- Linux or Windows both works properly

## Installation

1. Clone the repository:

```bash
git clone https://github.com/mohamidsaiid/uniclipboard.git
cd uniclipboard
```

2. Install dependencies:

```bash
go mod download
go mod tidy
```

3. Create a `.env` file in the project root with the following content:

```env
BASE_URL=192.168.1
PORT=:8080
SECRET_KEY_PORT=:3000
ORIGINAL_SECRET_KEY=your-initial-secret-key
```

Note:

- Adjust the BASE_URL according to your local network configuration
- Choose a different PORT if 8080 is already in use
- SECRET_KEY_PORT is for the web interface to manage secret keys
- ORIGINAL_SECRET_KEY is the initial key used before setting up a custom one

## Usage

1. Run the application:

```bash
go run cmd/main.go
```

2. Initial Setup:

   - When you first run the application, visit `http://localhost:3000/secretkey`
   - Set up your secret key through the web interface
   - Use the same secret key on all devices you want to connect

3. The application will:
   - Start a web server for secret key management
   - Search for existing UniShare servers on your network
   - If no server is found, start a new server
   - Automatically connect to the available server using your secret key
   - Begin synchronizing clipboard content with other authenticated clients

## Security

- All clipboard sharing requires secret key authentication
- Only devices with matching secret keys can share clipboard content
- Secret keys can be updated through the web interface
- Connections are validated during WebSocket handshake

## How It Works

- The application uses a client-server architecture with authentication
- Secret keys are stored locally in a SQLite database
- When started, it searches for existing servers on the local network
- If no server is found, it creates a new one
- All connections are authenticated using the configured secret key
- Clipboard changes are automatically detected and shared
- Supports both text and binary data (images)
- Uses WebSocket for real-time bidirectional communication

## Technical Details

- Uses `gorilla/websocket` for WebSocket communication
- SQLite database for secret key management
- Implements mutex-based synchronization for thread safety
- Automatic server discovery in the IP range xxx.xxx.1.2 to xxx.xxx.1.254
- Clipboard monitoring using `golang.design/x/clipboard`

## Troubleshooting

1. Connection Issues:

   - Ensure all devices use the same secret key
   - Check if the server is accessible on your network
   - Verify your firewall allows connections on the specified ports

2. If the application fails to start:

   - Check if the `.env` file is properly configured
   - Ensure required ports are not in use
   - Verify you have proper permissions for clipboard access

3. For clipboard issues:

   - Check system permissions for clipboard access

4. Secret Key Problems:
   - Make sure the web interface is accessible
   - Check if the SQLite database has proper permissions
   - Verify the secret key matches across all devices

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security Note

While UniShare implements secret key authentication, it's designed for use within trusted local networks. The clipboard data transfer itself is not encrypted beyond the authentication layer.
