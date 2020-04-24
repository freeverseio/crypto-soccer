echo '***************************************************************************************'
echo copy the following ID to field bootstrap-peer-id in env-configmap.yaml
echo copy the following  private encode it with base64 and then copy it to field bootstrap-peer-priv-key in secret.yaml
echo '***************************************************************************************'
ipfs-key | base64

echo
echo
echo '***************************************************************************************'
echo copy the following encoded hash to cluster-secret field in secret.yaml
echo '***************************************************************************************'
secret=$(od  -vN 32 -An -tx1 /dev/urandom | tr -d ' \n')
encoded_secret=`echo -n $secret | base64`
echo CLUSTER_SECRET=$secret
echo ENCODED_CLUSTER_SECRET=$encoded_secret
