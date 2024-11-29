package internal_test

import (
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nickbryan/slogutil/internal"
)

func TestAttrGroupTreeWithAttrs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attrGroupTree        internal.AttrGroupTree
		attrs                []slog.Attr
		shouldReturnReceiver bool
		want                 []slog.Attr
	}{
		"calling WithAttrs on an empty AttrGroup returns nil": {
			attrGroupTree:        internal.NewAttrGroupTree(),
			shouldReturnReceiver: true,
			want:                 []slog.Attr{},
		},
		"calling WithAttrs with nil returns the receiver": {
			attrGroupTree:        internal.NewAttrGroupTree(),
			attrs:                nil,
			shouldReturnReceiver: true,
			want:                 []slog.Attr{},
		},
		"calling WithAttrs with nil on a AttrGroupTree that has existing attrs returns the receiver": {
			attrGroupTree:        internal.NewAttrGroupTree().WithAttrs([]slog.Attr{slog.String("a", "aVal"), slog.String("b", "bVal")}),
			attrs:                nil,
			shouldReturnReceiver: true,
			want:                 []slog.Attr{slog.String("a", "aVal"), slog.String("b", "bVal")},
		},
		"calling WithAttrs with an empty []slog.Attr returns the receiver": {
			attrGroupTree:        internal.NewAttrGroupTree(),
			attrs:                make([]slog.Attr, 0),
			shouldReturnReceiver: true,
			want:                 []slog.Attr{},
		},
		"calling WithAttrs with an empty []slog.Attr on a AttrGroupTree that has existing attrs returns the receiver": {
			attrGroupTree:        internal.NewAttrGroupTree().WithAttrs([]slog.Attr{slog.String("a", "aVal"), slog.String("b", "bVal")}),
			attrs:                make([]slog.Attr, 0),
			shouldReturnReceiver: true,
			want:                 []slog.Attr{slog.String("a", "aVal"), slog.String("b", "bVal")},
		},
		"calling WithAttrs with a []slog.Attr returns a new *AttrGroup with the passed Attrs": {
			attrGroupTree:        internal.NewAttrGroupTree(),
			attrs:                []slog.Attr{slog.String("a", "aVal"), slog.String("b", "bVal")},
			shouldReturnReceiver: false,
			want:                 []slog.Attr{slog.String("a", "aVal"), slog.String("b", "bVal")},
		},
		"calling WithAttrs with a []slog.Attr on a NewAttrGroupTree() that already has Attrs returns a new *AttrGroup with the passed Attrs concatenated with the existing Attrs": {
			attrGroupTree:        internal.NewAttrGroupTree().WithAttrs([]slog.Attr{slog.String("a", "aVal")}),
			attrs:                []slog.Attr{slog.String("b", "bVal"), slog.String("c", "cVal")},
			shouldReturnReceiver: false,
			want:                 []slog.Attr{slog.String("a", "aVal"), slog.String("b", "bVal"), slog.String("c", "cVal")},
		},
		"calling WithAttrs with a []slog.Attr on a NewAttrGroupTree() that already has grouped Attrs appends the Attrs to the group": {
			attrGroupTree:        internal.NewAttrGroupTree().WithGroup("myGroup").WithAttrs([]slog.Attr{slog.String("a", "aVal")}),
			attrs:                []slog.Attr{slog.String("b", "bVal")},
			shouldReturnReceiver: false,
			want:                 []slog.Attr{slog.Group("myGroup", slog.String("a", "aVal"), slog.String("b", "bVal"))},
		},
		"calling WithAttrs with duplicate Attr keys on the root marks duplicates accordingly": {
			attrGroupTree: internal.NewAttrGroupTree().WithAttrs([]slog.Attr{
				slog.String("a", "aVal"),
			}),
			attrs: []slog.Attr{
				slog.String("a", "aValDup"),
				slog.Int("a", 123),
			},
			shouldReturnReceiver: false,
			want: []slog.Attr{
				slog.String("a", "aVal"),
				slog.String("a#01", "aValDup"),
				slog.Int("a#02", 123),
			},
		},
		"calling WithAttrs with duplicate Attr group keys on the root marks duplicates accordingly": {
			attrGroupTree: internal.NewAttrGroupTree().WithAttrs([]slog.Attr{
				slog.Group("g", slog.String("a", "aVal")),
			},
			),
			attrs: []slog.Attr{
				slog.Group("g", slog.String("a", "aGVal")),
			},
			shouldReturnReceiver: false,
			want: []slog.Attr{
				slog.Group("g", slog.String("a", "aVal")),
				slog.Group("g#01", slog.String("a", "aGVal")),
			},
		},
		"calling WithAttrs with duplicate Attr keys in a group deduplicates accordingly": {
			attrGroupTree: internal.NewAttrGroupTree().WithAttrs([]slog.Attr{
				slog.Group("g", slog.String("a", "aVal"), slog.String("a", "aValDup")),
			}),
			attrs: []slog.Attr{
				slog.Group("g", slog.String("a", "aGVal"), slog.String("a", "aGValDup")),
			},
			shouldReturnReceiver: false,
			want: []slog.Attr{
				slog.Group("g", slog.String("a", "aVal"), slog.String("a#01", "aValDup")),
				slog.Group("g#01", slog.String("a", "aGVal"), slog.String("a#01", "aGValDup")),
			},
		},
		"calling WithAttrs with many duplicate Attr deduplicates accordingly": {
			attrGroupTree: internal.NewAttrGroupTree().WithAttrs([]slog.Attr{
				slog.String("a", "v"),
				slog.Int("a", 123),
				slog.String("a", "v2"),
				slog.Group("a", slog.String("a", "v")),
			}),
			attrs: []slog.Attr{
				slog.String("a", "v"),
				slog.Int("a", 123),
				slog.String("a", "v2"),
				slog.Group("a", slog.String("a", "v")),
			},
			shouldReturnReceiver: false,
			want: []slog.Attr{
				slog.String("a", "v"),
				slog.Int("a#01", 123),
				slog.String("a#02", "v2"),
				slog.Group("a#03", slog.String("a", "v")),
				slog.String("a#04", "v"),
				slog.Int("a#05", 123),
				slog.String("a#06", "v2"),
				slog.Group("a#07", slog.String("a", "v")),
			},
		},
		"calling WithAttrs with an empty group drops the group": {
			attrGroupTree:        internal.NewAttrGroupTree(),
			attrs:                []slog.Attr{slog.Group("g")},
			shouldReturnReceiver: false,
			want:                 []slog.Attr{},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			got := testCase.attrGroupTree.WithAttrs(testCase.attrs)

			if testCase.shouldReturnReceiver {
				if !cmp.Equal(testCase.attrGroupTree, got, cmp.AllowUnexported(internal.AttrGroup{}, internal.AttrGroupTree{})) {
					t.Errorf("calling WithAttrs(%+[1]v), got: %[2]T(%+[2]v), want: %[3]T(%+[3]v)", testCase.attrs, got, testCase.attrGroupTree)
				}
			} else {
				if cmp.Equal(testCase.attrGroupTree, got, cmp.AllowUnexported(internal.AttrGroup{}, internal.AttrGroupTree{})) {
					t.Errorf("calling WithAttrs(%+[1]v), got: %[2]T(%+[2]v), want: return value to != receiver)", testCase.attrs, got)
				}
			}

			if !cmp.Equal(testCase.want, got.History().DeduplicatedAttrs()) {
				t.Errorf("calling WithAttrs(%+[1]v).History().DeduplicatedAttrs(), got: %[2]T(%+[2]v), want: %[3]T(%+[3]v)", testCase.attrs, got.History().DeduplicatedAttrs(), testCase.want)
			}
		})
	}
}

