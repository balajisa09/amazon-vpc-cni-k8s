#!/usr/bin/env bash

# script to run integration tests (no cluster creation & deletion/installing addons)
# use case: run script after CNI images are updated to the image to be verified (see update-cni-images.sh)

# CLUSTER_NAME: name of the cluster to run the test
# VPC_ID: cluster VPC ID
# REGION: default us-west-2
# KUBE_CONFIG_PATH: path to the kubeconfig file, default ~/.kube/config
# NG_LABEL_KEY: nodegroup label key, default "kubernetes.io/os"
# NG_LABEL_VAL: nodegroup label val, default "linux"
# CNI_METRICS_HELPER: cni metrics helper image tag, default "602401143452.dkr.ecr.us-west-2.amazonaws.com/cni-metrics-helper:v1.11.4"

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
INTEGRATION_TEST_DIR="$SCRIPT_DIR/../test/integration"

source "$SCRIPT_DIR"/lib/cluster.sh
source "$SCRIPT_DIR"/lib/integration.sh

function run_integration_test() {
  : "${NG_LABEL_KEY:=kubernetes.io/os}"
  : "${NG_LABEL_VAL:=linux}"
  TEST_RESULT=success
  echo "Running cni integration tests"
  START=$SECONDS
  cd $INTEGRATION_TEST_DIR/cni && CGO_ENABLED=0 ginkgo $EXTRA_GINKGO_FLAGS -v -timeout 60m --no-color --fail-on-pending -- --cluster-kubeconfig="$KUBE_CONFIG_PATH" --cluster-name="$CLUSTER_NAME" --aws-region="$REGION" --aws-vpc-id="$VPC_ID" --ng-name-label-key="$NG_LABEL_KEY" --ng-name-label-val="$NG_LABEL_VAL" || TEST_RESULT=fail
  echo "cni test took $((SECONDS - START)) seconds."
  echo "Running ipamd integration tests"
  START=$SECONDS
  # NOTE: skipping ipamd_event_test.go until it can be triaged further
  cd $INTEGRATION_TEST_DIR/ipamd && CGO_ENABLED=0 ginkgo $EXTRA_GINKGO_FLAGS --skip-file=ipamd_event_test.go -v -timeout 90m --no-color --fail-on-pending -- --cluster-kubeconfig="$KUBE_CONFIG_PATH" --cluster-name="$CLUSTER_NAME" --aws-region="$REGION" --aws-vpc-id="$VPC_ID" --ng-name-label-key="$NG_LABEL_KEY" --ng-name-label-val="$NG_LABEL_VAL" || TEST_RESULT=fail
  echo "ipamd test took $((SECONDS - START)) seconds."

  : "${CNI_METRICS_HELPER:=602401143452.dkr.ecr.us-west-2.amazonaws.com/cni-metrics-helper:v1.11.4}"
  REPO_NAME=$(echo $CNI_METRICS_HELPER | cut -d ":" -f 1)
  TAG=$(echo $CNI_METRICS_HELPER | cut -d ":" -f 2)
  echo "Running cni-metrics-helper image($CNI_METRICS_HELPER) tests"
  START=$SECONDS
  cd $INTEGRATION_TEST_DIR/metrics-helper && CGO_ENABLED=0 ginkgo $EXTRA_GINKGO_FLAGS -v -timeout 15m --no-color --fail-on-pending -- --cluster-kubeconfig="$KUBE_CONFIG_PATH" --cluster-name="$CLUSTER_NAME" --aws-region="$REGION" --aws-vpc-id="$VPC_ID" --ng-name-label-key="$NG_LABEL_KEY" --ng-name-label-val="$NG_LABEL_VAL" --cni-metrics-helper-image-repo=$REPO_NAME --cni-metrics-helper-image-tag=$TAG || TEST_RESULT=fail
  echo "cni-metrics-helper test took $((SECONDS - START)) seconds."
  if [[ "$TEST_RESULT" == fail ]]; then
      echo "Integration test failed."
      exit 1
  fi
  echo "Integration tests completed successfully!"
}

if [[ -n "${ENDPOINT}" ]]; then
  ENDPOINT_FLAG="--endpoint $ENDPOINT"
fi

echo "Running release tests on cluster: $CLUSTER_NAME in region: $REGION"

load_cluster_details
START=$SECONDS
run_integration_test

echo "Completed running all tests in $((SECONDS - START)) seconds."
