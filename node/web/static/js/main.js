const sender = document.getElementById("sender")

const receiver = document.getElementById("receiver")

const text_area = document.getElementById("sender-area")
const area = document.getElementById("area")

sender.addEventListener("submit", handle)

area.scrollTop = area.scrollHeight

function submitOnEnter(event) {
    if (event.which === 13) {
        if (!event.repeat) {
            const newEvent = new Event("submit", {cancelable: true});
            event.target.form.dispatchEvent(newEvent);
        }
        event.preventDefault(); // Prevents the addition of a new line in the text field
    }
}

sender.addEventListener("keydown", submitOnEnter)

async function handle() {
    event.preventDefault()
    if (text_area.value.trim() === "") {
        return
    }
    const form = event.currentTarget
    const url = "/send?receiver="+receiver.innerHTML.trim()
    const formData = new FormData(form)
    text_area.innerHTML = ""

    const responseData = await send({ url, formData });
    window.location.href="/?current="+receiver.innerHTML
    if (responseData.data !== "ok") {
        console.log(responseData.data)
        console.log(responseData.error)
    }
}

async function send({ url, formData }) {
    const plainFormData = Object.fromEntries(formData.entries());
    const formDataJsonString = JSON.stringify(plainFormData);
    if (plainFormData.toString() === "") {
        return
    }
    const fetchOptions = {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json"
        },
        body: formDataJsonString
    };
    const response = await fetch(url, fetchOptions);

    return response.json();
}

function handleFiles(files) {
    ([...files]).forEach(uploadFile)
}

function uploadFile(file) {
    const url = "/file?receiver="+receiver.innerHTML.trim()
    const formData = new FormData()
    formData.append('file', file)
    fetch(url, {
        method: 'POST',
        body: formData
    }).then((response) => {
        return response.json()
    }).then(async (response) => {
        window.location.href="/?current="+receiver.innerHTML
        if (response.data !== "ok") {
            console.log(response.data)
            console.log(response.error)
        }
    })
}
