# deploi

Deploi is a tool that manages the deployment of build artifacts into different cluster environments. For the start it focuses on Jenkins and Kubernetes.

## Idea

Jenkins has an arbitrary number of build pipelines. Some of those build pipelines produce artifacts that can be deployed to different cluster environments, e.g. a staging cluster or the production environment. Managing the different deploy targets (e.g. different Kubernetes clusters and/or namespaces) for a number of projects and branches inside Jenkins can be cumbersome. Therefore deploi acts as a small bridge between Jenkins as the supplier of new builds and different cluster environments as recipients of new artifacts. deploi allows to automate recurring deployment patterns. For instance a merge to master automatically triggers Jenkins to create a new artifact. Every artifact that comes out of this pipeline is automatically deployed to a staging environment.

## Concept

deploi consists of multiple components:

* The central deploi daemon (deploid)
* A number of deploi agents, one per target environment
* A deploi CLI to interact with the deploid API

The last step in every build pipeline that produces a deployable artifact is an API call to the deploid that announces metadata of the new build. That includes:

* The name of the project
* A unique version number that identifies the build (can be for example the git commit hash or the Jenkins build number)
* Some sort of URL where the build can be found (e.g. a docker registry)
* A link back to the Jenkins build
* The branch name

deploid stores the information about the new artifact in its database. Based on preconfigured rules or ad-hoc commands through the CLI, deploid plans a deployment of an artifact in one of the environments. Because of firewalls and/or complex ingress setups, the agents are not reachable from outside the cluster. Instead each agent pulls in a regular interval if there is anything new it should do. As soon as a job is done, the agent informs deploid.
