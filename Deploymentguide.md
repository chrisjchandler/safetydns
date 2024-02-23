DNS Proxy Server Deployment Guide

This guide provides instructions for deploying the DNS Proxy Server application using Docker and Kubernetes. The application enhances network resilience by serving stale DNS data during upstream DNS outages and incorporates periodic health checks for upstream DNS availability.

Docker Deployment

Building the Docker Image
Clone the repository and navigate to the directory containing the Dockerfile.

Build the Docker image with the following command, replacing your-image-name with your preferred image name:

docker build -t your-image-name:latest .

Verify the image is built successfully by listing all Docker images:

docker images

Running the Image Locally
Run your Docker image locally to test its functionality:

docker run -d --name dnsproxy-server -p 53:53/udp -p 53:53/tcp your-image-name:latest
This command runs the DNS proxy server in detached mode, mapping port 53 on both UDP and TCP from the container to port 53 on the host.

Kubernetes Deployment

Prerequisites

A Kubernetes cluster
kubectl configured to communicate with your cluster

Docker image pushed to a container registry accessible by your Kubernetes cluster

Deployment and Service Configuration
Create a deployment.yaml file for the DNS proxy server with the following content, adjusting the image field to point to your Docker image:

apiVersion: apps/v1
kind: Deployment
metadata:
  name: dnsproxy-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dnsproxy
  template:
    metadata:
      labels:
        app: dnsproxy
    spec:
      containers:
      - name: dnsproxy
        image: your-image-name:latest
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            add:
            - NET_BIND_SERVICE
          runAsNonRoot: true
        ports:
        - containerPort: 53
          name: dns
          protocol: UDP
        - containerPort: 53
          name: dns-tcp
          protocol: TCP
Create a service.yaml file to expose the DNS proxy server within your cluster:


apiVersion: v1
kind: Service
metadata:
  name: dnsproxy-service
spec:
  selector:
    app: dnsproxy
  ports:
    - protocol: UDP
      port: 53
      targetPort: 53
    - protocol: TCP
      port: 53
      targetPort: 53
  type: ClusterIP
Deploying to Kubernetes
Deploy the DNS proxy server and service to your Kubernetes cluster:

kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

Verifying the Deployment
Check the status of the deployment and service:

kubectl get deployments
kubectl get services
Ensure the dnsproxy-deployment is running and the dnsproxy-service is available.

Additional Information

Custom Configuration: Adjust deployment configurations based on your environment and requirements.

Security Considerations: Review Kubernetes security contexts and Docker security best practices.
Monitoring and Logs: Monitor the application's logs for health check outcomes and DNS query handling.
