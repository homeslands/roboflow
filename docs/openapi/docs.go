package openapi

import (
	"embed"
	"fmt"
)

//go:embed build/openapi.yml
var OpenapiSpec embed.FS

// GetTemplate returns the HTML template for Swagger UI
func GetTemplate(specPath string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta name="description" content="SwaggerUI" />
  <title>SwaggerUI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://static1.smartbear.co/swagger/media/assets/swagger_fav.png" sizes="32x32" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '%s',
      dom_id: '#swagger-ui',
      deepLinking: true,
	  showExtensions: true,
	  showCommonExtensions: true,
    });
  };
</script>
</body>
</html>
`, specPath)
}
