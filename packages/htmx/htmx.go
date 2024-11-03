package htmx

func SuccessRegister() string {
	return `
	<h5 style="color:green;">Succesfull registration. You can now login here:</h5>
	<br />
    <button
        class="login-button"
        hx-get="/relogin"
        hx-target="#content"
            >
        Log In
    </button>
	`
}

func ReturnReloginString() string {
	return `
 <div id="content">
            <form
                id="passwordForm"
                hx-post="/submit-password"
                hx-target="#content"
                hx-swap="innerHTML"
                hx-on="htmx:responseError: document.getElementById('Badresponse').innerHTML = event.detail.xhr.responseText; document.getElementById('Badresponse').style.display = 'block';"
            >
            <h2>Login</h2>
                <input
                    type="text"
                    name="username"
                    required
                    placeholder="Enter Username"
                />
                <input
                    type="password"
                    name="password"
                    required
                    placeholder="Enter Password"
                />
                <br />
                <button type="submit">Submit</button>
            </form>
            <div id="Badresponse"></div>
            <br />
            <button
                class="register-button"
                hx-get="/register"
                hx-target="#content"
            >
                Register
            </button>
        </div>
	`
}

func ReturnRegisterString() string {

	return `
	<div id="content">
    <form
        id="registrationForm"
        onsubmit="return checkPassword()"
        hx-post="/submit-registration"
        hx-target="#content"
        hx-swap="innerHTML"
        hx-on="htmx:responseError: document.getElementById('Badresponse').innerHTML = event.detail.xhr.responseText; document.getElementById('Badresponse').style.display = 'block';"
    >
        <h2>Register</h2>
        <input
            type="text"
            name="username"
            required
            placeholder="Enter Username"
        />
        <input
            type="password"
            name="password"
            id="password1"
            required
            placeholder="Enter Password"
        />
        <input
            type="password"
            name="confirm_password"
            id="password2"
            required
            placeholder="Confirm Password"
        />
        <br />
        <button type="submit">Register</button>
    </form>
    <div id="Badresponse" style="color: red;"></div>
    <br />
    <button
        class="login-button"
        hx-get="/relogin"
        hx-target="#content"
            >
        Already have an account? Log In
    </button>
<script>
function isStringLengthValid(password) {
    const encoder = new TextEncoder();
    const byteArray = encoder.encode(password);
    // Check if the length of the byte array is less than or equal to 72 cuz bcrypt doesnt support more than that
    return byteArray.length <= 72;
}
function checkPassword() {

    let password1 = document.getElementById("password1").value;
    let password2 = document.getElementById("password2").value;
    let errorMessage = document.getElementById("Badresponse");

    if (password1 !== password2) {
        errorMessage.innerHTML = "Passwords do not match. Please make them match.";
        alert("Passwords do not match");
        return false;
    }
    return isStringLengthValid(password1);

}
</script>
</div>`
}
