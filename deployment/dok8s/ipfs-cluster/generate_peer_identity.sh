echo copy the following ID to field peer-id-n in env-configmap.yaml
echo copy the following encoded private key directly to field peer-priv-key-n into secret.yaml
ipfs-key -type Ed25519  | base64 | base64
