package internal

import (
	"fmt"
	"log/slog"
	"strings"
)

type (
	// AttrGroup is an immutable collection of [slog.Attr] values that will all be
	// qualified by the group name. The path is used to track nesting of groups from the root.
	AttrGroup struct {
		name, path string
		attrs      []slog.Attr
	}

	// AttrGroupTree is an immutable tree structure of one or more [AttrGroup] objects.
	AttrGroupTree struct {
		AttrGroup
		ancestor *AttrGroupTree
	}

	// AttrGroupHistory is a set of [AttrGroup] objects that are in historical order from the
	// oldest ancestor group (the root) through to the most recent descendant group. It also handles
	// the tracking and deduplication of attr keys.
	AttrGroupHistory struct {
		groups            []AttrGroup
		duplicateAttrKeys map[string]int
	}
)

// NewAttrGroupTree creates an empty [AttrGroupTree].
func NewAttrGroupTree() AttrGroupTree {
	return AttrGroupTree{
		AttrGroup: AttrGroup{
			name:  "",
			path:  "",
			attrs: make([]slog.Attr, 0),
		},
		ancestor: nil,
	}
}

// WithAttrs returns a new copy of the [AttrGroupTree] with the given attributes added
// to the current [AttrGroup].
func (agt AttrGroupTree) WithAttrs(attrs []slog.Attr) AttrGroupTree {
	if len(attrs) == 0 {
		return agt
	}

	return AttrGroupTree{
		AttrGroup: AttrGroup{
			name:  agt.name,
			path:  agt.path,
			attrs: append(agt.attrs, attrs...),
		},
		ancestor: agt.ancestor,
	}
}

// WithGroup returns a new [AttrGroupTree], with the previous [AttrGroupTree] as its
// ancestor, creating a new [AttrGroup] with the given group name.
func (agt AttrGroupTree) WithGroup(name string) AttrGroupTree {
	if name == "" {
		return agt
	}

	return AttrGroupTree{
		AttrGroup: AttrGroup{
			name:  name,
			path:  groupPath(agt.path, name),
			attrs: nil,
		},
		ancestor: &agt,
	}
}

// History returns a new [AttrGroupHistory] for the given [AttrGroupTree].
func (agt AttrGroupTree) History() *AttrGroupHistory {
	numberOfGroups := 0
	for group := &agt; group != nil; group = group.ancestor {
		numberOfGroups++
	}

	groups := make([]AttrGroup, numberOfGroups)
	for group := &agt; group != nil; group = group.ancestor {
		groups[numberOfGroups-1] = group.AttrGroup
		numberOfGroups--
	}

	return &AttrGroupHistory{groups: groups, duplicateAttrKeys: make(map[string]int)}
}

// PushFront adds the given attrs to the beginning of the list of attrs
// for the root group of the [AttrGroupHistory].
func (agh *AttrGroupHistory) PushFront(attrs []slog.Attr) {
	if len(agh.groups) == 0 || len(attrs) == 0 {
		return
	}

	agh.groups[0].attrs = append(attrs, agh.groups[0].attrs...)
}

// DeduplicatedAttrs returns the flattened slice of [slog.Attr] with attrs
// properly nested within the desired groups. Where there are multiple attrs at
// the same group level with the same key, the first attr's key will be left as is
// and every subsequent duplicate attr's key will be suffixed with #0x
// incrementally. This logic also applies to groups.
func (agh *AttrGroupHistory) DeduplicatedAttrs() []slog.Attr {
	return agh.resolve()
}

// resolve returns the [AttrGroupHistory] as a flattened slice of resolved
// [slog.Attr] values, qualified by all applicable group names and ready for a
// [slog.Handler].
func (agh *AttrGroupHistory) resolve() []slog.Attr {
	if len(agh.groups) == 0 {
		return nil
	}

	resolvedAttrs := agh.resolveAttrs(agh.groups[0].path, agh.groups[0].attrs)

	if len(agh.groups) > 1 {
		descendentGroups := &AttrGroupHistory{groups: agh.groups[1:], duplicateAttrKeys: agh.duplicateAttrKeys}
		resolvedAttrs = append(resolvedAttrs, descendentGroups.resolve()...)
	}

	if agh.groups[0].name == "" {
		return resolvedAttrs
	}

	key := agh.groups[0].name

	pathWithKey := groupPath(agh.groups[0].path, key)
	if _, ok := agh.duplicateAttrKeys[pathWithKey]; ok {
		agh.duplicateAttrKeys[pathWithKey]++
		key = duplicateKey(key, agh.duplicateAttrKeys[pathWithKey])
	} else {
		agh.duplicateAttrKeys[pathWithKey] = 0
	}

	return []slog.Attr{{
		Key:   key,
		Value: slog.GroupValue(resolvedAttrs...),
	}}
}

// resolveAttrs resolves the values of the given [slog.Attr] slice ready for handling
// by a [slog.Handler] (see [slog.LogValuer]). Empty [slog.Attr] values are ignored,
// groups without a group key are inlined, empty groups are ignored and named
// groups are recursively handled. We also track any duplicates here.
func (agh *AttrGroupHistory) resolveAttrs(path string, attrs []slog.Attr) []slog.Attr {
	if len(attrs) == 0 {
		return attrs
	}

	resolvedAttrs := make([]slog.Attr, 0, len(attrs))

	for _, attr := range attrs {
		if attr.Equal(slog.Attr{}) {
			continue
		}

		attr.Value = attr.Value.Resolve()

		pathWithKey := groupPath(path, attr.Key)
		if _, ok := agh.duplicateAttrKeys[pathWithKey]; ok {
			agh.duplicateAttrKeys[pathWithKey]++
		} else {
			agh.duplicateAttrKeys[pathWithKey] = 0
		}

		if attr.Value.Kind() != slog.KindGroup {
			if agh.duplicateAttrKeys[pathWithKey] > 0 {
				attr.Key = duplicateKey(attr.Key, agh.duplicateAttrKeys[pathWithKey])
			}

			resolvedAttrs = append(resolvedAttrs, attr)

			continue
		}

		if attr.Key == "" {
			resolvedAttrs = append(resolvedAttrs, agh.resolveAttrs(pathWithKey, attr.Value.Group())...)
			continue
		}

		if len(attr.Value.Group()) == 0 {
			continue
		}

		if agh.duplicateAttrKeys[pathWithKey] > 0 {
			attr.Key = duplicateKey(attr.Key, agh.duplicateAttrKeys[pathWithKey])
			pathWithKey = duplicateKey(pathWithKey, agh.duplicateAttrKeys[pathWithKey])
		}

		groupedAttrs := agh.resolveAttrs(pathWithKey, attr.Value.Group())
		resolvedAttrs = append(resolvedAttrs, slog.Attr{Key: attr.Key, Value: slog.GroupValue(groupedAttrs...)})
	}

	return resolvedAttrs
}

func groupPath(path, key string) string {
	const delimiter = "[.]"

	if strings.Contains(key, delimiter) {
		key = strings.ReplaceAll(key, delimiter, "___")
	}

	return strings.TrimLeft(path+delimiter+key, delimiter)
}

func duplicateKey(key string, duplicates int) string {
	return fmt.Sprintf("%s#%02d", key, duplicates)
}
