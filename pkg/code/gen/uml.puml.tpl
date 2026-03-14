@startuml
top to bottom direction
skinparam package {
    BackgroundColor #fafafa
}
skinparam defaulttextalignment center
skinparam noteTextAlignment left

{{- $nameColour := (mapColour "name") }}
{{- $embeddedColour := (mapColour "embedded") }}
{{- $fieldTypeColour := (mapColour "fieldType") }}
{{- $methodSignatureColour := (mapColour "methodSignature") }}
{{- $funcSignatureColour := (mapColour "funcSignature") }}

{{- if (not skipLegend) }}
package legend {
    class "type" as legend_Type << ({{ getTypeInitial "type" }}, {{ mapTypeColour "type" }}) >> {}
    class "primitive" as legend_Primitive << ({{ getTypeInitial "int" }}, {{ mapTypeColour "int" }}) >> {}
    class "composite" as legend_Composite << ({{ getTypeInitial "map[]" }}, {{ mapTypeColour "map[]" }}) >> {}
    class "pointer" as legend_Pointer << ({{ getTypeInitial "*" }}, {{ mapTypeColour "*" }}) >> {}
    class "struct" as legend_Struct << ({{ getTypeInitial "struct" }}, {{ mapTypeColour "struct" }}) >> {}
    class "interface" as legend_Interface << ({{ getTypeInitial "interface" }}, {{ mapTypeColour "interface" }}) >> {}
    class "var" as legend_Var << ({{ getTypeInitial "var" }}, {{ mapTypeColour "var" }}) >> {}
    class "const" as legend_Const << ({{ getTypeInitial "const" }}, {{ mapTypeColour "const" }}) >> {}
    class "func" as legend_Func << ({{ getTypeInitial "func" }}, {{ mapTypeColour "func" }}) >> {}
    class "external" as legend_External << ({{ getTypeInitial "external" }}, {{ mapTypeColour "external" }}) >> {}
    class "TypeA" as legend_TypeA << ({{ getTypeInitial "type" }}, {{ mapTypeColour "type" }}) >> {
        + <b>ExportedField</b> string
        - <font color={{ $nameColour }}>notExportedField</font> string
        + <b>ExportedMethod</b> <font color={{ $funcSignatureColour }}>(input string)</font>
        - <font color={{ $nameColour }}>notExportedMethod</font> <font color={{ $funcSignatureColour }}>(input string)</font>
    }
    class "TypeB" as legend_TypeB << ({{ getTypeInitial "type" }}, {{ mapTypeColour "type" }}) >> {
        + <b>TypeA</b> <font color={{ $fieldTypeColour }}>TypeA</font>
    }
    class "TypeC" as legend_TypeC << ({{ getTypeInitial "type" }}, {{ mapTypeColour "type" }}) >> {
        <font color={{ $embeddedColour }}>TypeA</font>
    }
    class "TypeD" as legend_TypeD << ({{ getTypeInitial "string" }}, {{ mapTypeColour "string" }}) >> {}
    class "string" as legend_string << ({{ getTypeInitial "string" }}, {{ mapTypeColour "string" }}) >> {}
    class "TypeE" as legend_TypeE << ({{ getTypeInitial "type" }}, {{ mapTypeColour "type" }}) >> {
        + <b>GetTypeD</b> <font color={{ $funcSignatureColour }}>() TypeD</font>
    }
    class "Interface" as legend_Interface << ({{ getTypeInitial "interface" }}, {{ mapTypeColour "interface" }}) >> {
        + <b>ExportedMethod</b> <font color={{ $funcSignatureColour }}>(input string)</font>
    }
    legend_TypeB o--> legend_TypeA: Aggregation
    legend_TypeC *--> legend_TypeA: Composition
    legend_TypeA ..|> legend_Interface: Inheritance / Generalisation
    legend_TypeC ..|> legend_Interface: Inheritance / Generalisation
    legend_TypeD #..> legend_string: Type Alias or Underlying Type Link
    legend_TypeE ..> legend_TypeD: Dependency
}
{{- end }}
{{- range $pkgPath, $codeMap := . }}
{{- $pkg := (generateElementID $pkgPath) }}
package "{{ $pkgPath }}" as {{ $pkg }} {
{{- if gt (len $codeMap.Type) 0 }}
    {{- range $type := $codeMap.Type }}
        {{- if not (and (not $type.IsExported) skipNotExported) }}
    class "{{ if $type.IsExported }}<b>{{ end }}<font color={{ $nameColour }}>{{ $type.Name }}</font>{{ if $type.IsExported }}</b>{{ end }}{{ if $type.Params }}<{{ $type.Params.Label }}>{{ end }}" as {{ $pkg }}_{{ generateElementID $type.Name }} << ({{ getTypeInitial $type.Type.Label }}, {{ mapTypeColour $type.Type.Label }}) >> {
            {{- range $embedded := $type.Embedded }}
        <font color={{ $embeddedColour }}>{{ $embedded.Label }}</font>
            {{- end }}
            {{- range $field := $type.Fields }}
                {{- if $field.IsExported }}
        + {{ $field.Name }} <font color={{ $fieldTypeColour }}>{{ if $field.Type }}{{ $field.Type.Label }}{{ else }}{{ $field.TypeInfo }}{{ end }}</font>
                {{- end }}
            {{- end }}
            {{- range $field := $type.Fields }}
                {{- if (and (not $field.IsExported) (not skipNotExported)) }}
        - {{ $field.Name }} <font color={{ $fieldTypeColour }}>{{ if $field.Type }}{{ $field.Type.Label }}{{ else }}{{ $field.TypeInfo }}{{ end }}</font>
                {{- end }}
            {{- end }}
            {{- range $method := $type.Methods }}
                {{- if $method.IsExported }}
        + {{ $method.Name }} <font color={{ $methodSignatureColour }}>{{ $method.Signature.Label }}</font>
                {{- end }}
            {{- end }}
            {{- range $method := $type.Methods }}
                {{- if (and (not $method.IsExported) (not skipNotExported)) }}
        - {{ $method.Name }} <font color={{ $methodSignatureColour }}>{{ $method.Signature.Label }}</font>
                {{- end }}
            {{- end }}
    }
        {{- end }}
    {{- end }}
{{- end }}
{{- if not skipVar }}
    {{- if gt (len $codeMap.Var) 0 }}
        {{- range $var := $codeMap.Var }}
            {{- if not (and (not $var.IsExported) skipNotExported) }}
    class "{{ if $var.IsExported }}<b>{{ end }}<font color={{ $nameColour }}>{{ $var.Name }}</font>{{ if $var.IsExported }}</b>{{ end }}{{ if $var.Params }}<{{ $var.Params.Label }}>{{ end }}" as {{ $pkg }}_{{ generateElementID $var.Name }} << ({{ getTypeInitial "var" }}, {{ mapTypeColour "var" }}) >> {}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
{{- if not skipConst }}
    {{- if gt (len $codeMap.Const) 0 }}
        {{- range $const := $codeMap.Const }}
            {{- if not (and (not $const.IsExported) skipNotExported) }}
    class "{{ if $const.IsExported }}<b>{{ end }}<font color={{ $nameColour }}>{{ $const.Name }}</font>{{ if $const.IsExported }}</b>{{ end }}" as {{ $pkg }}_{{ generateElementID $const.Name }} << ({{ getTypeInitial "const" }}, {{ mapTypeColour "const" }}) >> {}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
{{- if not skipFunc }}
    {{- if gt (len $codeMap.Func) 0 }}
        {{- range $func := $codeMap.Func }}
            {{- if not (and (not $func.IsExported) skipNotExported) }}
    class "{{ if $func.IsExported }}<b>{{ end }}<font color={{ $nameColour }}>{{ $func.Name }}</font>{{ if $func.IsExported }}</b>{{ end }}<font color={{ $funcSignatureColour }}>{{ $func.Signature.Label }}</font>{{ if $func.Params }}<{{ $func.Params.Label }}>{{ end }}" as {{ $pkg }}_{{ generateElementID $func.Name }} << ({{ getTypeInitial "func" }}, {{ mapTypeColour "func" }}) >> {}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
}
{{- end }}

