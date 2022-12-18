load("ext://namespace", "namespace_create", "namespace_inject")
load("ext://restart_process", "docker_build_with_restart")
load("ext://helm_remote", "helm_remote")
load("ext://local_output", "local_output")

image_repo = os.getenv("TILT_IMAGE_REPO", default="registry.local:5000")
k8s_context = os.getenv("TILT_K8S_CONTEXT", default="perspex-local")
deploy_namespace = os.getenv("TILT_DEPLOY_NAMESPACE", default="perspex")
remote_cluster = os.getenv("TILT_REMOTE_CLUSTER", default=False)
values_file="infrastructure/tilt/values-dev.yaml"

allow_k8s_contexts(k8s_context)
namespace_create(deploy_namespace)

if remote_cluster == True :
  docker_login_cmd="aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 450613976047.dkr.ecr.us-east-1.amazonaws.com"
  local(docker_login_cmd)

services={
  "backend": os.getenv("TILT_BACKEND_ENABLED", default="true"),
  "contour": os.getenv("TILT_CONTOUR_ENABLED", default="false"),
  "gateway": os.getenv("TILT_GATEWAY_ENABLED", default="false"),
  "migration": os.getenv("TILT_MIGRATION_ENABLED", default="true"),
  "postgresql": os.getenv("TILT_POSTGRES_ENABLED", default="true"),
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
    "contour.enabled={s}".format(s=services["contour"]),
    "gateway.enabled={s}".format(s=services["gateway"]),
    "migration.enabled={s}".format(s=services["migration"]),
    "postgresql.enabled={s}".format(s=services["postgresql"]),
  ]
)

k8s_yaml(perspex)

#########################
# Services
#########################

if services["backend"] == "true": 
  include('services/backend/Tiltfile')

if services["contour"] == "true":
  k8s_resource(workload="perspex-contour-contour", labels=["contour"])
  k8s_resource(workload="perspex-contour-contour-certgen", labels=["contour"])
  k8s_resource(workload="perspex-contour-envoy", labels=["contour"])

if services["gateway"] == "true": 
  include('services/gateway/Tiltfile')

if services["migration"] == "true": 
  include("services/migration/Tiltfile")

if services["postgresql"] == "true":
  k8s_resource(workload="postgresql", labels=["postgresql"], port_forwards=[port_forward(5433, 5432, name="postgres")])
