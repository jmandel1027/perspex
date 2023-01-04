{{/*
Expand the name of the chart.
*/}}
{{- define "frontend.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "frontend.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "frontend.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "frontend.labels" -}}
helm.sh/chart: {{ include "frontend.chart" . }}
{{ include "frontend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "frontend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "frontend.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "frontend.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "frontend.uniqueName" -}}
{{- $root := index . "root" -}}
{{- $fullname := include "frontend.fullname" $root -}}
{{- $job := index . "job" -}}
{{- $unique := index . "unique" -}}
{{- $now := ( now | unixEpoch ) -}}
{{- if $unique -}}
{{- printf "%s-%s-%s" $fullname $job $now }}
{{- else }}
{{- printf "%s-%s" $fullname $job }}
{{- end }}
{{- end -}}

{{/*
Define default environment variables
*/}}
{{- define "frontend.environmentVariables" -}}
- name: frontend_HOST
  value: {{ .Values.global.config.frontend.host | default .Values.config.frontend.host | quote }}
- name: frontend_PORT
  value: {{ .Values.service.port | quote }}
- name: frontend_SCHEME
  value: {{ .Values.global.config.frontend.scheme | default  .Values.config.frontend.scheme | quote }}
- name: POSTGRES_HOST
  value: {{ .Values.global.config.database.host | default .Values.config.database.host | quote }}
- name: POSTGRES_PORT
  value: {{ .Values.global.config.database.port | default .Values.config.database.port | quote }}
- name: POSTGRES_USER
  value: {{ .Values.global.config.database.user | default .Values.config.database.user | quote }}
- name: POSTGRES_SCHEMA
  value: {{ .Values.global.config.database.schema | default .Values.config.database.schema | quote }}
- name: POSTGRES_DB
  value: {{ .Values.global.config.database.name | default .Values.config.database.name | quote }}
- name: POSTGRES_MAX_OPEN_CONNECTIONS
  value: {{ .Values.global.config.database.maxOpenConnections | default  .Values.config.database.maxOpenConnections | quote }}
- name: POSTGRES_MAX_IDLE_CONNECTIONS
  value: {{ .Values.global.config.database.maxIdleConnections | default .Values.config.database.maxIdleConnections | quote }}
- name: POSTGRES_CONNECTION_LIFESPAN
  value: {{ .Values.global.config.database.connectionLifespan | default .Values.config.database.connectionLifespan | quote }}
- name: POSTGRES_DEBUG
  value: {{ .Values.global.config.database.debug | default .Values.config.database.debug | quote }}
{{- range $key, $value := .Values.extraEnv }}
- name: {{ $key | quote }}
  value: {{ $value | quote }}
{{- end }}
{{- end }}
