# Scenario 1

## PVC: universedb-storage-universedb-0

```
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: universedb-storage-universedb-0
  namespace: freeverse
  selfLink: >-
    /api/v1/namespaces/freeverse/persistentvolumeclaims/universedb-storage-universedb-0
  uid: 59db8741-dd43-4d2c-95e4-45e3a2d7ba21
  resourceVersion: '8930284'
  creationTimestamp: '2020-08-04T10:02:48Z'
  labels:
    app.kubernetes.io/app: universedb
    app.kubernetes.io/component: universedb
    app.kubernetes.io/part-of: cryptosoccer
    app.kubernetes.io/version: 1.0.0
  annotations:
    pv.kubernetes.io/bind-completed: 'yes'
    pv.kubernetes.io/bound-by-controller: 'yes'
    volume.beta.kubernetes.io/storage-provisioner: dobs.csi.digitalocean.com
  finalizers:
    - kubernetes.io/pvc-protection
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 200Gi
  volumeName: pvc-59db8741-dd43-4d2c-95e4-45e3a2d7ba21
  storageClassName: do-block-storage
  volumeMode: Filesystem
status:
  phase: Bound
  accessModes:
    - ReadWriteOnce
  capacity:
    storage: 200Gi

```

## PV: pvc-59db8741-dd43-4d2c-95e4-45e3a2d7ba21

```
kind: PersistentVolume
apiVersion: v1
metadata:
  name: pvc-59db8741-dd43-4d2c-95e4-45e3a2d7ba21
  selfLink: /api/v1/persistentvolumes/pvc-59db8741-dd43-4d2c-95e4-45e3a2d7ba21
  uid: 59480b0e-53b7-4da6-95bc-74262daade92
  resourceVersion: '8930390'
  creationTimestamp: '2020-08-04T10:02:52Z'
  annotations:
    pv.kubernetes.io/provisioned-by: dobs.csi.digitalocean.com
  finalizers:
    - kubernetes.io/pv-protection
    - external-attacher/dobs-csi-digitalocean-com
spec:
  capacity:
    storage: 200Gi
  csi:
    driver: dobs.csi.digitalocean.com
    volumeHandle: a915e943-d639-11ea-9a23-0a58ac148067
    fsType: ext4
    volumeAttributes:
      storage.kubernetes.io/csiProvisionerIdentity: 1594669292138-8081-dobs.csi.digitalocean.com
  accessModes:
    - ReadWriteOnce
  claimRef:
    kind: PersistentVolumeClaim
    namespace: freeverse
    name: universedb-storage-universedb-0
    uid: 59db8741-dd43-4d2c-95e4-45e3a2d7ba21
    apiVersion: v1
    resourceVersion: '8930115'
  persistentVolumeReclaimPolicy: Delete
  storageClassName: do-block-storage
  volumeMode: Filesystem
status:
  phase: Bound

```

Note that the PV has `persistentVolumeReclaimPolicy: Delete` which on deleting the PVC will trigger the deletion of the volume and the deletion of the bound PV(although there is a bug in some k8s versions that would avoid this last thing to happen).

### Tests

Cluster with K8S 1.16.10-do.0

1. Scaling to 0 statefulset universedb, the PVC and PV will stay as before.

2. Deleting the PVC will delete also the PV, and since the PV had delete as reclaim policy it will also delete the volume from digital ocean block storage.

# Scenario 2

PVC

```
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: universedb-storage-universedb-0
  namespace: freeverse
  selfLink: >-
    /api/v1/namespaces/freeverse/persistentvolumeclaims/universedb-storage-universedb-0
  uid: 865fe6fe-0365-4348-b581-9bd5bc79d95f
  resourceVersion: '22432588'
  creationTimestamp: '2020-10-19T11:54:40Z'
  labels:
    app.kubernetes.io/app: universedb
    app.kubernetes.io/component: universedb
    app.kubernetes.io/part-of: cryptosoccer
    app.kubernetes.io/version: 1.0.0
  annotations:
    pv.kubernetes.io/bind-completed: 'yes'
    pv.kubernetes.io/bound-by-controller: 'yes'
    volume.beta.kubernetes.io/storage-provisioner: dobs.csi.digitalocean.com
  finalizers:
    - kubernetes.io/pvc-protection
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 200Gi
  volumeName: pvc-865fe6fe-0365-4348-b581-9bd5bc79d95f
  storageClassName: do-block-storage
  volumeMode: Filesystem
status:
  phase: Bound
  accessModes:
    - ReadWriteOnce
  capacity:
    storage: 200Gi

```

