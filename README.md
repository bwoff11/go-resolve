# Go-Resolve DNS Server
Go-Resolve is a high-performance DNS server developed in Go, designed with simplicity and adaptability in mind. It utilizes a single YAML configuration file, making it perfectly suited for Kubernetes or serverless deployments where streamlined configuration and deployment processes are essential.

## Features
- Caching Mechanism: Enhances performance by caching DNS query responses, reducing latency and upstream server load.
- Blocklisting: Offers the ability to block domains using customizable blocklists to improve network security.
- Custom Local Records: Allows defining custom DNS records for local network overrides.
- Prometheus Metrics Integration: Provides comprehensive metrics on DNS query processing, cache performance, and blocklist efficiency, facilitating easy monitoring.
- UDP and TCP Support: Handles DNS queries over both UDP and TCP protocols, ensuring compatibility with various clients and network configurations.

## Configuration
The DNS server is configurable via a JSON configuration file, allowing you to specify upstream DNS servers, cache settings, blocklists, and local DNS records. See the config.example.json for a template and documentation on each setting.

## Installation
- Clone the Repository: git clone https://github.com/yourrepo/go-resolve.git
- Build the Server: Navigate to the project directory and run go build -o go-resolve.
- Configure: Copy config.example.json to config.json and adjust the settings according to your network environment.

## Usage
To start the DNS server, run ```./go-resolve```

Ensure your DNS client system or router is configured to use the server's IP address as the DNS server.

## Prometheus Metrics
The DNS server exposes Prometheus-compatible metrics at /metrics endpoint. Configure your Prometheus server to scrape metrics from this endpoint to monitor performance and query statistics.

## Contributing
Contributions to Go-Resolve are welcome! Whether it's feature enhancements, bug fixes, or documentation improvements, please feel free to fork the repository and submit a pull request.

## Development Setup
- Ensure you have Go installed (version 1.15 or later recommended).
- Install necessary Go packages with go get -d ./....
- Use go test to run unit tests.

## Submitting Contributions
- Fork the repository.
- Create a new branch for your feature or fix.
- Commit your changes with clear, descriptive messages.
- Push your branch and submit a pull request against the main branch.

## License
This project is licensed under the MIT License - see the LICENSE file for details.