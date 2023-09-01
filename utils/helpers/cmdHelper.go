package helpers

import (
	"github.com/spf13/cobra"
	"log"
)

// MustFlags 获取flag参数值
func MustFlags(cmd *cobra.Command, key string, valueType string) interface{} {
	switch valueType {
	case "string":
		if v, err := cmd.Flags().GetString(key); err == nil && v != "" {
			return v
		}
	case "int":
		if v, err := cmd.Flags().GetInt(key); err == nil {
			return v
		}
	}

	log.Fatalln(key, "is required")
	return nil
}
