# mono-client
> single-client / single-device per account

This app is built with golang showing how to allow only one device linked to an account. This tracks the clients IP address and ensures that only one device can be authorized at a time for the selected account. The authotized device's IP addresses are cached in redis memory with their corresponding token and userID. 

A device can signin to an account if it has no device linked to that account. if there exists another client or device connected to that account, the server prevents it from signing in to that account (meaning there shouldn't exist any device connected to that account before signing in).

### Tech stack: 
Golang, Redis, Pg


### Local starter

There are two starter configs in the starters, starter-docker.sh and starter-podman.sh (if you are using Podman daemon).

Make these shell scripts executable using 

bash ```
chmod +x ./starters/starter-*.sh

./starters/starter-*
```


## Kubernetes

A starter Kubernetes manifest is available at `k8s/mono-client.yaml`.

It includes:
- Namespace
- Postgres + Redis Deployments/Services/PVCs
- App Deployment/Service
- Secrets for Postgres credentials and app `.env` values

### Deploy

```bash
kubectl apply -f k8s/mono-client.yaml
```

> Note: the app deployment uses image `mono-client:latest` by default. Build and push that image (or edit the manifest to your registry image) before deploying.
