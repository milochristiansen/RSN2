<template>
<div>
	<Form :submit="submit">
		<div>
			<Field type="text" placeholder="Email" v-model="user" name="email" :rules="notempty"/>
			<ErrorMessage class="error" name="email" />
		</div>
		<div>
			<Field type="password" placeholder="Password" v-model="password" name="password" :rules="notempty"/>
			<ErrorMessage class="error" name="password" />
		</div>
		<div>
			<input type="submit" value="Login">
			<span v-if="loginfail" class="error">Invalid username or password.</span>
		</div>
	</Form>
	<section name="nav">
		<router-link to="/forgotpass">Forgot your password?</router-link> |
		<router-link to="/newuser">Sign Up</router-link>
	</section>
</div>
</template>

<script>
import { Form, Field, ErrorMessage } from "vee-validate";

export default {
	name: 'Login',
	components: {
		Form,
		Field,
		ErrorMessage
	},
	data() {
		return {
			user: "",
			password: "",
			loginfail: false,
		}
	},
	methods: {
		notempty(value) {
			if (value == "") {
				return "Must provide a value."
			}
			return true
		},
		submit() {
			let self = this;
			fetch("/api/user/login", {
				method: "POST",
				body: JSON.stringify({
					Email: String(this.user),
					Password: String(this.password),
				})
			})
				.then(function(res) {
					if (res.ok) {
						self.$router.push("/user/unread");
						return
					}
					throw new Error(res.status);
				})
				.catch(error => {
					console.error(error.message)
					self.loginfail = true
				});
		}
	},
	created() {
		let self = this;
		fetch("/api/user/logged-in")
			.then(function(res) {
				if (res.ok) {
					self.$router.push("/user/unread");
					return
				}
				throw new Error(res.status);
			})
			.catch(error => {
				console.log("User is not logged in.")
			});
	}
}
</script>

<style scoped lang="scss">
form {
	display: flex;
	flex-direction: column;

	div {
		position: relative;
		display: flex;
		flex-direction: column;
		justify-content: center;
		
		& > span {
			text-align: center;
			background-color: var(--bg-color);
			font-size: .8em;
			margin-bottom: 5px;
		}
	}

	input {
		max-width: 400px;
		width: 90%;
		padding: 1px 2px;
		align-self: center;

		border-radius: 5px;
		border-style: outset;
		border-width: 3px;
		border-color: var(--secondary-color);

		color: var(--font-color);
		background-color: var(--bg-color);

		&[type=submit] {
			cursor: pointer;
		}
		&[type=text] {
			border-style: inset;
		}
	}
}

section[name=nav] {
	padding-top: .4em;
	font-size: .6em;
	display: flex;
	flex-direction: row;
	justify-content: center;

	a {
		color: var(--font-color);
	}
}

.error {
	color: red;
}
</style>
