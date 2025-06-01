package umonolang

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/umono-cms/umono-lang/internal/utils/mocks"
	"github.com/umono-cms/umono-lang/internal/utils/test"
)

type umonoLangTestSuite struct {
	suite.Suite
	umonoLang *UmonoLang
}

func (s *umonoLangTestSuite) SetupTest() {
	s.umonoLang = New(new(mocks.Converter))
}

func (s *umonoLangTestSuite) TestConvert() {
	inputFileReader := test.NewFileReader("test_assets/main/inputs", "ul")
	outputFileReader := test.NewFileReader("test_assets/main/outputs", "mock")

	for _, scene := range []struct {
		file      string
		converter func(string) string
	}{
		{"000_empty", mocks.BasicConverter},
		{"001", mocks.BasicConverter},
		{"002", mocks.BasicConverter},
		{"003", mocks.BasicConverter},
		{"004", mocks.BasicConverter},
		{"005", mocks.BasicConverter},
		{"006", mocks.BasicConverter},
		{"007", mocks.BasicConverter},
		{"008", mocks.BasicConverter},
		{"009", mocks.BasicConverter},
		{"010", mocks.BasicConverter},
		{"011", func(ul string) string {
			return strings.ToUpper(ul)
		}},
		{"012", func(ul string) string {
			return strings.ToLower(ul)
		}},
		{"013", func(ul string) string {
			return strings.ReplaceAll(ul, "xyz", "abc")
		}},
		{"014", func(ul string) string {
			return strings.ReplaceAll(ul, "u", "a")
		}},
		{"015", func(ul string) string {
			return strings.ToUpper(ul)
		}},
		{"016", func(ul string) string {
			return strings.ToUpper(ul)
		}},
		{"017", func(ul string) string {
			return strings.ReplaceAll(ul, "-", "#")
		}},
		{"018", func(ul string) string {
			return ul
		}},
		{"019", func(ul string) string {
			return strings.ToUpper(ul)
		}},
		{"020", func(ul string) string {
			return strings.ToUpper(ul)
		}},
		{"021", func(ul string) string {
			return strings.ToLower(ul)
		}},
		{"022", func(ul string) string {
			return ul
		}},
		{"023", func(ul string) string {
			return ul
		}},
		{"024", func(ul string) string {
			return ul
		}},
		{"025", func(ul string) string {
			return ul
		}},
		{"026", func(ul string) string {
			return ul
		}},
		{"027", func(ul string) string {
			return ul
		}},
		{"028", func(ul string) string {
			return ul
		}},
	} {
		inputCont, err := inputFileReader.Read(scene.file, false)
		require.Nil(s.T(), err)

		outputCont, err := outputFileReader.Read(scene.file, true)
		require.Nil(s.T(), err)

		ul := New(mocks.NewMockConverter(scene.converter))

		assert.Equal(s.T(), outputCont, ul.Convert(inputCont), "input file name: "+scene.file+".ul")
	}
}

func (s *umonoLangTestSuite) TestSetGlobalComponentOK() {
	s.umonoLang.SetGlobalComponent("HELLO_WORLD", "hello!")

	require.Equal(s.T(), int(1), len(s.umonoLang.globalComps))

	hello := s.umonoLang.globalComps[0]

	require.Equal(s.T(), "hello!", hello.RawContent())
}

func (s *umonoLangTestSuite) TestSetGlobalComponentSyntaxError() {
	err := s.umonoLang.SetGlobalComponent("HELLO WORLD", "hello!")
	require.NotNil(s.T(), err)
	require.True(s.T(), strings.HasPrefix(err.Error(), "SYNTAX_ERROR"))
}

func (s *umonoLangTestSuite) TestRemoveGlobalComponentOK() {
	s.umonoLang.SetGlobalComponent("HELLO_WORLD", "hello!")
	err := s.umonoLang.RemoveGlobalComponent("HELLO_WORLD")
	require.Nil(s.T(), err)

	_, found := findCompByName(s.umonoLang.globalComps, "HELLO_WORLD")
	require.Nil(s.T(), found)
}

func (s *umonoLangTestSuite) TestRemoveGlobalComponentSyntaxError() {
	err := s.umonoLang.RemoveGlobalComponent("hello world")
	require.NotNil(s.T(), err)
	require.True(s.T(), strings.HasPrefix(err.Error(), "SYNTAX_ERROR"))
}

func (s *umonoLangTestSuite) TestRemoveGlobalComponentNotFound() {
	err := s.umonoLang.RemoveGlobalComponent("HELLO_WORLD")
	require.NotNil(s.T(), err)
	require.True(s.T(), strings.HasPrefix(err.Error(), "NOT_FOUND"))
}

func TestUmonoLangTestSuite(t *testing.T) {
	suite.Run(t, new(umonoLangTestSuite))
}
