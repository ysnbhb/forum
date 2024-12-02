function showPss(idButtom, idPassword) {
    let hiden = true;
    const button = document.getElementById(idButtom);
    button.addEventListener("click", (event) => {
        event.preventDefault();
        const password = document.getElementById(idPassword);
        if (hiden) {
            password.setAttribute("type", "text");
            button.innerHTML = "visibility_off";
        } else {
            password.setAttribute("type", "password");
            button.innerHTML = "visibility";
        }
        hiden = !hiden;
    });
}

function confermetPas() {
    const password = document.getElementById("password");
    const confpassword = document.getElementById("confpassword");

    if (password && confpassword) {
        confpassword.addEventListener("input", () => {
            if (password.value !== confpassword.value) {
                confpassword.style.border = "1px solid red";
                return true
            } else {
                confpassword.style.border = "1px solid green";
                return false
            }
        });
    } else {
        console.error("Password or confirm password input not found!");
    }
}


function runButtom(idButtom) {
    const buttom = document.getElementById(idButtom)
    if (confermetPas()) {
        buttom.addEventListener("click" , (env)=> {
            env.defaultPrevented()
        })
    }
}
showPss("changTYpe", "password")
showPss("changContype", "confpassword")

runButtom("sing-up")
