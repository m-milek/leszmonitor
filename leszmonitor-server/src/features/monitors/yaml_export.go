package monitors

import (
	"context"
	"io"
	"slices"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/probe"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/config"
	"gopkg.in/yaml.v3"
)

func ExportToYAML(ctx context.Context, monitor Monitor, tagsDAO tags.ITagDAO, writer io.Writer) error {
	monitorNode, err := buildMonitorNode(monitor)
	if err != nil {
		return err
	}

	tagNames, err := getTagNames(ctx, monitor.TagIDs, tagsDAO)
	if err != nil {
		return err
	}
	if len(tagNames) > 0 {
		tagsNode, err := buildTagsNode(tagNames)
		if err != nil {
			return err
		}
		monitorNode.Content = append(monitorNode.Content, stringNode("tags"), tagsNode)
	}

	rootNode := yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{stringNode(monitor.Slug), monitorNode},
	}

	return yaml.NewEncoder(writer).Encode(&rootNode)
}

func buildMonitorNode(monitor Monitor) (*yaml.Node, error) {
	monitorNode := yaml.Node{}
	if err := monitorNode.Encode(monitor); err != nil {
		return nil, err
	}

	probeNode, err := buildProbeNode(monitor)
	if err != nil {
		return nil, err
	}

	var content []*yaml.Node
	for pair := range slices.Chunk(monitorNode.Content, 2) {
		key, value := pair[0], pair[1]
		if slices.Contains(config.ForbiddenMonitorConfigKeys, key.Value) {
			continue
		}
		if key.Value == "probeConfig" {
			value = probeNode
		}
		content = append(content, key, value)
	}
	monitorNode.Content = content

	return &monitorNode, nil
}

func buildProbeNode(monitor Monitor) (*yaml.Node, error) {
	parsedProbe, err := probe.Parse[probe.Probe](monitor.Type, monitor.ProbeConfig)
	if err != nil {
		return nil, err
	}

	probeNode := yaml.Node{}
	if err = probeNode.Encode(parsedProbe); err != nil {
		return nil, err
	}

	return &probeNode, nil
}

func buildTagsNode(tagNames []string) (*yaml.Node, error) {
	tagsNode := yaml.Node{}
	if err := tagsNode.Encode(tagNames); err != nil {
		return nil, err
	}
	tagsNode.Style = yaml.FlowStyle

	return &tagsNode, nil
}

func stringNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

func getTagNames(ctx context.Context, tagIDs []uuid.UUID, tagsDAO tags.ITagDAO) ([]string, error) {
	allTags, err := tagsDAO.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}

	tagsByID := make(map[uuid.UUID]tags.Tag, len(allTags))
	for _, tag := range allTags {
		tagsByID[tag.ID] = tag
	}

	var tagNames []string
	for _, tagID := range tagIDs {
		if tag, found := tagsByID[tagID]; found {
			tagNames = append(tagNames, tag.Name)
		}
	}

	return tagNames, nil
}