func TestAttrGroupTreeWithGroup(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attrGroupTree        internal.AttrGroupTree
		groupName            string
		withAttrs            []slog.Attr
		shouldReturnReceiver bool
		want                 []slog.Attr
	}{
		"calling WithGroup with an empty groupName name returns the receiver": {
			attrGroupTree:        internal.NewAttrGroupTree(),
			groupName:            "",
			shouldReturnReceiver: true,
		},
		"calling WithGroup on the root adds all future calls to WithAttrs to the group": {
			attrGroupTree:        internal.NewAttrGroupTree(),
			groupName:            "rootGroup",
			withAttrs:            []slog.Attr{slog.String("a", "aVal")},
			shouldReturnReceiver: false,
			want:                 []slog.Attr{slog.Group("rootGroup", slog.String("a", "aVal"))},
		},
		"calling WithGroup on the root with an empty group name adds all future calls to WithAttrs to the root": {
			attrGroupTree:        internal.NewAttrGroupTree(),
			groupName:            "",
			withAttrs:            []slog.Attr{slog.String("a", "aVal")},
			shouldReturnReceiver: false,
			want:                 []slog.Attr{slog.String("a", "aVal")},
		},
		"calling WithGroup on a group and not adding attributes drops the outer group": {
			attrGroupTree:        internal.NewAttrGroupTree().WithGroup("groupA"),
			groupName:            "groupB",
			withAttrs:            nil,
			shouldReturnReceiver: false,
			want:                 []slog.Attr{slog.Group("groupA")},
		},
		"calling WithGroup on a group and then adding attributes nests the groups": {
			attrGroupTree:        internal.NewAttrGroupTree().WithGroup("groupA"),
			groupName:            "groupB",
			withAttrs:            []slog.Attr{slog.String("a", "aVal")},
			shouldReturnReceiver: false,
			want:                 []slog.Attr{slog.Group("groupA", slog.Group("groupB", slog.String("a", "aVal")))},
		},
		"calling WithGroup with a group name that duplicates an existing attribute gets marked accordingly": {
			attrGroupTree:        internal.NewAttrGroupTree().WithAttrs([]slog.Attr{slog.String("a", "aVal")}),
			groupName:            "a",
			withAttrs:            []slog.Attr{slog.String("a", "aVal")},
			shouldReturnReceiver: false,
			want:                 []slog.Attr{slog.String("a", "aVal"), slog.Group("a#01", slog.String("a", "aVal"))},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			got := testCase.attrGroupTree.WithGroup(testCase.groupName)

			if testCase.withAttrs != nil {
				got = got.WithAttrs(testCase.withAttrs)
			}

			if testCase.shouldReturnReceiver {
				if !cmp.Equal(testCase.attrGroupTree, got, cmp.AllowUnexported(internal.AttrGroup{}, internal.AttrGroupTree{})) {
					t.Errorf("calling WithGroup(%+[1]v), got: %[2]T(%+[2]v), want: %[3]T(%+[3]v)", testCase.groupName, got, testCase.attrGroupTree)
				}
			} else {
				if cmp.Equal(testCase.attrGroupTree, got, cmp.AllowUnexported(internal.AttrGroup{}, internal.AttrGroupTree{})) {
					t.Errorf("calling WithGroup(%+[1]v), got: %[2]T(%+[2]v), want: return value to != receiver)", testCase.groupName, got)
				}

				if !cmp.Equal(testCase.want, got.History().DeduplicatedAttrs()) {
					t.Errorf("calling WithGroup(%+[1]v).History().DeduplicatedAttrs(), got: %[2]T(%+[2]v), want: %[3]T(%+[3]v)", testCase.groupName, got.History().DeduplicatedAttrs(), testCase.want)
				}
			}
		})
	}
}