{{- range $pkgPath, $codeMap := . }}
{{- $pkg := (generateElementID $pkgPath) }}
package "{{ $pkgPath }}" as {{ $pkg }} {
{{- if gt (len $codeMap.Type) 0 }}
    {{- range $type := $codeMap.Type }}
        {{- if not (and (not $type.IsExported) skipNotExported) }}
            {{- range $link := $type.Type.Links }}
                {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                    {{- if not (typeExists (concat $link.PackagePath "." $link.Name)) }}
                        {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    package "{{ $link.PackagePath }}" as .{{ generateElementID $link.PackagePath }} {
                        {{- end }}
{{ if and $link.PackageName (ne $pkgPath $link.PackagePath) }}    {{ end }}    class "{{ $link.Name }}" as {{ generateElementID $link.PackagePath }}_{{ generateElementID $link.Name }} << ({{ getTypeInitial $link.Underlying }}, {{ mapTypeColour $link.Underlying }}) >> {}
                        {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    }
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- if not (or (isAny $type.Type.Label) (isInterface $type.Type.Label) (isStruct $type.Type.Label) (isSelector $type.Type.Label) (isPrimitive $type.Type.Label)) }}
    class "{{ $type.Type.Label }}" as {{ $pkg }}_{{ generateElementID $type.Type.Label }} << ({{ getTypeInitial $type.Type.Label }}, {{ mapTypeColour $type.Type.Label }}) >> {}
            {{- end }}
            {{- range $embedded := $type.Embedded }}
                {{- range $link := $embedded.Links }}
                    {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                        {{- if not (typeExists (concat $link.PackagePath "." $link.Name)) }}
                            {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    package "{{ $link.PackagePath }}" as .{{ generateElementID $link.PackagePath }} {
                            {{- end }}
{{ if and $link.PackageName (ne $pkgPath $link.PackagePath) }}    {{ end }}    class "{{ $link.Name }}" as {{ generateElementID $link.PackagePath }}_{{ generateElementID $link.Name }} << ({{ getTypeInitial $link.Underlying }}, {{ mapTypeColour $link.Underlying }}) >> {}
                            {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    }
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- if $type.Params }}
                {{- range $item := $type.Params.List }}
                    {{- range $link := $item.Constrain.Links }}
                        {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                            {{- if not (typeExists (concat $link.PackagePath "." $link.Name)) }}
                                {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    package "{{ $link.PackagePath }}" as .{{ generateElementID $link.PackagePath }} {
                                {{- end }}
{{ if and $link.PackageName (ne $pkgPath $link.PackagePath) }}    {{ end }}    class "{{ $link.Name }}" as {{ generateElementID $link.PackagePath }}_{{ generateElementID $link.Name }} << ({{ getTypeInitial $link.Underlying }}, {{ mapTypeColour $link.Underlying }}) >> {}
                                {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    }
                                {{- end }}
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- range $interface := $type.Interfaces }}
                {{- if not (typeExists (concat $interface.PackagePath "." $interface.Name)) }}
                    {{- if and $interface.PackageName (ne $pkgPath $interface.PackagePath) }}
    package "{{ $interface.PackagePath }}" as .{{ generateElementID $interface.PackagePath }} {
                    {{- end }}
{{ if and $interface.PackageName (ne $pkgPath $interface.PackagePath) }}    {{ end }}    class "{{ $interface.Name }}" as {{ generateElementID $interface.PackagePath }}_{{ generateElementID $interface.Name }} << ({{ getTypeInitial "interface" }}, {{ mapTypeColour "interface" }}) >> {}
                    {{- if and $interface.PackageName (ne $pkgPath $interface.PackagePath) }}
    }
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- range $field := $type.Fields }}
                {{- range $link := $field.Type.Links }}
                    {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                        {{- if not (typeExists (concat $link.PackagePath "." $link.Name)) }}
                            {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    package "{{ $link.PackagePath }}" as .{{ generateElementID $link.PackagePath }} {
                            {{- end }}
{{ if and $link.PackageName (ne $pkgPath $link.PackagePath) }}    {{ end }}    class "{{ $link.Name }}" as {{ generateElementID $link.PackagePath }}_{{ generateElementID $link.Name }} << ({{ getTypeInitial $link.Underlying }}, {{ mapTypeColour $link.Underlying }}) >> {}
                            {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    }
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- range $method := $type.Methods }}
                {{- range $link := $method.Signature.Links }}
                    {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                        {{- if not (typeExists (concat $link.PackagePath "." $link.Name)) }}
                            {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    package "{{ $link.PackagePath }}" as .{{ generateElementID $link.PackagePath }} {
                            {{- end }}
{{ if and $link.PackageName (ne $pkgPath $link.PackagePath) }}    {{ end }}    class "{{ $link.Name }}" as {{ generateElementID $link.PackagePath }}_{{ generateElementID $link.Name }} << ({{ getTypeInitial $link.Underlying }}, {{ mapTypeColour $link.Underlying }}) >> {}
                            {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    }
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
{{- if not skipVar }}
    {{- if gt (len $codeMap.Var) 0 }}
        {{- range $var := $codeMap.Var }}
            {{- if not (and (not $var.IsExported) skipNotExported) }}
                {{- if $var.Type }}
                    {{- range $link := $var.Type.Links }}
                        {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                            {{- if not (typeExists (concat $link.PackagePath "." $link.Name)) }}
                                {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    package "{{ $link.PackagePath }}" as .{{ generateElementID $link.PackagePath }} {
                                {{- end }}
{{ if and $link.PackageName (ne $pkgPath $link.PackagePath) }}    {{ end }}    class "{{ $link.Name }}" as {{ generateElementID $link.PackagePath }}_{{ generateElementID $link.Name }} << ({{ getTypeInitial $link.Underlying }}, {{ mapTypeColour $link.Underlying }}) >> {}
                                {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    }
                                {{- end }}
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- else }}
                    {{- if not (and (isPrimitive $var.TypeInfo) skipPrimitive) }}
                        {{- if not (typeExists $var.TypeInfo) }}
                        {{- $typeInfoParts := (getFullPathParts $var.TypeInfo)}}
                            {{- if and $typeInfoParts $typeInfoParts.PackagePath (ne $pkgPath $typeInfoParts.PackagePath) }}
    package "{{ $typeInfoParts.PackagePath }}" as .{{ generateElementID $typeInfoParts.PackagePath }} {
                            {{- end }}
{{ if and $typeInfoParts $typeInfoParts.PackagePath (ne $pkgPath $typeInfoParts.PackagePath) }}    {{ end }}    class "{{ if $typeInfoParts }}{{ $typeInfoParts.Name }}{{ else }}{{ $var.TypeInfo }}{{ end }}" as {{ if $typeInfoParts }}{{ generateElementID $typeInfoParts.PackagePath }}_{{ generateElementID $typeInfoParts.Name }}{{ else }}{{ generateElementID $var.TypeInfo }}{{ end }} << ({{ getTypeInitial $var.TypeInfo }}, {{ mapTypeColour $var.TypeInfo }}) >> {}
                            {{- if and $typeInfoParts $typeInfoParts.PackagePath (ne $pkgPath $typeInfoParts.PackagePath) }}
    }
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
{{- if not skipConst }}
    {{- if gt (len $codeMap.Const) 0 }}
        {{- range $const := $codeMap.Const }}
            {{- if not (and (not $const.IsExported) skipNotExported) }}
                {{- if $const.Type }}
                    {{- range $link := $const.Type.Links }}
                        {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                            {{- if not (typeExists (concat $link.PackagePath "." $link.Name)) }}
                                {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    package "{{ $link.PackagePath }}" as .{{ generateElementID $link.PackagePath }} {
                                {{- end }}
{{ if and $link.PackageName (ne $pkgPath $link.PackagePath) }}    {{ end }}    class "{{ $link.Name }}" as {{ generateElementID $link.PackagePath }}_{{ generateElementID $link.Name }} << ({{ getTypeInitial $link.Underlying }}, {{ mapTypeColour $link.Underlying }}) >> {}
                                {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    }
                                {{- end }}
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- else }}
                    {{- if not (and (isPrimitive $const.TypeInfo) skipPrimitive) }}
                        {{- if not (typeExists $const.TypeInfo) }}
                        {{- $typeInfoParts := (getFullPathParts $const.TypeInfo)}}
                            {{- if and $typeInfoParts $typeInfoParts.PackagePath (ne $pkgPath $typeInfoParts.PackagePath) }}
    package "{{ $typeInfoParts.PackagePath }}" as .{{ generateElementID $typeInfoParts.PackagePath }} {
                            {{- end }}
{{ if and $typeInfoParts $typeInfoParts.PackagePath (ne $pkgPath $typeInfoParts.PackagePath) }}    {{ end }}    class "{{ if $typeInfoParts }}{{ $typeInfoParts.Name }}{{ else }}{{ $const.TypeInfo }}{{ end }}" as {{ if $typeInfoParts }}{{ generateElementID $typeInfoParts.PackagePath }}_{{ generateElementID $typeInfoParts.Name }}{{ else }}{{ generateElementID $const.TypeInfo }}{{ end }} << ({{ getTypeInitial $const.TypeInfo }}, {{ mapTypeColour $const.TypeInfo }}) >> {}
                            {{- if and $typeInfoParts $typeInfoParts.PackagePath (ne $pkgPath $typeInfoParts.PackagePath) }}
    }
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
{{- if not skipFunc }}
    {{- if gt (len $codeMap.Func) 0 }}
        {{- range $func := $codeMap.Func }}
            {{- if not (and (not $func.IsExported) skipNotExported) }}
                {{- range $link := $func.Signature.Links }}
                    {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                        {{- if not (typeExists (concat $link.PackagePath "." $link.Name)) }}
                            {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    package "{{ $link.PackagePath }}" as .{{ generateElementID $link.PackagePath }} {
                            {{- end }}
{{ if and $link.PackageName (ne $pkgPath $link.PackagePath) }}    {{ end }}    class "{{ $link.Name }}" as {{ generateElementID $link.PackagePath }}_{{ generateElementID $link.Name }} << ({{ getTypeInitial $link.Underlying }}, {{ mapTypeColour $link.Underlying }}) >> {}
                            {{- if and $link.PackageName (ne $pkgPath $link.PackagePath) }}
    }
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
}
{{- end }}

{{- range $pkgPath, $codeMap := . }}
{{- $pkg := (generateElementID $pkgPath) }}
package "{{ $pkgPath }}" as {{ $pkg }} {
{{- if gt (len $codeMap.Type) 0 }}
    {{- range $type := $codeMap.Type }}
        {{- if not (and (not $type.IsExported) skipNotExported) }}
            {{- range $link := $type.Type.Links }}
                {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                    {{- $link := (concat $pkg "_" (generateElementID $type.Name) " #...up...> " (generateElementID $link.PackagePath) "_" (generateElementID $link.Name)) }}
                    {{- if not (linkInCache $link) }}
    {{ $link }}
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- if not (or (isAny $type.Type.Label) (isInterface $type.Type.Label) (isStruct $type.Type.Label) (isSelector $type.Type.Label) (isPrimitive $type.Type.Label)) }}
                {{- $link := (concat $pkg "_" (generateElementID $type.Name) " #...up...> " $pkg "_" (generateElementID $type.Type.Label))}}
                {{- if not (linkInCache $link) }}
    {{ $link }}
                {{- end }}
            {{- end }}
            {{- range $embedded := $type.Embedded }}
                {{- range $link := $embedded.Links }}
                    {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                        {{- $link := (concat $pkg "_" (generateElementID $type.Name) " *---up---> " (generateElementID $link.PackagePath) "_" (generateElementID $link.Name)) }}
                        {{- if not (linkInCache $link) }}
    {{ $link }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- if $type.Params }}
                {{- range $item := $type.Params.List }}
                    {{- range $link := $item.Constrain.Links }}
                        {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                            {{- $link := (concat $pkg "_" (generateElementID $type.Name) " #...up...> " (generateElementID $link.PackagePath) "_" (generateElementID $link.Name)) }}
                            {{- if not (linkInCache $link) }}
    {{ $link }}
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- range $interface := $type.Interfaces }}
                {{- if not (and (isPrimitive $interface.Name) skipPrimitive) }}
                    {{- $link := (concat $pkg "_" (generateElementID $type.Name) " ...up...|> " (generateElementID $interface.PackagePath) "_" (generateElementID $interface.Name)) }}
                    {{- if not (linkInCache $link) }}
    {{ $link }}
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- range $field := $type.Fields }}
                {{- range $link := $field.Type.Links }}
                    {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                        {{- $link := (concat $pkg "_" (generateElementID $type.Name) " o---up---> " (generateElementID $link.PackagePath) "_" (generateElementID $link.Name)) }}
                        {{- if not (linkInCache $link) }}
    {{ $link }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- range $method := $type.Methods }}
                {{- range $link := $method.Signature.Links }}
                    {{- if not (and (not (isPrimitive $link.Name)) (not (isError $link.Name)) (not (isInterface $link.Name)) (not (isAny $link.Name)) (not (isExported $link.Name)) skipNotExported) }}
                        {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                            {{- $link := (concat $pkg "_" (generateElementID $type.Name) " ...up...> " (generateElementID $link.PackagePath) "_" (generateElementID $link.Name)) }}
                            {{- if not (linkInCache $link) }}
    {{ $link }}
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
{{- if not skipVar }}
    {{- if gt (len $codeMap.Var) 0 }}
        {{- range $var := $codeMap.Var }}
            {{- if not (and (not $var.IsExported) skipNotExported) }}
                {{- if $var.Type }}
                    {{- range $link := $var.Type.Links }}
                        {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                            {{- $link := (concat $pkg "_" (generateElementID $var.Name) " #...down...> " (generateElementID $link.PackagePath) "_" (generateElementID $link.Name)) }}
                            {{- if not (linkInCache $link) }}
    {{ $link }}
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- else }}
                    {{- if not (and (isPrimitive $var.TypeInfo) skipPrimitive) }}
                        {{- $typeInfoParts := (getFullPathParts $var.TypeInfo)}}
                        {{- $target := (or (and $typeInfoParts (concat (generateElementID $typeInfoParts.PackagePath) "_" (generateElementID $typeInfoParts.Name))) (generateElementID $var.TypeInfo)) }}
                        {{- $link := (concat $pkg "_" (generateElementID $var.Name) (ternary (isError $var.TypeInfo) " ...up...|> " " #...up...> ") $target) }}
                        {{- if not (linkInCache $link) }}
    {{ $link }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
{{- if not skipConst }}
    {{- if gt (len $codeMap.Const) 0 }}
        {{- range $const := $codeMap.Const }}
            {{- if not (and (not $const.IsExported) skipNotExported) }}
                {{- if $const.Type }}
                    {{- range $link := $const.Type.Links }}
                        {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                            {{- $link := (concat $pkg "_" (generateElementID $const.Name) " #...down...> " (generateElementID $link.PackagePath) "_" (generateElementID $link.Name)) }}
                            {{- if not (linkInCache $link) }}
    {{ $link }}
                            {{- end }}
                        {{- end }}
                    {{- end }}
                {{- else }}
                    {{- if not (and (isPrimitive $const.TypeInfo) skipPrimitive) }}
                        {{- $typeInfoParts := (getFullPathParts $const.TypeInfo)}}
                        {{- $target := (or (and $typeInfoParts (concat (generateElementID $typeInfoParts.PackagePath) "_" (generateElementID $typeInfoParts.Name))) (generateElementID $const.TypeInfo)) }}
                        {{- $link := (concat $pkg "_" (generateElementID $const.Name) " #...up...> " $target) }}
                        {{- if not (linkInCache $link) }}
    {{ $link }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
{{- if not skipFunc }}
    {{- if gt (len $codeMap.Func) 0 }}
        {{- range $func := $codeMap.Func }}
            {{- if not (and (not $func.IsExported) skipNotExported) }}
                {{- range $link := $func.Signature.Links }}
                    {{- if not (and (isPrimitive $link.Name) skipPrimitive) }}
                        {{- $link := (concat $pkg "_" (generateElementID $func.Name) " ....up....> " (generateElementID $link.PackagePath) "_" (generateElementID $link.Name)) }}
                        {{- if not (linkInCache $link) }}
    {{ $link }}
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
}
{{- end }}

{{- if not skipDoc }}
    {{- range $pkgPath, $codeMap := . }}
    {{- $pkg := (generateElementID $pkgPath) }}
package "{{ $pkgPath }}" as {{ $pkg }} {
        {{- if gt (len $codeMap.Type) 0 }}
            {{- range $type := $codeMap.Type }}
                {{- if not (and (not $type.IsExported) skipNotExported) }}
                    {{- if $type.Doc }}
note top of {{ concat $pkg "_" (generateElementID $type.Name) }}
{{ $type.Doc }}
end note
                    {{- end }}
                    {{- range $field := $type.Fields }}
                        {{- if $field.Doc }}
note left of {{ concat $pkg "_" (generateElementID $type.Name) "::" $field.Name }}
{{ $field.Doc }}
end note
                        {{- end }}
                    {{- end }}
                    {{- range $method := $type.Methods }}
                        {{- if $method.Doc }}
note left of {{ concat $pkg "_" (generateElementID $type.Name) "::" $method.Name }}
{{ $method.Doc }}
end note
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
        {{- if not skipVar }}
            {{- if gt (len $codeMap.Var) 0 }}
                {{- range $var := $codeMap.Var }}
                    {{- if not (and (not $var.IsExported) skipNotExported) }}
                        {{- if $var.Doc }}
note left of {{ concat $pkg "_" (generateElementID $var.Name) }}
{{ $var.Doc }}
end note
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
        {{- if not skipConst }}
            {{- if gt (len $codeMap.Const) 0 }}
                {{- range $const := $codeMap.Const }}
                    {{- if not (and (not $const.IsExported) skipNotExported) }}
                         {{- if $const.Doc }}
 note left of {{ concat $pkg "_" (generateElementID $const.Name) }}
 {{ $const.Doc }}
 end note
                         {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
        {{- if not skipFunc }}
            {{- if gt (len $codeMap.Func) 0 }}
                {{- range $func := $codeMap.Func }}
                    {{- if not (and (not $func.IsExported) skipNotExported) }}
                        {{- if $func.Doc }}
note left of {{ concat $pkg "_" (generateElementID $func.Name) }}
{{ $func.Doc }}
end note
                        {{- end }}
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
}
    {{- end }}
{{- end }}
@enduml
