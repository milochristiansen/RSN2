<template>
	<div>
		<p v-if="confirmed == null">Checking your token...</p>
		<p v-else-if="confirmed == true">Your unconfirmed account has been deleted from our database. Maybe you would like to <router-link to="/">make another one?</router-link></p>
		<p v-else>There was an error deleting your unconfirmed account. Please ensure you used the correct link and that your account is actualy unconfirmed.</p>
	</div>
</template>

<script>
export default {
	name: 'DeleteEmail',
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
		fetch("/api/user/delete-email?token="+this.token)
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
