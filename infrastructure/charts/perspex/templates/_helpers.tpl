{{/*
Expand the name of the chart.
*/}}
{{- define "perspex.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "perspex.fullname" -}}
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
{{- define "perspex.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "perspex.labels" -}}
helm.sh/chart: {{ include "perspex.chart" . }}
{{ include "perspex.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "perspex.selectorLabels" -}}
app.kubernetes.io/name: {{ include "perspex.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "perspex.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "perspex.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "perspex.uniqueName" -}}
{{- $root := index . "root" -}}
{{- $fullname := include "perspex.fullname" $root -}}
{{- $job := index . "job" -}}
{{- $unique := index . "unique" -}}
{{- $now := ( now | unixEpoch ) -}}
{{- if $unique -}}
{{- printf "%s-%s-%s" $fullname $job $now }}
{{- else }}
{{- printf "%s-%s" $fullname $job }}
{{- end }}
{{- end -}}


{{- define "perspex.environmentVariables" -}}
- name: POSTGRES_HOST
  value: {{ .Values.global.config.database.writer.host | quote }}
- name: POSTGRES_PORT
  value: {{ .Values.global.config.database.writer.port | quote }}
- name: POSTGRES_USER
  value: {{ .Values.global.config.database.writer.user | quote }}
- name: POSTGRES_SCHEMA
  value: {{ .Values.global.config.database.writer.schema | quote }}
- name: POSTGRES_DB
  value: {{ .Values.global.config.database.writer.name | quote }}
{{- end }}
