<template>
	<form @submit.prevent="addfeed">
		<input type="text" placeholder="Feed URL" v-model="url"/>
		<input type="text" placeholder="Feed Name" v-model="name"/>
		<input type="submit" value="Subscribe Feed">
		<span v-if="addstate === false" class="error">Failed adding feed.</span>
		<span v-else-if="addstate === true">Feed added!</span>
	</form>
</template>

<script>
export default {
	name: 'AddFeed',

	emits: {
		added: null,
	},
	
	data() {
		return {
			url: "",
			name: "",
			addstate: null
		}
	},

	methods: {
		addfeed() {
			if (this.url == "" || this.name == "") {
				this.addstate = false
				setTimeout(() => this.addstate = null, 5000)
				return
			}

			let self = this;
			fetch("/api/feed/subscribe", {
				method: "POST",
				body: JSON.stringify({
					URL: String(this.url),
					Name: String(this.name),
				})
			})
				.then(function(res) {
					if (res.ok) {
						self.addstate = true
						setTimeout(() => self.addstate = null, 3000)
						self.$emit("added")
						return
					}
					throw new Error(res.status);
				})
				.catch(error => {
					console.error(error.message)
					self.addstate = false
					setTimeout(() => self.addstate = null, 5000)
				});
		}
	}
}
</script>

<style lang="scss">
	form {
		width: 100%;
		display: flex;
		flex-direction: row;
		padding-top: 5px;

		margin-top: auto;
		margin-bottom: 5px;

		@media (max-width: 500px) {
			flex-direction: column;
		}
		
		input {
			flex: 1 1 auto;

			margin: 0;
			padding-left: 5px;
			padding-right: 5px;

			border-radius: 5px;
			border-style: outset;
			border-width: 3px;
			border-color: var(--secondary-color);

			color: var(--font-color);
			background-color: var(--bg-color);
			
			&[type=submit] {
				cursor: pointer;
				
				@media (min-width: 500px) {
					max-width: 10em;
				}
			}
			&[type=text] {
				border-style: inset;
			}
		}

		.error {
			color: red;
		}

		span {
			margin-left: 5px;
			color: var(--font-color);
		}
	}
</style>
