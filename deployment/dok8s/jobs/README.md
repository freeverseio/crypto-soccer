## Running jobs
Each time that you want to run a job in the kubernetes cluster you should build, tag and push the latest version to `freeverseio/jobs` or at least be sure that the image in that repository is the one that you want to use.

To do that you can use the `dockerpush.sh` file in this same folder. Add the `DOCKER_ID` and `DOCKER_PSSWD` needed to push.

Now if you want to run the job in your k8s cluster, you should change the variable `JOB_NAME` to the value of the job that you want to execute and add the environment variables that your job is going to need, then you can execute `run.sh` which will apply the `job.yml` in your current context.

`BE SURE TO CHECK THE CLUSTER CONTEXT BEFORE EXECUTING`


Basically the `freeverseio/jobs` image encapsulates all the binaries of the jobs, and lets you choose with an environment variable `JOB_NAME` which job you want to run.