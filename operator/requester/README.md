# GKE Test Cluster Requester

This is a simple Go program that serves as a client to the GKE Test Cluster Operator.

I can be use by CI jobs as well as developers.


## Developer Usage

To run this program outside CI, you must ensure that Google Cloud SDK Application credentials
are setup correctly, to do so, run:
```
gcloud auth application-default login
```

Next, build it:
```
go build ./
```

Run:
```
./requester --namespace=test-clusters-dev --description="<your name and purpose of this cluster>"
```

## CI Usage

This program supports traditional `GOOGLE_APPLICATION_CREDENTIALS` environment variable, but also
for convenience it has `GCP_SERVICE_ACCOUNT_KEY` that is expected to contain a base64-encoded
JSON service account key (i.e. no need to have the data written to a file).

For GitHub Actions, it's recommended to use the official image:

```
      - name: Request GKE test cluster
        uses: docker://quay.io/isovalent/gke-test-cluster-requester:ad06d7c2151d012901fc2ddc92406044f2ffba2d
        env:
          GCP_SERVICE_ACCOUNT_KEY: ${{ secrets.GCP_SERVICE_ACCOUNT_KEY }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          args: --namespace=... --image=...
```
