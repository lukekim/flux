// DO NOT EDIT: This file is autogenerated via the builtin command.

package profiler

import (
	ast "github.com/influxdata/flux/ast"
	runtime "github.com/influxdata/flux/runtime"
)

func init() {
	runtime.RegisterPackage(pkgAST)
}

var pkgAST = &ast.Package{
	BaseNode: ast.BaseNode{
		Comments: nil,
		Errors:   nil,
		Loc:      nil,
	},
	Files: []*ast.File{&ast.File{
		BaseNode: ast.BaseNode{
			Comments: nil,
			Errors:   nil,
			Loc: &ast.SourceLocation{
				End: ast.Position{
					Column: 31,
					Line:   17,
				},
				File:   "profiler.flux",
				Source: "package profiler\n\n\n// EnabledProfilers sets a list of profilers that should be enabled during execution.\n//\n// Available profilers are:\n//   * query - Profiles time spent in the various phases of query execution.\n//   * operator - Profiles time spent in each operator of the query.\n//\n// Example:\n//\n//    import \"profiler\"\n//    option profiler.enabledProfilers = [\"query\", \"operator\"]\n//\noption enabledProfilers = [\"\"]",
				Start: ast.Position{
					Column: 1,
					Line:   3,
				},
			},
		},
		Body: []ast.Statement{&ast.OptionStatement{
			Assignment: &ast.VariableAssignment{
				BaseNode: ast.BaseNode{
					Comments: nil,
					Errors:   nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 31,
							Line:   17,
						},
						File:   "profiler.flux",
						Source: "enabledProfilers = [\"\"]",
						Start: ast.Position{
							Column: 8,
							Line:   17,
						},
					},
				},
				ID: &ast.Identifier{
					BaseNode: ast.BaseNode{
						Comments: nil,
						Errors:   nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 24,
								Line:   17,
							},
							File:   "profiler.flux",
							Source: "enabledProfilers",
							Start: ast.Position{
								Column: 8,
								Line:   17,
							},
						},
					},
					Name: "enabledProfilers",
				},
				Init: &ast.ArrayExpression{
					BaseNode: ast.BaseNode{
						Comments: nil,
						Errors:   nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 31,
								Line:   17,
							},
							File:   "profiler.flux",
							Source: "[\"\"]",
							Start: ast.Position{
								Column: 27,
								Line:   17,
							},
						},
					},
					Elements: []ast.Expression{&ast.StringLiteral{
						BaseNode: ast.BaseNode{
							Comments: nil,
							Errors:   nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 30,
									Line:   17,
								},
								File:   "profiler.flux",
								Source: "\"\"",
								Start: ast.Position{
									Column: 28,
									Line:   17,
								},
							},
						},
						Value: "",
					}},
					Lbrack: nil,
					Rbrack: nil,
				},
			},
			BaseNode: ast.BaseNode{
				Comments: []ast.Comment{ast.Comment{Text: "// EnabledProfilers sets a list of profilers that should be enabled during execution.\n"}, ast.Comment{Text: "//\n"}, ast.Comment{Text: "// Available profilers are:\n"}, ast.Comment{Text: "//   * query - Profiles time spent in the various phases of query execution.\n"}, ast.Comment{Text: "//   * operator - Profiles time spent in each operator of the query.\n"}, ast.Comment{Text: "//\n"}, ast.Comment{Text: "// Example:\n"}, ast.Comment{Text: "//\n"}, ast.Comment{Text: "//    import \"profiler\"\n"}, ast.Comment{Text: "//    option profiler.enabledProfilers = [\"query\", \"operator\"]\n"}, ast.Comment{Text: "//\n"}},
				Errors:   nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 31,
						Line:   17,
					},
					File:   "profiler.flux",
					Source: "option enabledProfilers = [\"\"]",
					Start: ast.Position{
						Column: 1,
						Line:   17,
					},
				},
			},
		}},
		Eof:      nil,
		Imports:  nil,
		Metadata: "parser-type=rust",
		Name:     "profiler.flux",
		Package: &ast.PackageClause{
			BaseNode: ast.BaseNode{
				Comments: []ast.Comment{ast.Comment{Text: "// Profiler exposes an API to profile queries.\n"}, ast.Comment{Text: "// Profile results are returned as an extra result in the response named according to the profiles which are enabled.\n"}},
				Errors:   nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 17,
						Line:   3,
					},
					File:   "profiler.flux",
					Source: "package profiler",
					Start: ast.Position{
						Column: 1,
						Line:   3,
					},
				},
			},
			Name: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Comments: nil,
					Errors:   nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 17,
							Line:   3,
						},
						File:   "profiler.flux",
						Source: "profiler",
						Start: ast.Position{
							Column: 9,
							Line:   3,
						},
					},
				},
				Name: "profiler",
			},
		},
	}},
	Package: "profiler",
	Path:    "profiler",
}
