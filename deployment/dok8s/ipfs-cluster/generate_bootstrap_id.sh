echo copy the following ID to field bootstrap-peer-id in env-configmap.yaml
echo copy the following encoded private key to field bootstrap-peer-priv-key in secret.yaml
ipfs-key | base64 | base64
