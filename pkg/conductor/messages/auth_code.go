package messages

func AuthCode(code string) string {
	return `<div>
				<h1>Код : ` + code + `</h1>
				<p>Используйте этот код для авторизации на сайте.</p>
			</div>`
}
