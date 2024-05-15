kubectl scale deployment --replicas=0  -n kube-system cilium-operator
make kind-build-image-operator
minikube image rm localhost:5000/cilium/operator-generic:local
minikube image load localhost:5000/cilium/operator-generic:local
kubectl scale deployment --replicas=1 -n kube-system cilium-operator

# helm repo add cilium https://helm.cilium.io/

# helm install cilium cilium/cilium --version 1.15.4 \
#   --namespace kube-system \
    # --set kubeProxyReplacement=true \
    # --set gatewayAPI.enabled=true
  #   --set hubble.relay.enabled=true \
  #  --set hubble.ui.enabled=true
