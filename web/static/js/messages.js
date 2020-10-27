"use strict";

document.addEventListener("DOMContentLoaded", function () {
    const messagesDiv  = document.querySelector("#_messages");
    const url = messagesDiv.dataset.url

    fetch(url)
        .then(response => response.json())
        .then(messages => messages.forEach(message => {
            renderElement(messagesDiv, message);
            attachDeleteListeners();
        }));

});

function renderElement(destination, message) {
    let element =
    `<div class="row mt-4">
        <div class="col-1 center">`;

    if (message.palindrome) {
        element += `<i class="palindrome-result mt-2 fa fa-check-circle success"></i>`
    } else {
        element += `<i class="palindrome-result mt-2 fa fa-times-circle error"></i>`
    }

    element +=
        `</div>
        <div class="col-11 menu-item p-2">${message.Content}</div>
    </div>
    <div class="row">
        <div class="ml-auto">
            <a class="_delete btn btn-warning text-light" href="/v1/users/${message.UserID}/messages/${message.id}">Delete</a>
            <a class="btn btn-primary text-light" href="/users/${message.UserID}/edit-message/${message.id}">Edit</a>
        </div>
    </div>`;

    destination.innerHTML = destination.innerHTML + element;
}

function attachDeleteListeners() {
    const deleteButtons  = document.querySelectorAll("._delete");
    deleteButtons.forEach(button => {
        button.addEventListener("click", function(event) {
            event.preventDefault()
            let url = event.target.attributes.href.value;

            fetch(url, {
                method: "DELETE"
            }).then(response => {
                if (response.status !== 204) {
                    console.log(response);

                    return;
                }

                window.location.href = window.location.href;
            })
        })
    })
}