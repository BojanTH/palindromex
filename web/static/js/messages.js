"use strict";

document.addEventListener("DOMContentLoaded", function () {
    const messagesDiv  = document.querySelector("#_messages");
    const url = messagesDiv.dataset.url

    let xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState !== 4) {
            return;
        }

        let data = this.responseText
        if (!data) {
            return;
        }
        let messages = JSON.parse(data)
        if (!messages || !messages.length) {
            return;
        }

        messages.forEach(message => {
            console.log(message);
            renderElement(messagesDiv, message);
        });
    };

    xhr.open("GET", url);
    xhr.send()
});

function renderElement(destination, message) {
    let element = "<div class='row'>";
    if (message.palindrome) {
        element += `<div class='col-1 center'><i class='palindrome-result fa fa-check-circle'></i></div>`
    } else {
        element += "<i class='palindrome-result fa fa-times-circle'></i>"
    }
    element += `<div class="col-11 menu-item p-2 mt-4">${message.Content}</div>`;
    element += "</div>";

    destination.innerHTML = destination.innerHTML + element;
}