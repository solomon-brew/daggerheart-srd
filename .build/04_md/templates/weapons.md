# {{ upper .name }}

**Trait:** {{ .trait }}; **Range:** {{ .range }}; **Damage:** {{ .damage }}; **Burden:** {{ .burden }}

{{- if gt (len .feature) 0 }}
{{- $feat := index .feature 0 }}

**Feature:** **_{{ $feat.name }}:_** {{ $feat.text }}
{{- else }}

**Feature:** —
{{- end }}

_{{ .primary_or_secondary }}{{ if .physical_or_magical }} {{ .physical_or_magical }}{{ end }} Weapon - Tier {{ .tier }}_
