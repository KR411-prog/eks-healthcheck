check_env:
ifndef ACCOUNT_ID
	$(error ACCOUNT_ID is undefined)
endif
ifndef CLUSTER
	$(error CLUSTER is undefined)
endif
ifndef REGION
 	$(error REGION is undefined)
endif

gomod:
	go mod download

testeks:  setenv  gomod
	set -e;
	clusterName=${CLUSTER} region=${REGION} account_id=${ACCOUNT_ID} go test -v -run TestEKSComplete -timeout 30m


testtag: setenv gomod
	set -e;
	clusterName=${CLUSTER} region=${REGION} account_id=${ACCOUNT_ID} go test -v -run TestTagging -timeout 30m

.PHONY: check_env setenv gomod testeks testtag
