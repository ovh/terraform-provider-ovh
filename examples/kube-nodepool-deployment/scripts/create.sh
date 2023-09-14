#!/bin/bash

exec &> logs

NB_CLUSTER=1
let "MAX=NB_CLUSTER-1"

terraform init

START=$(date +%s)
date

terraform plan -var="nb_cluster=${NB_CLUSTER}"

terraform apply -var="nb_cluster=${NB_CLUSTER}" -auto-approve

date

apply_deploy() {
    KUBECONFIG_FILE="my-kube-cluster-$1.yml"
    kubectl --kubeconfig=$KUBECONFIG_FILE get nodes
    kubectl --kubeconfig=$KUBECONFIG_FILE apply -f hello.yaml -n default
    kubectl --kubeconfig=$KUBECONFIG_FILE get all -n default
    kubectl --kubeconfig=$KUBECONFIG_FILE get services -n default -l app=hello-world

    ip=""
    while [ -z $ip ]; do
        echo "Waiting for external IP"
        ip=$(kubectl --kubeconfig=${KUBECONFIG_FILE} -n default get service hello-world -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
        [ -z "$ip" ] && sleep 10
    done
    echo 'Found external IP: '$ip
    export APP_IP=$ip
    echo $APP_IP
    curl $APP_IP
}

wait 

# Deploy on each cluster as jobs
for ((i=0;i<=$MAX;i++)); 
do
    apply_deploy $i &
done

# Wait for jobs to finish
wait

date
END=$(date +%s)

echo $((END-START)) | awk '{printf "%d:%02d:%02d", $1/3600, ($1/60)%60, $1%60}'
