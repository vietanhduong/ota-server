package profile

import "github.com/vietanhduong/ota-server/pkg/apis/v1/metadata"

func ConvertMetadataListToMap(input []*metadata.Metadata) map[string]string {
	result := make(map[string]string)
	for _, m := range input {
		result[m.Key] = m.Value
	}
	return result
}
