const form = document.getElementById('form')

form.addEventListener('submit', (e) => {
    e.preventDefault()
    const formData = new FormData(e.target)
    const json = JSON.stringify(Object.fromEntries(formData));
    fetch("/login", {
        method: 'POST',
        body: json
    })
        .then((response) => {
            return response.json()
        })
        .then(async (response) => {
            let data = await response
            if (data["error"] === "") {
                statusHidden(true)
                statusColor("green")
                statusText("Успешно")
                window.location.href = "/"
            } else {
                statusHidden(true)
                statusColor("red")
                statusText(data["error"])
            }
        })
})