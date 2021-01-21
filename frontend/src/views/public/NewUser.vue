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
				<Field type="password" placeholder="Reenter Password" name="password2" :rules="match1"/>
				<ErrorMessage class="error" name="password2" />
			</div>
			<input type="submit" value="Create Account">
		</Form>
</div>
</template>

<script>
import { Form, Field, ErrorMessage } from "vee-validate";

export default {
	name: 'NewUser',
	components: {
		Form,
		Field,
		ErrorMessage
	},
	data() {
		return {
			user: "",
			password: ""
		}
	},
	methods: {
		notempty(value) {
			if (value == "") {
				return "Must provide a value."
			}
			return true
		},
		match1(value) {
			if (value != this.password) {
				return "Passwords do not match."
			}
			return true
		},
		submit() {
			let self = this;
			fetch("/api/user/new", {
				method: "POST",
				body: JSON.stringify({
					Email: String(this.user),
					Password: String(this.password),
				})
			})
				.then(function(res) {
					if (res.ok) {
						self.$router.push("/confirm");
						return
					}
					throw new Error(res.status);
				})
				.catch(error => console.error(error.message));
		}
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

.error {
	color: red;
}
</style>