PV

```
kind: PersistentVolume
apiVersion: v1
metadata:
  name: pvc-865fe6fe-0365-4348-b581-9bd5bc79d95f
  selfLink: /api/v1/persistentvolumes/pvc-865fe6fe-0365-4348-b581-9bd5bc79d95f
  uid: 5e6d3d11-c566-497b-af97-2fb49b2b9613
  resourceVersion: '22432643'
  creationTimestamp: '2020-10-19T11:54:43Z'
  annotations:
    pv.kubernetes.io/provisioned-by: dobs.csi.digitalocean.com
  finalizers:
    - kubernetes.io/pv-protection
    - external-attacher/dobs-csi-digitalocean-com
spec:
  capacity:
    storage: 200Gi
  csi:
    driver: dobs.csi.digitalocean.com
    volumeHandle: dffb39d6-1201-11eb-9e6b-0a58ac148105
    fsType: ext4
    volumeAttributes:
      storage.kubernetes.io/csiProvisionerIdentity: 1599393852872-8081-dobs.csi.digitalocean.com
  accessModes:
    - ReadWriteOnce
  claimRef:
    kind: PersistentVolumeClaim
    namespace: freeverse
    name: universedb-storage-universedb-0
    uid: 865fe6fe-0365-4348-b581-9bd5bc79d95f
    apiVersion: v1
    resourceVersion: '22432507'
  persistentVolumeReclaimPolicy: Retain
  storageClassName: do-block-storage
  volumeMode: Filesystem
status:
  phase: Bound

```

1. Wont' scale universdb to 0 pods

2. Trying to delete pvc does nothing*.

3. Scaling statefulset to 0 deletes pvc and pv but does not delete dobs, it leaves it unattached to any droplet. Maybe it deletes the pvc and pv because it's empty? (Deleted it due to Terminating status in pvc)

4. Scaling statefulset universedb back up will create a new pvc and a new pv and a new volume in dobs. (It's not reusing the volume from before). It seems like it created a PV with reclaim policy delete.

5. Let's insert some data into the db in the newly created PV. Now what will happen if we try to delete the pvc while there is 1 pod running and then scale it down.

6. Deleting the pvc does nothing*.

7. Scaling the universedb statefulset deleted the pvc and the pv and also the dobs volume.

8. Let's scale back up to 1. New PVC, PV and dobs Volume created.

9. Without trying to delete the pvc let's scale down. Now the PVC and PV are still alive. The dobs Volume is alive but dettached! Maybe trying to delete the PVC sets it in a deleting state that when the pods is scaled to 0 triggers the deletion.

10. Scaling statefulset universedb will connect the pod with the pvc, pv and therefore with the dobs Volume from before.

11. Let's try to delete the pvc(which has the PV on delete reclaim policy) while the universedb pod is up and running and let's inspect its status through kubectl. Trying to delete did nothing visible... but executing `kubectl get pvc -n freeverse` lent this results:

```
NAME                                                 STATUS        VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS       AGE
universedb-storage-universedb-0                      Terminating   pvc-ce07c50b-511d-435f-864f-b692e7dca71c   200Gi      RWO            do-block-storage   7m12s

```

The pvc will keep in status Terminating, so when we try to scale down the statefulset which has a pod bound with this PVC the status Terminating will trigger the deletion from the PVC, PV and then due to the reclaim policy delete the dobs Volume.

* Deleting the pvc actually sets the status Terminating that will execute the deletion of the pvc when possible, in this case when the pods are scaled down to 0.



## Conclusions

If we want to swap a Persistent Volume, and therefore its dobs Volume, from a current universedb statefulset to another new statefulset we should:

1. Set PV reclaim policy to retain

2. Scale universedb statefulset to 0

3. Delete the PVC which will leave the PV in a Released Status and the dobs Volume dettached from any droplet.

4. Update the PV and delete the claimRef section of the yaml, after doing that the status will change from Released to Available. (If you delete an available pv the dobs Volume will be deleted or not based on reclaim policy)

5. Apply the new statefulset that we wish to couple with the released PV, with a PVC that matches the available PV specs.


Scaling down to 0 a statefulset and then scaling it back to 1 will allow the PVC, PV and dobs Volume to recouple succesfully and keep the same data.

