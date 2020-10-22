"use strict";

document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("_signup-form");
    const email = document.getElementById("_email");
    const password = document.getElementById("_password");
    const confirmPassword = document.getElementById("_confirm-password");

    form.addEventListener("submit", function (event) {
        event.preventDefault();
        if (!email.value || password.value !== confirmPassword.value) {
            console.log("The confirm password is not the same as the password");
            return;
        }
        this.submit();
    });
});