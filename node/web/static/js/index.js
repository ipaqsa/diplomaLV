const status = document.getElementById("status")

function statusText(text) {
    status.innerHTML = text
}

function statusColor(color) {
    status.style.color = color
}

function statusHidden(s) {
    if (s === true) {
        status.style.display = "block"
    } else {
        status.style.display = "none"
    }
}