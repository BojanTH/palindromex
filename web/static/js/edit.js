"use strict";
import { toJSONString } from "./form-handler";

document.addEventListener("DOMContentLoaded", function () {
    const messagesForm  = document.querySelector("#_message-form");
    const url = messagesForm.action
    const redirectURL = messagesForm.dataset.redirect

    fetch(url)
      .then(response => response.json())
      .then(data => renderMessage(data));

    messagesForm.addEventListener("submit", function(event) {
        event.preventDefault();

        let data = toJSONString(event.target)
        console.log(data)
        fetch(url, {
            method: "PUT",
            body: data
        }).then(response => {
            if (response.status !== 204) {
                console.log(response);

                return;
            }

            window.location.href = redirectURL;
        })
    })
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