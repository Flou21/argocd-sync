# argocd-sync

little container to wait for argocd sync in ci pipelines

## usage

### parameter

#### --app

the app that should be synced and waited for

### env vars

#### ARGOCD_SERVER

something like argocd.domain.tld

#### ARGOCD_AUTH_TOKEN

see this for more information: https://argoproj.github.io/argo-cd/operator-manual/security/


### example

````
docker run -e ARGOCD_SERVER="argocd.domain.tld" -e ARGOCD_AUTH_TOKEN="<your-token>" <some-registry>/argocd-sync argocd-sync --app <your-app>
````