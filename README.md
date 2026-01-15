## Nocta

A lightweight GUI application for monitoring and managing open ports on your system.
Replace repetitive terminal commands like `ss`, `netstat`, or `lsof`.

## Features

- **Port Monitoring**: Real-time monitoring of TCP/UDP ports
- **Process Details**: View detailed information about processes using ports
- **Port Management**: Terminate processes associated with specific ports
- **Search & Filter**: Search by port number, process name, or filter by protocol
- **Configurable**: YAML-based configuration for customization
- **Export**: Export port data to CSV/JSON formats
- **Cross-platform**: Built with Fyne for Linux, macOS, and Windows

## Project Structure

```
nocta/
├── cmd/nocta/           # Application entry point
├── internal/            # Private application code
│   ├── config/         # Configuration management
│   ├── gui/            # GUI components and application
│   ├── logger/         # Structured logging
│   ├── models/         # Data models and types
│   ├── scanner/        # Port scanning logic
│   └── service/        # Business logic layer
├── pkg/                # Public libraries
├── configs/            # Default configuration files
└── README.md
```

## Installation

### From Source

1. Clone the repository:

```bash
git clone https://github.com/yofabr/nocta.git
cd nocta
```

2. Build the application:

```bash
go build ./cmd/nocta
```

3. Run the application:

```bash
./nocta
```

### From Release (Coming Soon)

You can install Nocta by downloading the latest release tarball:

1. Download the tar file:

```bash
wget https://github.com/yofabr/nocta/releases/download/v1.0.0/nocta.tar.xz
```

2. Extract the archive:

```bash
tar -xf nocta.tar.xz
```

3. Run the installer

```bash
make install
```

4. Enter the extracted directory and run the application:

```bash
cd nocta
./nocta
```

> Note: Make sure the `nocta` file is executable. If not, run `chmod +x nocta`.

## Configuration

Nocta uses YAML configuration stored at `~/.config/nocta/config.yaml`. The default configuration includes:

```yaml
gui:
  width: 900
  height: 520
  split: 0.35
refresh:
  interval: 30
  auto: false
filter:
  default_protocol: ALL
  show_closed: false
notifications:
  enabled: false
  on_new: true
  on_close: true
```

## Usage

1. **View Ports**: Launch Nocta to see all active network ports
2. **Search**: Use the search bar to filter by port number or process name
3. **Protocol Filter**: Filter ports by TCP, UDP, or view all
4. **Port Details**: Click on any port to see detailed information
5. **Terminate**: Use the Terminate button to stop processes (when available)
6. **Refresh**: Click the refresh button or enable auto-refresh for updates

## Screenshot

![Screenshot](./screenshots/main.png)

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
