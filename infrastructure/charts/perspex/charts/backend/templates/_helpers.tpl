{{/*
Expand the name of the chart.
*/}}
{{- define "backend.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "backend.fullname" -}}
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
{{- define "backend.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "backend.labels" -}}
helm.sh/chart: {{ include "backend.chart" . }}
{{ include "backend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "backend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "backend.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "backend.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "backend.name" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Define default environment variables
*/}}
{{- define "backend.environmentVariables" -}}
- name: BACKEND_HOST
  value: {{ .Values.global.config.backend.host | default .Values.config.backend.host | quote }}
- name: BACKEND_HTTP_PORT
  value: {{ .Values.service.port | quote }}
- name: BACKEND_GRPC_PORT
  value: {{ .Values.service.grpcPort | quote }}
- name: BACKEND_LOG_MODE
  value: {{ .Values.global.config.backend.logmode | default  .Values.config.backend.logmode | quote }}
- name: WRITER_POSTGRES_HOST
  value: {{ .Values.global.config.database.writer.host | default .Values.config.database.writer.host | quote }}
- name: WRITER_POSTGRES_PORT
  value: {{ .Values.global.config.database.writer.port | default .Values.config.database.writer.port | quote }}
- name: WRITER_POSTGRES_USER
  value: {{ .Values.global.config.database.writer.user | default .Values.config.database.writer.user | quote }}
- name: WRITER_POSTGRES_SCHEMA
  value: {{ .Values.global.config.database.writer.schema | default .Values.config.database.writer.schema | quote }}
- name: WRITER_POSTGRES_DB
  value: {{ .Values.global.config.database.writer.name | default .Values.config.database.writer.name | quote }}
- name: WRITER_POSTGRES_MAX_OPEN_CONNECTIONS
  value: {{ .Values.global.config.database.writer.maxOpenConnections | default  .Values.config.database.writer.maxOpenConnections | quote }}
- name: WRITER_POSTGRES_MAX_IDLE_CONNECTIONS
  value: {{ .Values.global.config.database.writer.maxIdleConnections | default .Values.config.database.writer.maxIdleConnections | quote }}
- name: WRITER_POSTGRES_CONNECTION_LIFESPAN
  value: {{ .Values.global.config.database.writer.connectionLifespan | default .Values.config.database.writer.connectionLifespan | quote }}
- name: WRITER_POSTGRES_DEBUG
  value: {{ .Values.global.config.database.writer.debug | default .Values.config.database.writer.debug | quote }}
- name: READER_POSTGRES_HOST
  value: {{ .Values.global.config.database.reader.host | default .Values.config.database.reader.host | quote }}
- name: READER_POSTGRES_PORT
  value: {{ .Values.global.config.database.reader.port | default .Values.config.database.reader.port | quote }}
- name: READER_POSTGRES_USER
  value: {{ .Values.global.config.database.reader.user | default .Values.config.database.reader.user | quote }}
- name: READER_POSTGRES_SCHEMA
  value: {{ .Values.global.config.database.reader.schema | default .Values.config.database.reader.schema | quote }}
- name: READER_POSTGRES_DB
  value: {{ .Values.global.config.database.reader.name | default .Values.config.database.reader.name | quote }}
- name: READER_POSTGRES_MAX_OPEN_CONNECTIONS
  value: {{ .Values.global.config.database.reader.maxOpenConnections | default  .Values.config.database.reader.maxOpenConnections | quote }}
- name: READER_POSTGRES_MAX_IDLE_CONNECTIONS
  value: {{ .Values.global.config.database.reader.maxIdleConnections | default .Values.config.database.reader.maxIdleConnections | quote }}
- name: READER_POSTGRES_CONNECTION_LIFESPAN
  value: {{ .Values.global.config.database.reader.connectionLifespan | default .Values.config.database.reader.connectionLifespan | quote }}
- name: READER_POSTGRES_DEBUG
  value: {{ .Values.global.config.database.reader.debug | default .Values.config.database.reader.debug | quote }}
{{- range $key, $value := .Values.extraEnv }}
- name: {{ $key | quote }}
  value: {{ $value | quote }}
{{- end }}
{{- end }}
