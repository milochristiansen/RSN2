<template>
	<div>
		<p v-if="confirmed == null">Checking your confirmation token...</p>
		<p v-else-if="confirmed == true">Thank you for confirming your email! <router-link to="/">Now go login!</router-link></p>
		<p v-else>There was an error confirming your email. Please ensure you used the correct link.</p>
	</div>
</template>

<script>
export default {
	name: 'EmailConfirmed',
	props: {
		token: {
			type: [String],
			default: "invalid"
		}
	},
	data() {
		return {
			confirmed: null,
		}
	},
	created() {
		let self = this
		fetch("/api/user/confirm-email?token="+this.token)
			.then(function(res) {
				if (res.ok) {
					self.confirmed = true
					return
				}
				self.confirmed = false
				throw new Error(res.status);
			})
			.catch(error => console.error(error.message));
	}
}
</script>

<style lang="scss">
	p {
		color: var(--font-color);
	}
	a {
		color: var(--secondary-color);
	}
</style>
