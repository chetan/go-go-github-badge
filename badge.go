package go_go_github_badge

import (
	"fmt"
	"html/template"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func Generate(username string) (string, error) {
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	if c == nil {
		return "", fmt.Errorf("error creating context")
	}
	r.SetHTMLTemplate(template.Must(template.New("badge.gohtml").Parse(badgeTemplate)))

	c.Params = append(c.Params, gin.Param{Key: "username", Value: username})
	generateBadge(c)

	if w.Code != 200 {
		return "", fmt.Errorf("error generating badge")
	}

	return w.Body.String(), nil
}
