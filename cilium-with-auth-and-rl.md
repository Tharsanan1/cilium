

# Cilium Installation and Configuration Guide

This guide provides instructions on installing Cilium with authentication and rate-limiting capabilities 

## Prerequisites

- Kubernetes cluster setup.
- Helm installed.
- Make utility installed.

## Installation Steps

1. Clone the repository

2. Install Cilium using Helm with authentication and rate-limiting:

   ```bash
   make install-cilium-with-auth-rl
   ```

3. Wait for deployments to be ready before proceeding.
4. If you are using cloud clusters (AKS, EKS, GKE) skip this step.
    - If you are using minikube run this command to expose loadbalancer ips to local network. `minikube tunnel`
    - Or else run `install-metallb`

## Setup Authentication

1. Apply resources and security policy:

   ```bash
   make setup-with-auth
   ```

2. Wait for deployments to be ready.

## Test Without Authentication Header

```bash
make test-without-auth-header
```

Gateway will respond with 401 as we have not provided a valid token.

## Test With Authentication Header

```bash
make test-with-auth-header
```

You should see a 200 response. 

## Setup Rate Limiting

1. Apply resources and rate limit policy:

   ```bash
   make setup-with-rate-limit
   ```

2. Wait for deployments to be ready.

## Test ratelimit

- Test rate limiting with custom headers:

  ```bash
  make test-with-rate-limit-custom-headers
  ```

- Test simple rate limiting:

  ```bash
  make test-with-simple-rate-limit
  ```