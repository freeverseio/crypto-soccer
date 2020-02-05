```bash
eksctl create cluster --name freeverse-test-fargate --version 1.14 --region eu-west-1 --fargate --profile eaylon-freeverse
```
output:
[ℹ]  eksctl version 0.13.0
[ℹ]  using region eu-west-1
[ℹ]  setting availability zones to [eu-west-1a eu-west-1b eu-west-1c]
[ℹ]  subnets for eu-west-1a - public:192.168.0.0/19 private:192.168.96.0/19
[ℹ]  subnets for eu-west-1b - public:192.168.32.0/19 private:192.168.128.0/19
[ℹ]  subnets for eu-west-1c - public:192.168.64.0/19 private:192.168.160.0/19
[ℹ]  using Kubernetes version 1.14
[ℹ]  creating EKS cluster "freeverse-test-fargate" in "eu-west-1" region with Fargate profile
[ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=eu-west-1 --cluster=freeverse-test-fargate'
[ℹ]  CloudWatch logging will not be enabled for cluster "freeverse-test-fargate" in "eu-west-1"
[ℹ]  you can enable it with 'eksctl utils update-cluster-logging --region=eu-west-1 --cluster=freeverse-test-fargate'
[ℹ]  Kubernetes API endpoint access will use default of {publicAccess=true, privateAccess=false} for cluster "freeverse-test-fargate" in "eu-west-1"
[ℹ]  1 task: { create cluster control plane "freeverse-test-fargate" }
[ℹ]  building cluster stack "eksctl-freeverse-test-fargate-cluster"
[ℹ]  deploying stack "eksctl-freeverse-test-fargate-cluster"
[✔]  all EKS cluster resources for "freeverse-test-fargate" have been created
[✔]  saved kubeconfig as "/Users/eaylon/.kube/config"
[ℹ]  creating Fargate profile "fp-default" on EKS cluster "freeverse-test-fargate"
[ℹ]  created Fargate profile "fp-default" on EKS cluster "freeverse-test-fargate"
[ℹ]  "coredns" is now schedulable onto Fargate
[ℹ]  "coredns" pods are now scheduled onto Fargate
[ℹ]  "coredns" is now scheduled onto Fargate
[ℹ]  "coredns" pods are now scheduled onto Fargate
[ℹ]  kubectl command should work with "/Users/eaylon/.kube/config", try 'kubectl get nodes'
[✔]  EKS cluster "freeverse-test-fargate" in "eu-west-1" region is ready

# get cluster info
```bash
kubectl get clusster-info
```

# deploy freeverse-admin service account
```bash
kubectl apply  -f freeverse-admin-service-account.yaml
```

# get secret token for service account to use when logging in dashboard
```bash
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep freeverse-admin | awk '{print $1}')
```

# install kubernetes dashboard
follow https://docs.aws.amazon.com/eks/latest/userguide/dashboard-tutorial.html
but use the following deployment instead:
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml
```
finally run
```bash
kubectl proxy
```
and browse to <http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/overview?namespace=default>
