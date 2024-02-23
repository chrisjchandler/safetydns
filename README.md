Stale DNS Proxy Server

The Stale DNS Proxy Server is a robust DNS forwarding service designed to enhance network resilience and reliability. This application acts as an intermediary between clients and upstream DNS servers, caching DNS queries to reduce latency and load on upstream servers. Uniquely, it features the capability to serve stale DNS data during upstream DNS outages, ensuring continuous domain name resolution based on cached data. Additionally, it periodically checks the availability of the upstream DNS server, dynamically adjusting its operation based on the server's status.

Features

DNS Query Caching: Reduces DNS query latency and upstream server load by caching DNS query responses.
Stale Data Serving: Enhances network resilience by serving cached DNS data when the upstream server is unreachable, ensuring continued operation.
Upstream Health Checks: Periodically checks the availability of the upstream DNS server and refreshes the cache with updated data when possible.
Easy Configuration: Offers simple configuration options for specifying the upstream DNS server and tuning cache behavior.
Logging and Monitoring: Provides detailed logs for monitoring the health of the DNS proxy server and the status of DNS queries.

Getting Started

Prerequisites

Go 1.15 or higher

Access to modify local network settings for DNS configuration
sudo privileges (for binding to port 53) See directions below for running without elevated privledges 

Installation

Ensure Go is installed and properly configured on your system.
Download the DNS proxy server source code to your local machine.
Navigate to the source code directory.

Build the application with Go:

go build -o dnsproxy
Running the Application
To start the DNS proxy server, run:

sudo ./dnsproxy

Note: Binding to port 53 requires elevated privileges; hence, sudo is necessary.

safetydns has no health check to determine if the name is resolvable again at the upstream resolver, cache times out at 24 hours and either the DNS is working or you're dead in the water

safetydnshc has a heatlh check that tests the resolver every 5 minutes to see if it's back so you can stop using the srvstale cache

Configuration

The application can be configured by editing specific variables within the source code:

upstreamDNS: Specify the address of the upstream DNS server you wish to use (e.g., "8.8.8.8:53" for Google DNS).

Health check interval and stale data serving period are adjustable within the scheduleHealthCheck and DNSCache functionalities.

Usage

After starting the DNS proxy server, configure your client device or router's DNS settings to point to the IP address where the DNS proxy server is running. This setup redirects all DNS queries through the proxy, utilizing its caching and stale serving capabilities.

Contributing

Contributions to the Stale DNS Proxy Server are welcome. Feel free to fork the repository, make your changes, and submit a pull request. We appreciate contributions that improve the application's functionality, performance, or reliability.


See [Deploymentguide.md] [https://github.com/chrisjchandler/safetydns/blob/main/Deploymentguide.md] for additional deployment details including docker & kubernetes 
