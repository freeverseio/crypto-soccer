echo copy the following ID to field peer-id-n in env-configmap.yaml
echo "copy the following private key and run echo -n <key> | base64 and then copy it into secret"
ipfs-key -type Ed25519 | base64
