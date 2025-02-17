package dtsfmt

const (
	// 完整定义见 https://github.com/joelspadin/tree-sitter-devicetree/blob/main/src/node-types.json
	// 这里仅列举需要处理格式化的节点类型

	NodeKindColon          = ":"
	NodeKindComma          = ","
	NodeKindEq             = "="
	NodeKindSemiColon      = ";"
	NodeKindLeftBracket    = "{"
	NodeKindRightBracket   = "}"
	NodeKindComment        = "comment"
	NodeKindFileVersion    = "file_version"
	NodeKindNode           = "node"
	NodeKindProperty       = "property"
	NodeKindStringLiteral  = "string_literal"
	NodeKindIntegerCells   = "integer_cells"
	NodeKindPreprocInclude = "preproc_include"
)
