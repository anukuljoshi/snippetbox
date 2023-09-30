package main

import "snippetbox.anukuljoshi/internals/models"

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}
