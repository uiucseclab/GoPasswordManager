package main

import (
	"bytes"
	"net/http"
	"path"
	"strings"

	"goji.io/pattern"

	"golang.org/x/net/context"
)

func apiPassName(s string) string {
	return strings.TrimSuffix(path.Base(s), ".gpg")
}

/*
GET /api/pass/* - get a password or a list of passwords
Reponse for files:
{
	"name": "base name of the file, minus the .gpg",
	"path": "full/path/to/file",
	"contents": "full file contents, base64 encoded",
	"recipients": ["key","ids","that","can","access"]
}

Reponse for directories:
{
	"children": [
		{
			"name": "name of the child",
			"path": "full path of the child",
			"type": "'dir' or 'file'"
		}
	],
	"recipients": ["key","ids","that","can","access","directory"]
}
*/
func handleGetPass(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	type responseFile struct {
		Name       string   `json:"name"`
		Path       string   `json:"path"`
		Contents   []byte   `json:"contents"`
		Recipients []string `json:"recipients"`
	}
	type responseDirEnt struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Type string `json:"type"`
	}
	type responseDir struct {
		Children   []responseDirEnt `json:"children"`
		Recipients []string         `json:"recipients"`
	}

	p := pattern.Path(ctx)
	ps := PassFromContext(ctx)
	var response interface{}
	if tx, err := ps.Begin(); err != nil {
		rlog(ctx, "Could not start transaction: ", err)
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	} else if exists, isFile := tx.Type(p); !exists {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	} else if isFile {
		if contents, err := tx.Get(p); err != nil {
			rlog(ctx, "Could not get file contents: ", err)
			http.Error(rw, "internal server error", http.StatusInternalServerError)
			return
		} else if recipients, err := getRecipients(bytes.NewReader(contents)); err != nil {
			rlog(ctx, "Could not get recipients: ", err)
			http.Error(rw, "internal server error", http.StatusInternalServerError)
			return
		} else {
			response = responseFile{
				Name:       apiPassName(p),
				Path:       path.Clean(p),
				Contents:   contents,
				Recipients: recipients,
			}
		}
	} else {
		if recipients, err := tx.Recipients(p); err != nil {
			rlog(ctx, "Could not get recipients: ", err)
			http.Error(rw, "internal server error", http.StatusInternalServerError)
			return
		} else if children, err := tx.List(p); err != nil {
			rlog(ctx, "Could not get directory listing: ", err)
			http.Error(rw, "internal server error", http.StatusInternalServerError)
			return
		} else {
			rChildren := make([]responseDirEnt, 0, len(children))
			for _, c := range children {
				if c.File && !strings.HasSuffix(c.Name, ".gpg") {
					continue
				}
				var ch responseDirEnt
				ch.Name = apiPassName(c.Name)
				ch.Path = path.Join(p, c.Name)
				ch.Type = "file"
				if !c.File {
					ch.Type = "dir"
				}
				rChildren = append(rChildren, ch)
			}

			response = responseDir{
				Children:   rChildren,
				Recipients: recipients,
			}
		}
	}

	if err := RenderFromContext(ctx).JSON(rw, http.StatusOK, response); err != nil {
		rlog(ctx, "Could not render JSON: ", err)
	}
}

func handlePostPass(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	// p := pattern.Path(ctx)
	// ps := PassFromContext(ctx)
	// u := UserFromContext(ctx)
	// tx := ps.Begin()
}

func handleDeletePass(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	// p := pattern.Path(ctx)
	// ps := PassFromContext(ctx)
	// u := UserFromContext(ctx)
	// tx := ps.Begin()
}

func handleGetPerm(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	// p := pattern.Path(ctx)
	// ps := PassFromContext(ctx)
	// u := UserFromContext(ctx)
	// tx := ps.Begin()
}

func handlePostPerm(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	// p := pattern.Path(ctx)
	// ps := PassFromContext(ctx)
	// u := UserFromContext(ctx)
	// tx := ps.Begin()
}
