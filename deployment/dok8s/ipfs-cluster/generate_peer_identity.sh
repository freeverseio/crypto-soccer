echo copy the following ID in env_configmap.yaml
echo copy the following private key directly into secret.yaml
ipfs-key -type Ed25519  | base64 | base64
