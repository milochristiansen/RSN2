<template>
<Form :submit="changepassword">
	<div>
		<Field type="password" placeholder="New Password" v-model="password" name="password" :rules="notempty"/>
		<ErrorMessage class="error" name="password" />
	</div>
	<div>
		<Field type="password" placeholder="Reenter New Password" name="password2" :rules="match1"/>
		<ErrorMessage class="error" name="password2" />
	</div>
	<div>
		<Field type="password" placeholder="Existing Password" v-model="oldpassword" name="oldpassword" :rules="notempty"/>
		<ErrorMessage class="error" name="oldpassword" />
	</div>
	<input type="submit" value="Change Password">
	<span v-if="passstate === false" class="error">Failed changing password.</span>
	<span v-else-if="passstate === true">Password changed!</span>
</Form>
<hr>
<Form :submit="changename">
	<p>
		WARNING! If you change your username you will not be able to login again until you confirm your email.
		Making a typo here could lock you out of your account.
	</p>
	<div>
		<Field type="text" placeholder="Email" v-model="user" name="email" :rules="notempty"/>
		<ErrorMessage class="error" name="email" />
	</div>
	<div>
		<Field type="password" placeholder="Password" v-model="oldpassword2" name="oldpassword2" :rules="notempty"/>
		<ErrorMessage class="error" name="oldpassword2" />
	</div>
	<input type="submit" value="Change Username">
	<span v-if="userstate === false" class="error">Failed changing username.</span>
	<span v-else-if="userstate === true">Username changed!</span>
</Form>
</template>

<script>
import { Form, Field, ErrorMessage } from "vee-validate";

export default {
	name: 'UserSettings',
	
	components: {
		Form,
		Field,
		ErrorMessage
	},
	
	data() {
		return {
			inflight: false,
		
			password: "",
			oldpassword: "",
			passstate: null,
			
			user: "",
			oldpassword2: "",
			userstate: null,
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
		changepassword() {
			if (this.inflight) {
				return
			}
			this.inflight = true

			let self = this;
			fetch("/api/user/new-pass", {
				method: "POST",
				body: JSON.stringify({
					Password: String(this.password),
					OldPassword: String(this.oldpassword),
				})
			})
				.then(function(res) {
					self.inflight = false
					if (res.ok) {
						self.password = ""
						self.oldpassword = ""
						
						self.passstate = true
						setTimeout(() => self.passstate = null, 3000)
						return
					}
					throw new Error(res.status);
				})
				.catch(error => {
					console.error(error.message)
					self.passstate = false
					setTimeout(() => self.passstate = null, 5000)
				});
		},
		changename() {
			if (this.inflight) {
				return
			}
			this.inflight = true

			let self = this;
			fetch("/api/user/new-name", {
				method: "POST",
				body: JSON.stringify({
					Email: String(this.user),
					Password: String(this.password2),
				})
			})
				then(function(res) {
					self.inflight = false
					if (res.ok) {
						self.user = ""
						self.oldpassword2 = ""
					
						self.userstate = true
						setTimeout(() => self.userstate = null, 3000)
						return
					}
					throw new Error(res.status);
				})
				.catch(error => {
					console.error(error.message)
					self.userstate = false
					setTimeout(() => self.userstate = null, 5000)
				});
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

	span, p {
		color: var(--font-color)
	}
}

.error {
	color: red;
}

hr {
	width: 90%;
	margin-top: 20px;
	margin-bottom: 20px;
}
</style>
