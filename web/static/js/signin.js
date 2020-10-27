"use strict";

import { toJSONString } from "./form-handler";

document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("_signin-form");
    const url = form.action
    const email = document.getElementById("_email");
    const password = document.getElementById("_password");

    form.addEventListener("submit", function (event) {
        event.preventDefault();
        if (!email.value || !password.value) {
            console.log("Invalid form");

            return;
        }

        let data = toJSONString(event.target)
        fetch(url, {
            method: "POST",
            body: data
        }).then(response => {
            if (response.status !== 200) {
                console.log(response);
            }

            return response.json()
        }).then(response => {
            if (response) {
                window.location.href = response["url"]
            }
        });
    });
});
