package file

import (
	"fmt"
	"testing"
)

type Case struct {
	path    string
	wantErr error
}

func TestValidateFileName(t *testing.T) {
	cases := []Case{
		// Success cases
		{
			path:    "test.txt",
			wantErr: nil,
		},
		{
			path:    "file-name.doc",
			wantErr: nil,
		},
		{
			path:    "file_name_test.test",
			wantErr: nil,
		},
		{
			path:    "file123",
			wantErr: nil,
		},
		{
			path:    "file-name_123.file",
			wantErr: nil,
		},
		{
			path:    "file.",
			wantErr: nil,
		},
		{
			path:    ".file-name",
			wantErr: nil,
		},
		{
			path:    "file.name.test.txt",
			wantErr: nil,
		},

		// Failed cases
		// Empty filename
		{
			path:    "",
			wantErr: ErrEmptyFileName,
		},

		// Forbidden characters
		{
			path:    "file/name.txt",
			wantErr: ErrForbiddenCharsFileName,
		},
		{
			path:    "test/file/",
			wantErr: ErrForbiddenCharsFileName,
		},
		{
			path:    "test/file",
			wantErr: ErrForbiddenCharsFileName,
		},
		{
			path:    "file\\name",
			wantErr: ErrForbiddenCharsFileName,
		},
		{
			path:    "file-name:test",
			wantErr: ErrForbiddenCharsFileName,
		},
		{
			path:    "file*.txt",
			wantErr: ErrForbiddenCharsFileName,
		},
		{
			path:    "<file>",
			wantErr: ErrForbiddenCharsFileName,
		},
		{
			path:    "file|name.txt",
			wantErr: ErrForbiddenCharsFileName,
		},
		{
			path:    "?<file|name:*>./txt",
			wantErr: ErrForbiddenCharsFileName,
		},

		// Dots only file name
		{
			path:    ".",
			wantErr: ErrDotsFileName,
		},
		{
			path:    "....",
			wantErr: ErrDotsFileName,
		},
		{
			path:    ".......",
			wantErr: ErrDotsFileName,
		},
	}

	for ind, item := range cases {
		caseName := fmt.Sprintf("case %d: [path %s]", ind, item.path)

		err := ValidateFileName(item.path)
		pe, ok := err.(*PathError)
		if ok {
			if pe.Err != item.wantErr {
				t.Fatalf("[%s] get error: %v, expected error: %v", caseName, pe.Err, item.wantErr)
			}
		} else {
			if err != item.wantErr {
				t.Fatalf("[%s] get error: %v, expected error: %v", caseName, err, item.wantErr)
			}
		}
	}
}

func TestValidateDirPath(t *testing.T) {
	cases := []Case{
		// Success cases
		{
			path:    "test/dir/path/",
			wantErr: nil,
		},
		{
			path:    "test/dir/path",
			wantErr: nil,
		},
		{
			path:    "test\\dir\\path",
			wantErr: nil,
		},
		{
			path:    "test.test/dir-name/",
			wantErr: nil,
		},
		{
			path:    "test",
			wantErr: nil,
		},
		{
			path:    "./test",
			wantErr: nil,
		},
		{
			path: "/test/",
			wantErr: nil,
		},

		// Failed cases
		// Empty dirpath
		{
			path:    "",
			wantErr: ErrEmptyDirPath,
		},

		// Forbidden characters
		{
			path:    "path/dir\"/test",
			wantErr: ErrForbiddenCharsDirPath,
		},
		{
			path:    "dir/<test>/",
			wantErr: ErrForbiddenCharsDirPath,
		},
		{
			path:    "dir?/test/*",
			wantErr: ErrForbiddenCharsDirPath,
		},
		{
			path:    "test|dir/",
			wantErr: ErrForbiddenCharsDirPath,
		},

		// Invalid separator using
		{
			path:    "test/dir\\",
			wantErr: ErrSepInvalidSyntax,
		},
		{
			path:    "//\\",
			wantErr: ErrSepInvalidSyntax,
		},

		// Dots only folder name
		{
			path:    ".../../test/dir",
			wantErr: ErrFolderDotsDirPath,
		},
		{
			path:    "../dir/",
			wantErr: ErrFolderDotsDirPath,
		},
		{
			path:    "....",
			wantErr: ErrFolderDotsDirPath,
		},
		{
			path: "test/./dir",
			wantErr: ErrFolderDotsDirPath,
		},
	}
	for ind, item := range cases {
		caseName := fmt.Sprintf("case %d: [path %s]", ind, item.path)

		err := ValidateDirPath(item.path)
		pe, ok := err.(*PathError)
		if ok {
			if pe.Err != item.wantErr {
				t.Fatalf("[%s] get error: %v, expected error: %v", caseName, pe.Err, item.wantErr)
			}
		} else {
			if err != item.wantErr {
				t.Fatalf("[%s] get error: %v, expected error: %v", caseName, err, item.wantErr)
			}
		}
	}
}
