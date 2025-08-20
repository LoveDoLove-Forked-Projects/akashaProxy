package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// yaml的读写由claude4和GPT4.1大力支持(

func setNestedValue(config map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, ".")
	current := config
	for i := 0; i < len(keys)-1; i++ {
		k := keys[i]
		if _, exists := current[k]; !exists {
			current[k] = make(map[string]interface{})
		}
		if nested, ok := current[k].(map[string]interface{}); ok {
			current = nested
		} else {
			current[k] = make(map[string]interface{})
			current = current[k].(map[string]interface{})
		}
	}

	current[keys[len(keys)-1]] = value
}

func createNewConfig(configFile string, key string, value interface{}) error {
	config := make(map[string]interface{})
	setNestedValue(config, key, value)

	output, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshaling yaml failed: %w", err)
	}

	err = os.WriteFile(configFile, output, 0644)
	if err != nil {
		return fmt.Errorf("writing config file failed: %w", err)
	}

	return nil
}

func setValueInNode(root *yaml.Node, key string, value interface{}) error {
	if root.Kind != yaml.DocumentNode {
		return fmt.Errorf("expected document node")
	}

	if len(root.Content) == 0 {
		root.Content = append(root.Content, &yaml.Node{
			Kind: yaml.MappingNode,
		})
	}

	mappingNode := root.Content[0]
	if mappingNode.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node")
	}

	keys := strings.Split(key, ".")
	return setNestedValueInNode(mappingNode, keys, value)
}

func setNestedValueInNode(node *yaml.Node, keys []string, value interface{}) error {
	if len(keys) == 0 {
		return fmt.Errorf("empty keys")
	}

	currentKey := keys[0]
	remainingKeys := keys[1:]

	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}

		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Value == currentKey {
			if len(remainingKeys) == 0 {
				*valueNode = *createValueNode(value)
				return nil
			} else {
				if valueNode.Kind != yaml.MappingNode {
					*valueNode = yaml.Node{Kind: yaml.MappingNode}
				}
				return setNestedValueInNode(valueNode, remainingKeys, value)
			}
		}
	}

	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: currentKey,
	}

	var valueNode *yaml.Node
	if len(remainingKeys) == 0 {
		valueNode = createValueNode(value)
	} else {
		valueNode = &yaml.Node{Kind: yaml.MappingNode}
		err := setNestedValueInNode(valueNode, remainingKeys, value)
		if err != nil {
			return err
		}
	}

	node.Content = append(node.Content, keyNode, valueNode)
	return nil
}

func createValueNode(value interface{}) *yaml.Node {
	switch v := value.(type) {
	case []interface{}:
		node := &yaml.Node{Kind: yaml.SequenceNode}
		for _, item := range v {
			node.Content = append(node.Content, createValueNode(item))
		}
		return node
	case []string:
		node := &yaml.Node{Kind: yaml.SequenceNode}
		for _, item := range v {
			node.Content = append(node.Content, createValueNode(item))
		}
		return node
	case []int:
		node := &yaml.Node{Kind: yaml.SequenceNode}
		for _, item := range v {
			node.Content = append(node.Content, createValueNode(item))
		}
		return node
	case map[string]interface{}:
		node := &yaml.Node{Kind: yaml.MappingNode}
		for k, val := range v {
			keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: k}
			valueNode := createValueNode(val)
			node.Content = append(node.Content, keyNode, valueNode)
		}
		return node
	case bool:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%t", v)}
	case int:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%d", v)}
	case float64:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%g", v)}
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: v}
	default:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", v)}
	}
}

func convertValue(value string) interface{} {
	value = strings.TrimSpace(value)

	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		arrayContent := strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
		if arrayContent == "" {
			return []interface{}{}
		}

		items := smartSplit(arrayContent, ",")
		result := make([]interface{}, 0, len(items))

		for _, item := range items {
			item = strings.TrimSpace(item)
			result = append(result, convertValue(item))
		}

		return result
	}

	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		objContent := strings.TrimPrefix(strings.TrimSuffix(value, "}"), "{")
		if objContent == "" {
			return make(map[string]interface{})
		}

		result := make(map[string]interface{})
		pairs := smartSplit(objContent, ",")

		for _, pair := range pairs {
			kv := smartSplit(strings.TrimSpace(pair), ":")
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				if strings.HasPrefix(key, "\"") && strings.HasSuffix(key, "\"") {
					key = strings.Trim(key, "\"")
				}
				result[key] = convertValue(val)
			}
		}

		return result
	}

	if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
		return strings.Trim(value, "\"'")
	}
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	if strings.Contains(value, ".") {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	} else {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return value
}

func smartSplit(s, sep string) []string {
	var result []string
	var current strings.Builder
	var depth int
	var inQuotes bool
	var quoteChar byte

	for i := 0; i < len(s); i++ {
		char := s[i]

		if (char == '"' || char == '\'') && !inQuotes {
			inQuotes = true
			quoteChar = char
		} else if char == quoteChar && inQuotes {
			inQuotes = false
		}

		if inQuotes {
			current.WriteByte(char)
			continue
		}

		if char == '[' || char == '{' || char == '(' {
			depth++
		} else if char == ']' || char == '}' || char == ')' {
			depth--
		}

		if depth == 0 && strings.HasPrefix(s[i:], sep) {
			result = append(result, current.String())
			current.Reset()
			i += len(sep) - 1
			continue
		}

		current.WriteByte(char)
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func deleteNestedValueInNode(node *yaml.Node, keys []string) error {
	if len(keys) == 0 {
		return fmt.Errorf("empty keys")
	}

	currentKey := keys[0]
	remainingKeys := keys[1:]

	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}

		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Value == currentKey {
			if len(remainingKeys) == 0 {
				node.Content = append(node.Content[:i], node.Content[i+2:]...)
				return nil
			} else {
				if valueNode.Kind != yaml.MappingNode {
					return fmt.Errorf("cannot delete nested key from non-mapping node")
				}
				return deleteNestedValueInNode(valueNode, remainingKeys)
			}
		}
	}

	return fmt.Errorf("key '%s' not found", currentKey)
}