func TestAttrGroupHistoryPushFront(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attrGroupHistory *internal.AttrGroupHistory
		pushAttrs        []slog.Attr
		want             []slog.Attr
	}{
		"calling PushFront with nil attrs does nothing": {
			attrGroupHistory: internal.NewAttrGroupTree().History(),
			pushAttrs:        nil,
			want:             []slog.Attr{},
		},
		"calling PushFront with empty attrs does nothing": {
			attrGroupHistory: internal.NewAttrGroupTree().History(),
			pushAttrs:        []slog.Attr{},
			want:             []slog.Attr{},
		},
		"calling PushFront with nil attrs on a group with attrs does nothing": {
			attrGroupHistory: internal.NewAttrGroupTree().WithAttrs([]slog.Attr{slog.String("aK", "aV")}).History(),
			pushAttrs:        nil,
			want:             []slog.Attr{slog.String("aK", "aV")},
		},
		"calling PushFront with empty attrs on a group with attrs does nothing": {
			attrGroupHistory: internal.NewAttrGroupTree().WithAttrs([]slog.Attr{slog.String("aK", "aV")}).History(),
			pushAttrs:        nil,
			want:             []slog.Attr{slog.String("aK", "aV")},
		},
		"calling PushFront with attrs on a new AttrGroupTree adds the attrs to the group": {
			attrGroupHistory: internal.NewAttrGroupTree().History(),
			pushAttrs:        []slog.Attr{slog.String("ak", "aV"), slog.Int("bK", 123)},
			want:             []slog.Attr{slog.String("ak", "aV"), slog.Int("bK", 123)},
		},
		"calling PushFront with attrs on an AttrGroupTree with existing attrs adds the attrs to the front of the group": {
			attrGroupHistory: internal.NewAttrGroupTree().WithAttrs([]slog.Attr{slog.String("eK", "eV")}).History(),
			pushAttrs:        []slog.Attr{slog.String("ak", "aV"), slog.Int("bK", 123)},
			want:             []slog.Attr{slog.String("ak", "aV"), slog.Int("bK", 123), slog.String("eK", "eV")},
		},
		"calling PushFront with attrs on an AttrGroupTree with multiple groups adds the attrs to the front of the group first group": {
			attrGroupHistory: internal.NewAttrGroupTree().WithAttrs([]slog.Attr{slog.String("rK", "rV")}).WithGroup("g1").WithAttrs([]slog.Attr{slog.Int("g1K", 123)}).History(),
			pushAttrs:        []slog.Attr{slog.String("ak", "aV"), slog.Int("bK", 123)},
			want:             []slog.Attr{slog.String("ak", "aV"), slog.Int("bK", 123), slog.String("rK", "rV"), slog.Group("g1", slog.Int("g1K", 123))},
		},
		"calling PushFront on a cloned AttrGroupTree does not persist previously pushed values": {
			attrGroupHistory: (func() *internal.AttrGroupHistory {
				agt1 := internal.NewAttrGroupTree()
				agt1.History().PushFront([]slog.Attr{slog.String("k", "v")})
				agt2 := agt1.WithAttrs([]slog.Attr{slog.String("k2", "v2")})

				return agt2.History()
			})(),
			pushAttrs: []slog.Attr{slog.String("k3", "v3")},
			want:      []slog.Attr{slog.String("k3", "v3"), slog.String("k2", "v2")},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			testCase.attrGroupHistory.PushFront(testCase.pushAttrs)

			got := testCase.attrGroupHistory.DeduplicatedAttrs()

			if !cmp.Equal(testCase.want, got) {
				t.Errorf("calling attrGroupHistory.DuplicateMarkedAttrs(), got: %v, want: %+v", got, testCase.want)
			}
		})
	}
}
