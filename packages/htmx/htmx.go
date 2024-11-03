package htmx

func SuccessRegister() string {
	return `
	<h5 style="color:green;">Succesfull registration. You can now login </h5>
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

func ErrorRegister() string {
	return `
		<h5 style="color:red;">Registration failed. Please try again or contact thw admin uwu.</h5>
		<br />
		`

}
func UnauthorizedRegister() string {
	return `
		<h5 style="color:red;">Registration failed, you arent an authorized user. Please contact the admin if you are actually authprized</h5>
		<br />
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

            <div id="Badresponse" style="color: red;"></div>

            <br />
            <button
                class="register-button"
                hx-get="/register"
                hx-target="#content"
                hx-swap="outerHTML"
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
        hx-target="#Badresponse"
        hx-swap="outerHTML"
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
        hx-swap="outerHTML"
            >
        Already have an account? Log In
    </button>
<script>

function isStringLengthValid(password) {
    const encoder = new TextEncoder();
    const byteArray = encoder.encode(password);
    // Check if the length of the byte array is less than or equal to 72 because bcrypt doesn't support more than that
    return byteArray.length <= 72;
}

function containsSemicolon(input) {
    return input.includes(';');
}

function checkPassword() {
    let password1 = document.getElementById("password1").value;
    let password2 = document.getElementById("password2").value;
    let username = document.getElementById("username").value; // Assuming you have a username field
    let errorMessage = document.getElementById("Badresponse");

    // Check for semicolons in the username and passwords
    if (containsSemicolon(username) || containsSemicolon(password1)) {
        errorMessage.innerHTML = "Username and password cannot contain a semicolon (;).";
        alert("Username and password cannot contain a semicolon (;).");
        return false;
    }

    if (password1 !== password2) {
        errorMessage.innerHTML = "Passwords do not match. Please make them match.";
        alert("Passwords do not match");
        return false;
    }

    if (!isStringLengthValid(password1)) {
        errorMessage.innerHTML = "Password must be 72 bytes or less.";
        alert("Password must be 72 bytes or less.");
        return false;
    }

    return true; // Return true if all checks pass
}

</script>
</div>`
}

func GetSubmissionSuccess() string {
	return `<h3>Thanks for your submission. it is now <span id="time">uwu</span> in my timezone, so i will see when i can get back at you!</h3>
			<script>	function sleep(ms) {
			return new Promise(resolve => setTimeout(resolve, ms));
		}

		async function time() {
		timeNow = document.getElementById("time");

		const options = {
		timeZone: 'Europe/Zurich',
		 dateStyle: 'full',
		 timeStyle: 'long',
		 /*hour: '2-digit',
		minute: '2-digit',
		second: '2-digit',
		hour12: false*/
		};

		const formatter = new Intl.DateTimeFormat('en-US', options);
		while (true) {

			let date = new Date();
			let formattedDate = formatter.format(date);
			timeNow.innerHTML = formattedDate;
			await sleep(1000);
		}
}

		time()
		</script>
			`
}
