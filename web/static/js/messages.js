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
    let element =
    `<div class="row mt-4">
        <div class="col-1 center">`;

    if (message.palindrome) {
        element += `<i class="palindrome-result mt-2 fa fa-check-circle"></i>`
    } else {
        element += `<i class="palindrome-result mt-2 fa fa-times-circle"></i>`
    }

    element +=
        `</div>
        <div class="col-11 menu-item p-2">${message.Content}</div>
    </div>
    <div class="row">
        <div class="ml-auto">
            <a class="btn btn-primary" href="/ui/users/${message.UserID}/edit-message/${message.id}">Edit</a>
        </div>
    </div>`;

    destination.innerHTML = destination.innerHTML + element;
}