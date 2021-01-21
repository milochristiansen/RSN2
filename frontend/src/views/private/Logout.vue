<template>
<p v-if="loggedout == null">Logging you out...</p>
<p v-else-if="loggedout == true">Thank you for testing RSN2! Please remember to report any bugs. Feature requests are also welcome. <router-link to="/">Return to login page.</router-link></p>
<p v-else>There was an error logging you out.</p>
</template>

<script>
import { Form, Field, ErrorMessage } from "vee-validate";

export default {
	name: 'Logout',
	data() {
		return {
			loggedout: null,
		}
	},
	created() {
		let self = this
		fetch("/api/user/logout")
			.then(function(res) {
				if (res.ok) {
					self.loggedout = true
					return
				}
				self.loggedout = false
				throw new Error(res.status);
			})
			.catch(error => console.error(error.message));
	}
}
</script>

<style scoped lang="scss">
p {
	color: var(--font-color);
	a {
		color: var(--secondary-color);
	}
}
</style>
