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

function search() {
    let input, filter, ul, li, a, i, txtValue;
    input = document.getElementById('searcher');
    filter = input.value.toUpperCase();
    ul = document.getElementById("contacts");
    li = ul.getElementsByTagName('li');

    for (i = 0; i < li.length; i++) {
        a = li[i].getElementsByTagName("a")[0];
        txtValue = a.textContent || a.innerText;
        if (txtValue.toUpperCase().indexOf(filter) > -1) {
            li[i].style.display = "";
        } else {
            li[i].style.display = "none";
        }
    }
}

function downloadFile() {
    let receivern = receiver.innerHTML.trim()
    let filename = event.currentTarget.getElementsByTagName("span")[0].innerHTML.trim()
    if ((filename === "")||(receivern === "")) {
        return
    }

    const url = "/download?receiver="+receivern+"&filename="+encodeURIComponent(filename)
    fetch(url, {
        method: 'GET',
    }).then((response) => {
        return response.blob()
    }).then(async (response) => {
        window.location.href="/?current="+receivern
        let file = window.URL.createObjectURL(response)
        window.location.assign(file)
    })
}