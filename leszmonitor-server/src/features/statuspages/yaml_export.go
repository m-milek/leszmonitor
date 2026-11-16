package statuspages

import (
	"context"
	"io"
	"slices"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/tags"
	"gopkg.in/yaml.v3"
)

func exportToYAML(ctx context.Context, page StatusPage, tagsDAO tags.ITagDAO, writer io.Writer) error {
	encodedNode := yaml.Node{}
	if err := encodedNode.Encode(page); err != nil {
		return err
	}

	pageNode := yaml.Node{Kind: yaml.MappingNode}
	for pair := range slices.Chunk(encodedNode.Content, 2) {
		if slices.Contains(forbiddenKeys, pair[0].Value) {
			continue
		}
		pageNode.Content = append(pageNode.Content, pair...)
	}

	tagNames, err := monitors.GetTagNames(ctx, page.TagIDs, tagsDAO)
	if err != nil {
		return err
	}
	tagsNode, err := monitors.BuildTagsNode(tagNames)
	if err != nil {
		return err
	}
	pageNode.Content = append(pageNode.Content, monitors.StringNode("tags"), tagsNode)

	rootNode := yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{monitors.StringNode(page.Slug), &pageNode},
	}

	return yaml.NewEncoder(writer).Encode(&rootNode)
}
