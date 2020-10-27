"use strict";

document.addEventListener("DOMContentLoaded", function () {
    const messagesForm  = document.querySelector("#_message-form");
    const url = messagesForm.action

    fetch(url)
      .then(response => response.json())
      .then(data => renderMessage(data));
})

function renderMessage(data) {
    const messageContent = document.querySelector("#_message-content");
    const messageResult = document.querySelector("#_message-result");
    
    messageContent.innerHTML = data.Content
    messageResult.classList.add("fa")
    if (data.palindrome) {
        messageResult.classList.add("success")
        messageResult.classList.add("fa-check-circle")
    } else {
        messageResult.classList.add("error")
        messageResult.classList.add("fa-times-circle")
    }
}