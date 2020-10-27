"use strict";

import { toJSONString } from "./form-handler";

document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("_signup-form");
    const url = form.action
    const redirectURL = form.dataset.redirect
    const email = document.getElementById("_email");
    const password = document.getElementById("_password");
    const confirmPassword = document.getElementById("_confirm-password");

    form.addEventListener("submit", function (event) {
        event.preventDefault();
        if (!email.value) {
            console.log("Empty email address");

            return;
        }
        if (password.value !== confirmPassword.value) {
            console.log("The confirm password is not the same as the password");

            return;
        }

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
    });
});