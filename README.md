### Quickly use rr+delve to debug openshift installer

0. modify `hack/config.sh` for installer branch
1. run `./hack/build.sh` ... five hours later ...
2. create `installer-dir/install-config.yaml`
3. run `./hack/run.sh` ... wait till explosion
4. run `./hack/debug.sh` ... delve running over port 2345
5. profit
