package emitter

import (
	"fmt"
	"strings"

	"sov.fleet/extendgo/81000-active-source/spec_ingest"
)

// EmitWebNFSN formats a slice of RawRules into a compliant .wag file conforming to weBNF.sn.
func EmitWebNFSN(rules []spec_ingest.RawRule, goVersion string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Wirth Augmented Grammar (WAG) for Sovereign Go v%s\n", goVersion))
	sb.WriteString("# Domain: Systems, Supercomputing & High-Precision Science\n")
	sb.WriteString("# Grammar Format: context-sensitive weBNF.sn (:com:velocikey:languagegrammar:go)\n\n")
	sb.WriteString(":com:velocikey:languagegrammar:go\n\n")

	for _, r := range rules {
		formatSingleRule(&sb, r)
	}

	return sb.String()
}

func formatSingleRule(sb *strings.Builder, r spec_ingest.RawRule) {
	switch r.Name {
	case "DeclStmt":
		sb.WriteString("DeclStmt = \"type\" Identifier [ TypeParamList ] ( StructType / InterfaceType / DataType )\n")
		sb.WriteString("         / \"var\" Identifier [ TypeParamList ] [ DataType ] [ \"=\" Expr ]\n")
		sb.WriteString("         / \"const\" Identifier [ DataType ] [ \"=\" Expr ]\n")
		sb.WriteString("         / Identifier [ TypeParamList ] [ \":\" DataType ] [ \"=\" Expr ] ;\n\n")

	case "FuncDecl":
		sb.WriteString("# Go 1.27 Generic Function and Method Declarations with Visibility Suffix\n")
		sb.WriteString(fmt.Sprintf("FuncDecl = %s ;\n\n", r.Body))

	case "TypeParamList":
		sb.WriteString("# Generic Type Parameters\n")
		sb.WriteString("TypeParamList = \"[\" TypeParamDecl { \",\" TypeParamDecl } \"]\" ;\n\n")

	case "InterfaceType":
		sb.WriteString("# Go 1.27 Interface Types (Enforcing negative constraint: no TypeParamList on InterfaceMethod)\n")
		sb.WriteString("InterfaceType = \"interface\" \"{\" { InterfaceElem \";\" } \"}\" ;\n\n")

	case "StructType":
		sb.WriteString("# Struct Types and Go 1.27 Field Selectors\n")
		sb.WriteString("StructType = \"struct\" \"{\" { FieldDecl \";\" } \"}\" ;\n\n")

	case "DataType":
		sb.WriteString("# Data Types (including pointers and slices for receivers and high-precision scientific types)\n")
		sb.WriteString("DataType = [ \"*\" / \"[]\" ] BaseType ;\n\n")

	case "BaseType":
		sb.WriteString("BaseType = \"float2048\" / \"float1024\" / \"float512\" / \"float256\" / \"float128\" / \"float64\" / \"float32\"\n")
		sb.WriteString("         / \"complex4096\" / \"complex2048\" / \"complex1024\" / \"complex512\" / \"complex256\" / \"complex128\" / \"complex64\"\n")
		sb.WriteString("         / \"int2048\" / \"int1024\" / \"int512\" / \"int256\" / \"int128\" / \"int64\" / \"int32\" / \"int16\" / \"int8\" / \"int\"\n")
		sb.WriteString("         / \"uint2048\" / \"uint1024\" / \"uint512\" / \"uint256\" / \"uint128\" / \"uint64\" / \"uint32\" / \"uint16\" / \"uint8\" / \"uint\"\n")
		sb.WriteString("         / DecimalType\n")
		sb.WriteString("         / GenericType\n")
		sb.WriteString("         / \"string\" / \"bool\" / \"byte\" / \"rune\" / \"error\" / Identifier ;\n\n")

	case "Expr":
		sb.WriteString("# Expressions and Primary Expressions with Go 1.27 Generic Method Invocations\n")
		sb.WriteString("Expr = PrimaryExpr / CompositeLit / Number / StringLit / Identifier ;\n\n")

	case "CompositeLit":
		sb.WriteString("# Go 1.27 Struct Literals with Field Selectors\n")
		sb.WriteString("CompositeLit = ( Identifier / GenericType / StructType ) LiteralValue ;\n\n")

	case "Keyword":
		sb.WriteString("Keyword = \"break\" / \"case\" / \"chan\" / \"const\" / \"continue\" / \"default\" / \"defer\" / \"else\"\n")
		sb.WriteString("        / \"fallthrough\" / \"for\" / \"func\" / \"go\" / \"goto\" / \"if\" / \"import\" / \"interface\"\n")
		sb.WriteString("        / \"map\" / \"package\" / \"range\" / \"return\" / \"select\" / \"struct\" / \"switch\"\n")
		sb.WriteString("        / \"type\" / \"var\" / \"decimal\" / \"public\" / \"private\" ;\n\n")

	default:
		sb.WriteString(fmt.Sprintf("%s = %s ;\n\n", r.Name, r.Body))
	}
}
