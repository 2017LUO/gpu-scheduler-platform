{{- define "gpu-scheduler-platform.name" -}}
gpu-scheduler-platform
{{- end }}

{{- define "gpu-scheduler-platform.namespace" -}}
{{- if .Values.namespaceOverride -}}
{{ .Values.namespaceOverride }}
{{- else -}}
{{ .Release.Namespace }}
{{- end -}}
{{- end }}

{{- define "gpu-scheduler-platform.fullname" -}}
{{ .Release.Name }}
{{- end }}

{{- define "gpu-scheduler-platform.serviceAccountName" -}}
{{- if .Values.serviceAccount.name -}}
{{ .Values.serviceAccount.name }}
{{- else -}}
{{ include "gpu-scheduler-platform.fullname" . }}
{{- end -}}
{{- end }}