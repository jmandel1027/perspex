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
namespace_create("localstack")

if remote_cluster == True :
  docker_login_cmd="aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 450613976047.dkr.ecr.us-east-1.amazonaws.com"
  local(docker_login_cmd)

services={
  "backend":   os.getenv("TILT_BACKEND_ENABLED", default="true"),
  "frontend":   os.getenv("TILT_FRONTEND_ENABLED", default="true"),
  "localstack": os.getenv("TILT_LOCALSTACK_ENABLED", default="true"),
  "migration": os.getenv("TILT_MIGRATION_ENABLED", default="true"),
  "postgres":  os.getenv("TILT_POSTGRES_ENABLED", default="true"),
  "traefik": os.getenv("TILT_TRAEFIK_ENABLED", default="true"),
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
    "frontend.enabled={s}".format(s=services["frontend"]),
    "migration.enabled={s}".format(s=services["migration"]),
    "postgres.enabled={s}".format(s=services["postgres"]),
  ]
)

k8s_yaml(perspex)

#########################
# Services
#########################

if services["backend"] == "true": 
  include('services/backend/Tiltfile')

if services["frontend"] == "true": 
  include('services/frontend/Tiltfile')

if services["migration"] == "true": 
  include("services/migration/Tiltfile")

if services["postgres"] == "true":
  k8s_resource(workload="postgresql", labels=["postgres"], port_forwards=[port_forward(5433, 5432, name="postgres")])

if services["traefik"] == "true":
  helm_remote(
    "traefik",
    namespace=deploy_namespace,
    repo_name="traefik",
    repo_url="https://traefik.github.io/charts",
    release_name="traefik",
    version="20.8.0",
    values="infrastructure/tilt/traefik-values.yaml",
  )
  k8s_resource(workload="traefik", labels=["traefik"], port_forwards=[port_forward(9000, 9000, name="traefik")])

if services["localstack"] == 'true':
  helm_remote(
    'localstack',
    namespace='localstack',
    repo_name='localstack',
    repo_url='https://helm.localstack.cloud',
    release_name='localstack',
    version='0.5.2',
    values="infrastructure/tilt/localstack-values.yaml",
    set=[
        'extraEnvVars[0].name=DEFAULT_REGION',
        'extraEnvVars[0].value=us-east-1',
        'enableStartupScripts=true'
    ]
  )
  k8s_resource(workload='localstack', labels=["localstack"])
