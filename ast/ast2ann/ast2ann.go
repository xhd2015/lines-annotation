package ast2ann

import (
	"go/token"
	"path/filepath"
	"strings"

	"go/ast"

	"github.com/xhd2015/go-coverage/cover"
	ast_load "github.com/xhd2015/lines-annotation/ast"
	"github.com/xhd2015/lines-annotation/model"
	"github.com/xhd2015/lines-annotation/model/coverage"
)

func CollectBlockProfilesNoModPath(astInfo ast_load.LoadInfo) (coverage.BinaryProfile, error) {
	return CollectBlockProfiles("", astInfo)
}

func CollectBlockProfiles(modPath string, astInfo ast_load.LoadInfo) (coverage.BinaryProfile, error) {
	profile := make(coverage.BinaryProfile)
	fset := astInfo.FileSet()
	astInfo.RangeFiles(func(f ast_load.File) bool {
		if f.SyntaxError() != nil {
			return true
		}
		// pkg file
		pkgFile := slashlizePath(f.RelPath())
		if modPath != "" {
			pkgFile = modPath + "/" + pkgFile
		}

		cover.RangeBlocks(fset, f.Ast(), f.Content(), func(insertPos, pos, end token.Pos, numStmts int, basicStmts []ast.Stmt) {
			profile[pkgFile] = append(profile[pkgFile], &coverage.BlockStats{
				Block: GetBlock(fset, pos, end),
			})
		})
		return true
	})
	profile.SortAll()
	return profile, nil
}

func slashlizePath(path string) string {
	if filepath.Separator == '/' {
		return path
	}
	return strings.ReplaceAll(path, string(filepath.Separator), "/")
}

func GetBlock(fset *token.FileSet, pos, end token.Pos) *model.Block {
	posBegin := fset.Position(pos)
	posEnd := fset.Position(end)
	return &model.Block{
		StartLine: posBegin.Line,
		StartCol:  posBegin.Column,
		EndLine:   posEnd.Line,
		EndCol:    posEnd.Column,
	}
}
