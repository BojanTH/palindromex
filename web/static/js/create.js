"use strict";

import { toJSONString } from "./form-handler";

document.addEventListener("DOMContentLoaded", function () {
    const messagesForm  = document.querySelector("#_message-form");
    const url = messagesForm.action
    const redirectURL = messagesForm.dataset.redirect

    messagesForm.addEventListener("submit", function(event) {
        event.preventDefault();

        let data = toJSONString(event.target)
        fetch(url, {
            method: "POST",
            body: data
        }).then(response => {
            if (response.status !== 201) {
                console.log(response);

                return;
            }

            window.location.href = redirectURL;
        })
    })
})
