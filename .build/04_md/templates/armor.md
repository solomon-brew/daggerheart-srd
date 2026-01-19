# {{ upper .name }}

> **Base Thresholds:** {{ .base_thresholds }} | **Base Score:** {{ .base_score }}

{{- if gt (len .feature) 0 }}
{{- $feat := index .feature 0 }}

**Feature:** **_{{ $feat.name }}:_** {{ $feat.text }}
{{- else }}

**Feature:** —
{{- end }}

_Armor - Tier {{ .tier }}_
