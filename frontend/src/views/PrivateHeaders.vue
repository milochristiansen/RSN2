<template>
	<header>
		<router-link to="/user/settings">⚙</router-link>
		<router-link to="/user/unread">Unread</router-link>
		<router-link to="/user/feeds">Feeds</router-link>
		<router-link to="/user/logout">Logout</router-link>
	</header>
	<router-view></router-view>
	<footer>
		<span>To report issues, get help, etc. send a message to: </span>
		<a href="&#109;&#97;&#105;&#108;&#116;&#111;&#58;&#109;&#105;&#108;&#111;&#64;&#104;&#116;&#116;&#112;&#115;&#99;&#111;&#108;&#111;&#110;&#115;&#108;&#97;&#115;&#104;&#115;&#108;&#97;&#115;&#104;&#119;&#119;&#119;&#46;&#99;&#111;&#109;">&#109;&#105;&#108;&#111;&#64;&#104;&#116;&#116;&#112;&#115;&#99;&#111;&#108;&#111;&#110;&#115;&#108;&#97;&#115;&#104;&#115;&#108;&#97;&#115;&#104;&#119;&#119;&#119;&#46;&#99;&#111;&#109;</a>
		<a href="dark-theme" class="theme-swap" v-if="theme == 'light'" @click.prevent="swaptheme()">☀</a>
		<a href="light-theme" class="theme-swap" v-else @click.prevent="swaptheme()">☼</a>
	</footer>
</template>

<script>
export default {
	name: 'PrivateHeaders',

	emits: ['theme'],

	props: {
		theme: String
	},

	methods: {
		swaptheme() {
			this.$emit('theme')
		}
	}
}
</script>

<style scoped lang="scss">
	header {
		margin-bottom: 5px;
		color: var(--bg-color);
		background-color: var(--heading-color);
		
		a {
			margin-left: 10px;
			color: var(--bg-color);
		}
	}
	footer {
		text-align: left;
		margin-top: 10px;
		font-size: .5em;
		color: var(--font-color);

		margin-top: auto;
		margin-bottom: 10px;
		
		a {
			color: var(--font-color);;
		}
		.theme-swap {
			float: right;
		}
	}
</style>
