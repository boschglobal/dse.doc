package ast

type YAMLVisitor interface {
	Visit(map[string]interface{}, *Index)
}

type InnerVisitor struct {
	YAMLVisitor
}

func (v InnerVisitor) Visit(node map[string]interface{}, md_doc *Index) {
	innerNodes, ok := node["inner"].([]interface{})
	if !ok {
		return
	}

	for _, innerNode := range innerNodes {
		if innerNodeMap, ok := innerNode.(map[string]interface{}); ok {
			v.YAMLVisitor.Visit(innerNodeMap, md_doc)
		}
	}
}

type FunctionDeclVisitor struct{}

func (v FunctionDeclVisitor) Visit(node map[string]interface{}, md_doc *Index) {
	kind, kindOk := node["kind"].(string)
	if !kindOk || kind != "FunctionDecl" {
		return
	}

	name, nameOk := node["name"].(string)
	if nameOk {
		md_doc.Functions = append(md_doc.Functions, name)
	}
}

type TypedefDeclVisitor struct {
	TypeList []string
}

func (v *TypedefDeclVisitor) Visit(node map[string]interface{}, md_doc *Index) {
	kind, kindOk := node["kind"].(string)
	if !kindOk || kind != "TypedefDecl" {
		return
	}

	name, nameOk := node["name"].(string)
	if nameOk {
		v.TypeList = append(v.TypeList, name)
	}
}

type RecordDeclVisitor struct {
	TypeList []string
}

func (v RecordDeclVisitor) Visit(node map[string]interface{}, md_doc *Index) {
	kind, kindOk := node["kind"].(string)
	if !kindOk || kind != "RecordDecl" {
		return
	}

	name, nameOk := node["name"].(string)
	if !nameOk || !itemExists(v.TypeList, name) {
		return
	}

	inner, innerOk := node["inner"].([]interface{})
	if !innerOk {
		return
	}

	for _, innerItem := range inner {
		if innerItemMap, innerItemOk := innerItem.(map[string]interface{}); innerItemOk {
			v.processFieldDecl(name, innerItemMap, md_doc)
		}
	}
}

func (v RecordDeclVisitor) processFieldDecl(typedefName string, innerItemMap map[string]interface{}, md_doc *Index) {
	innerKind, innerKindOk := innerItemMap["kind"].(string)
	if !innerKindOk || innerKind != "FieldDecl" {
		return
	}

	fieldName, fieldNameOk := innerItemMap["name"].(string)
	fieldTypeMap, fieldTypeMapOk := innerItemMap["type"].(map[string]interface{})
	if !fieldNameOk || !fieldTypeMapOk {
		return
	}

	fieldType, fieldTypeOk := fieldTypeMap["qualType"].(string)
	if fieldTypeOk {
		typeVar := fieldType + " " + fieldName
		md_doc.Typedefs[typedefName] = append(md_doc.Typedefs[typedefName], typeVar)
	}
}

func itemExists(list []string, targetItem string) bool {
	for _, item := range list {
		if item == targetItem {
			return true
		}
	}
	return false
}
