"use strict";

document.addEventListener("DOMContentLoaded", function () {
    const messagesDiv  = document.querySelector("#_messages");
    const url = messagesDiv.dataset.url

    fetch(url)
      .then(response => response.json())
      .then(messages => messages.forEach(message => {
          renderElement(messagesDiv, message)
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
            <a class="btn btn-primary" href="/users/${message.UserID}/edit-message/${message.id}">Edit</a>
        </div>
    </div>`;

    destination.innerHTML = destination.innerHTML + element;
}