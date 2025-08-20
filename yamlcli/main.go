package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func main() {
	var (
		key        string
		configFile = flag.String("f", "", "config file path")
		help       = flag.Bool("h", false, "help information")
	)

	flag.Usage = func() {
		fmt.Printf(`Usage:
	yamlcli -f <config_file> get <key>
	yamlcli -f <config_file> set <key> <value>
	yamlcli -f <config_file> del <key>

Options:
	-f		Config file path
	-h		Show help information

Examples:
	yamlcli -f config.yaml get # get all configuration
	yamlcli -f config.yaml get database.host # get key
	yamlcli -f config.yaml set database.enabled true # set boolean value
	yamlcli -f config.yaml set database.user "root" # set string value
	yamlcli -f config.yaml set database.port 3306 # set number value
	yamlcli -f config.yaml set servers "[server1,server2,server3]" # set array value
	yamlcli -f config.yaml set database "{host:localhost,port:3306,name:mydb}" # set object value
	yamlcli -f config.yaml del database.password # delete key
`)
	}

	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	if *configFile == "" {
		fmt.Fprintf(os.Stderr, "Error: Config file path is required\n")
		flag.Usage()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Must specify operation (get, set, or delete)\n\n")
		flag.Usage()
		return
	}

	operation := strings.ToLower(args[0])

	switch operation {
	case "get":
		if len(args) < 2 {
			key = ""
		} else {
			key = args[1]
		}

		value := GetConfig(*configFile, key)
		if value == nil {
			fmt.Printf("")
			return
		}

		fmt.Println(value)

	case "set":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s -f <config_file> set <key> <value>\n", os.Args[0])
			return
		}

		key := args[1]
		rawValue := args[2]

		// 智能转换值类型
		value := convertValue(rawValue)
		err := SetConfig(*configFile, key, value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to set configuration: %v\n", err)
			return
		}

		fmt.Printf("Successfully set '%s'\n", key)

	case "del":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: %s -f <config_file> delete <key>\n", os.Args[0])
			return
		}

		key := args[1]
		err := DeleteConfig(*configFile, key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete configuration: %v\n", err)
			return
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown operation '%s'\n", operation)
		flag.Usage()
		return
	}
}

func GetConfig(configFile string, key string) interface{} {
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	if key == "" {
		return v.AllSettings()
	}
	return v.Get(key)
}

func SetConfig(configFile string, key string, value interface{}) error {
	// 读取原始文件内容
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果文件不存在，创建新文件
			return createNewConfig(configFile, key, value)
		}
		return fmt.Errorf("reading config file failed: %w", err)
	}

	// 解析 YAML 保持格式
	var root yaml.Node
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return fmt.Errorf("parsing yaml failed: %w", err)
	}

	// 设置值
	err = setValueInNode(&root, key, value)
	if err != nil {
		return fmt.Errorf("setting value failed: %w", err)
	}

	// 写回文件，保持原格式
	output, err := yaml.Marshal(&root)
	if err != nil {
		return fmt.Errorf("marshaling yaml failed: %w", err)
	}

	err = os.WriteFile(configFile, output, 0644)
	if err != nil {
		return fmt.Errorf("writing config file failed: %w", err)
	}

	return nil
}

func DeleteConfig(configFile string, key string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("reading config file failed: %w", err)
	}

	var root yaml.Node
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return fmt.Errorf("parsing yaml failed: %w", err)
	}

	err = deleteValueInNode(&root, key)
	if err != nil {
		return fmt.Errorf("deleting value failed: %w", err)
	}
	output, err := yaml.Marshal(&root)
	if err != nil {
		return fmt.Errorf("marshaling yaml failed: %w", err)
	}

	err = os.WriteFile(configFile, output, 0644)
	if err != nil {
		return fmt.Errorf("writing config file failed: %w", err)
	}

	return nil
}

func deleteValueInNode(root *yaml.Node, key string) error {
	if root.Kind != yaml.DocumentNode {
		return fmt.Errorf("expected document node")
	}

	if len(root.Content) == 0 {
		return fmt.Errorf("document is empty")
	}

	mappingNode := root.Content[0]
	if mappingNode.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node")
	}

	keys := strings.Split(key, ".")
	return deleteNestedValueInNode(mappingNode, keys)
}
