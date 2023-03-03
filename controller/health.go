package controller

import (
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(
		`
		<html>
			<body style="
				min-height: 100vh;
				display: flex;
				justify-content: center;
				align-items: center;
			">
				<h2>
					pong
				</h2>
			</body>
		</html>
		`,
	))
}
