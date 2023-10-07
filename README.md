# Webhook

Small project which provides a webhook server that may receive webhook message under `/hook/*`.

Each webhook message's body will be pretty printed to the log output.

The primary use-case for this is within Red Hat Advanced Cluster Security to test audit log notifier integrations.

# Usage

You may start the webhook server by running:
```bash
./webhook serve --port 8080
```

After starting, all incoming requests under `/hook/` will be logged.

## Using with ACS

Deploy the webhook server on the same cluster as Central by running:
```bash
kubectl apply -f ./hack/deployment.yaml
```

This will expose the webhook server under `webhook.stackrox.default.svc:8080`.

For ease of configuration within ACS, a declarative configuration exists that you may apply via:
```bash
kubectl apply -f ./hack/audit-log-config.yaml
```

_Note that this requires the usual setup of declartive config within ACS installation method, adding the config map
`declarative-configuatins` as mount point._

## Next things

- Adding kustomize to generate the deployment manifests on-demand, including the latest image to use.
- Adding metrics for the received audit messages from Central grouped by their interaction and method.
