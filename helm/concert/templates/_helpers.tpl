{{- define "concert.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- define "concert.fullname" -}}
{{- if .Values.fullnameOverride }}{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}{{- else }}{{- printf "%s-%s" .Release.Name (include "concert.name" .) | trunc 63 | trimSuffix "-" }}{{- end }}
{{- end }}
{{- define "concert.selectorLabels" -}}
app.kubernetes.io/name: {{ include "concert.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
{{- define "concert.labels" -}}
{{ include "concert.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}
{{- define "concert.dbHost" -}}
{{- if .Values.database.external.enabled }}{{- .Values.database.external.host }}{{- else }}{{- printf "%s-postgresql" (include "concert.fullname" .) }}{{- end }}
{{- end }}
{{- define "concert.dbPassword" -}}
{{- if .Values.database.external.enabled }}{{- .Values.database.external.password }}{{- else if .Values.secrets.dbPassword }}{{- .Values.secrets.dbPassword }}{{- else }}{{- .Values.database.internal.postgresql.auth.postgresPassword }}{{- end }}
{{- end }}
{{- define "concert.temporalHost" -}}
{{- if .Values.temporal.external.enabled }}{{- .Values.temporal.external.host }}{{- else }}{{- printf "%s-temporal" (include "concert.fullname" .) }}{{- end }}
{{- end }}
