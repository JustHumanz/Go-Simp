### k8s deploy

first create PersistentVolume
```
kubectl apply -f vol/
```
after that apply all infra mainfest 
```
kubectl apply -f infra/
```
wait util all service/dev up and deploy migrate pod
```
kubectl apply -f apps/migrate.yaml
```
if migrate process already done deploy all mainfest in apps
```
kubectl apply -f apps
```    

#### Tips
Create configmaps for go-simp-web service