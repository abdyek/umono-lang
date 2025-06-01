package umonolang

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/umono-cms/umono-lang/interfaces"
	"github.com/umono-cms/umono-lang/internal/components"
	"github.com/umono-cms/umono-lang/internal/parameters"
)

type callTestSuite struct {
	suite.Suite
}

func (s *callTestSuite) TestAlreadyRead() {
	for sI, scene := range []struct {
		indexes [][2]int
		start   int
		end     int
		result  bool
	}{
		{
			[][2]int{
				[2]int{30, 40},
			},
			34,
			45,
			true,
		},
		{
			[][2]int{
				[2]int{30, 40},
				[2]int{5, 15},
			},
			1,
			3,
			false,
		},
		{
			[][2]int{
				[2]int{30, 40},
				[2]int{90, 120},
			},
			125,
			154,
			false,
		},
		{
			[][2]int{
				[2]int{30, 40},
				[2]int{90, 120},
			},
			50,
			60,
			false,
		},
	} {
		require.Equal(s.T(), scene.result, alreadyRead(scene.indexes, scene.start, scene.end), "scene index: "+strconv.Itoa(sI))
	}
}

func (s *callTestSuite) TestReadCalls() {
	for sI, scene := range []struct {
		content string
		comps   []interfaces.Component
		results []*call
	}{
		{
			"{{ ABC }}",
			[]interfaces.Component{
				components.NewCustom("ABC", "no-matter"),
			},
			[]*call{
				&call{components.NewCustom("ABC", "no-matter"), 0, 9, []interfaces.Argument{}},
			},
		},
		{
			"{{ ABC param-1=\"val-1\" param-2=\"val-2\" }}",
			[]interfaces.Component{
				components.NewCustomWithParams("ABC", "no-matter", []interfaces.Parameter{
					parameters.NewDynamicParam("param-1", "string", "val-1"),
					parameters.NewDynamicParam("param-2", "string", "val-2"),
				}),
			},
			[]*call{
				&call{components.NewCustomWithParams("ABC", "no-matter", []interfaces.Parameter{
					parameters.NewDynamicParam("param-1", "string", ""),
					parameters.NewDynamicParam("param-2", "string", ""),
				}), 0, 41, []interfaces.Argument{
					newArg("param-1", "string", "val-1"),
					newArg("param-2", "string", "val-2"),
				}},
			},
		},
		{
			"ABC XYZ",
			[]interfaces.Component{
				components.NewCustom("ABC", "no-matter"),
				components.NewCustom("XYZ", "no-matter"),
			},
			[]*call{
				&call{components.NewCustom("ABC", "no-matter"), 0, 3, []interfaces.Argument{}},
				&call{components.NewCustom("XYZ", "no-matter"), 4, 7, []interfaces.Argument{}},
			},
		},
		{
			"XYZ",
			[]interfaces.Component{
				components.NewCustomWithParams("XYZ", "no-matter", []interfaces.Parameter{
					parameters.NewDynamicParam("param", "string", "default-value"),
				}),
			},
			[]*call{
				&call{components.NewCustomWithParams("XYZ", "no-matter", []interfaces.Parameter{
					parameters.NewDynamicParam("param", "string", "default-value"),
				}), 0, 3, []interfaces.Argument{
					newArg("param", "string", "default-value"),
				}},
			},
		},
		{
			"{{ LINK url=\"https://umono.io\" text=\"click me!\" new-tab=true }}",
			[]interfaces.Component{
				components.NewCustomWithParams("LINK", "no-matter", []interfaces.Parameter{
					parameters.NewDynamicParam("url", "string", ""),
					parameters.NewDynamicParam("text", "string", ""),
					parameters.NewDynamicParam("new-tab", "bool", false),
				}),
			},
			[]*call{
				&call{
					components.NewCustomWithParams("LINK", "no-matter", []interfaces.Parameter{
						parameters.NewDynamicParam("url", "string", ""),
						parameters.NewDynamicParam("text", "string", ""),
						parameters.NewDynamicParam("new-tab", "bool", false),
					}), 0, 63, []interfaces.Argument{
						newArg("url", "string", "https://umono.io"),
						newArg("text", "string", "click me!"),
						newArg("new-tab", "bool", true),
					},
				},
			},
		},
		{
			"\\{{ ABC }}",
			[]interfaces.Component{
				components.NewCustom("ABC", "no-matter"),
			},
			[]*call{},
		},
		{
			"\\\\{{ ABC }}",
			[]interfaces.Component{
				components.NewCustom("ABC", "no-matter"),
			},
			[]*call{
				&call{components.NewCustom("ABC", "no-matter"), 0, 9, []interfaces.Argument{}},
			},
		},
		{
			"\\\\\\{{ ABC }}",
			[]interfaces.Component{
				components.NewCustom("ABC", "no-matter"),
			},
			[]*call{},
		},
	} {

		calls := readCalls(scene.content, scene.comps)

		errMsg := "scene index: " + strconv.Itoa(sI)

		require.Equal(s.T(), len(scene.results), len(calls), errMsg+" | Not equal scene results with calls.")

		for i := 0; i < len(scene.results); i++ {

			sCall := scene.results[i]
			call := calls[i]

			errMsg := errMsg + ", call index: " + strconv.Itoa(i)

			require.Equal(s.T(), sCall.Component().Name(), call.Component().Name(), errMsg)
			require.Equal(s.T(), sCall.Start(), call.Start(), errMsg)
			require.Equal(s.T(), sCall.End(), call.End(), errMsg)

			sArgs := sCall.Arguments()
			callArgs := call.Arguments()

			for j := 0; j < len(sArgs); j++ {
				errMsg := errMsg + ", param index: " + strconv.Itoa(j)
				require.Equal(s.T(), sArgs[j].Name(), callArgs[j].Name(), errMsg)
				if sArgs[j].Type() == "string" {
					require.Equal(s.T(), sArgs[j].Value().(string), callArgs[j].Value().(string), errMsg)
				} else if sArgs[j].Type() == "bool" {
					require.Equal(s.T(), sArgs[j].Value().(bool), callArgs[j].Value().(bool), errMsg)
				}
			}
		}
	}
}

func TestCallTestSuite(t *testing.T) {
	suite.Run(t, new(callTestSuite))
}
