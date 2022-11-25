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
  "backend": "true",
  "migration": "true",
  "postgres": "true"
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
    "postgres.enabled={s}".format(s=services["postgres"])
  ]
)

k8s_yaml(perspex)
k8s_resource(workload="postgresql", labels=["postgres"])

#########################
# Services
#########################

if services["migration"] == "true": 
  include("services/migration/Tiltfile")

if services["backend"] == "true": 
    include('services/backend/Tiltfile')

if services["postgres"] == "true": 
  local_resource(
    "postgresql-port-forward",
    serve_cmd="kubectl -n {deploy_namespace} port-forward service/postgresql 5433:5432".format(deploy_namespace=deploy_namespace),
    trigger_mode=TRIGGER_MODE_MANUAL,
    auto_init=False,
    labels=["postgres"]
  )
