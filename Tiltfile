load("ext://namespace", "namespace_create", "namespace_inject")
load("ext://restart_process", "docker_build_with_restart")
load("ext://helm_remote", "helm_remote")
load("ext://local_output", "local_output")

image_repo="registry.local:5000"
k8s_context="perspex-local"
deploy_namespace="perspex"
values_file="infrastructure/tilt/values-dev.yaml"

allow_k8s_contexts(k8s_context)
namespace_create(deploy_namespace)

services={
  "backend": "false",
  "migration": "true"
}

#########################
# Charts
#########################

perspex = helm(
  "infrastructure/charts/perspex/",
  name=deploy_namespace,
  namespace=deploy_namespace,
  values=values_file,
  set=[
    "backend.enabled={s}".format(s=services["backend"]),
    "migration.enabled={s}".format(s=services["migration"]),
  ]
)

k8s_yaml(perspex)
k8s_resource(workload="postgresql", labels=["postgres"])
#k8s_resource(workload="perspex-migration", labels=["postgres"])

#########################
# Services
#########################

# this doesn't  have a conditional block because we have settings inside this service
include("services/migration/Tiltfile")

if services["backend"] == "true": 
    include('services/backend/Tiltfile')


#########################
# Utils
#########################

local_resource(
  "postgresql-port-forward",
  serve_cmd="kubectl -n {deploy_namespace} port-forward service/postgresql 5433:5432".format(deploy_namespace=deploy_namespace),
  trigger_mode=TRIGGER_MODE_MANUAL,
  auto_init=False,
  labels=["postgres"]
)
