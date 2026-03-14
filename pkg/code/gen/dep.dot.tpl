digraph DepGraph {
    fontname="Courier New"
    splines="ortho"
    nodesep="1.4"
    ranksep="0.8"
    graph [rankdir="TB", pad="0.2"]
    node [shape="tab", fontname="Courier New", style="filled", fixedsize=false]
    edge [fontname="Courier New", color="grey66", arrowhead="normal", arrowtail="none", arrowsize="0.6", weight="0.5"]
    subgraph deps {
    {{- range .List -}}
        {{ $path := .Path }}
        "{{ $path }}" [label="{{ .Label }}", color="{{ strokeColour . }}", fillcolor="{{ fillColour . }}", penwidth=1]
            {{- range .DepPaths }}
              "{{ $path }}" -> "{{.}}"
            {{- end }}
    {{- end }}
    }
}
