<template>
	<router-view @theme="swaptheme()" :theme="theme"/>
</template>

<style>
	:root {
	    --primary-color: #302AE6;
	    --secondary-color: #536390;
	    --font-color: #424242;
	    --bg-color: #fff;
	    --heading-color: #292922;
	}
	
	:root[data-theme="dark"] {
	    --primary-color: #9A97F3;
	    --secondary-color: #818cab;
	    --font-color: #e1e1ff;
	    --bg-color: #161625;
	    --heading-color: #818cab;
	}
</style>

<style lang="scss">
html {
	height: 100%;
}
body {
	background-color: var(--bg-color);

	max-width: 800px;
	height: 100%;
	
	margin: 0 auto !important;
	float: none !important;
}

#app {
	height: 100%;
	display: flex;
	flex-direction: column;
}
</style>

<script>
export default {
	name: 'App',

	computed: {
		theme: {
			get() {
				return this.themedat;
			},
			set(v) {
				this.themedat = v
				document.documentElement.setAttribute('data-theme', v);
				localStorage.setItem('theme', v);
			}
		}
	},

	data() {
		return {
			themedat: localStorage.getItem('theme') ? localStorage.getItem('theme') : 'dark',
		}
	},

	methods: {
		swaptheme() {
			this.theme == "light" ? this.theme = "dark" : this.theme = "light"
		}
	},

	mounted() {
		document.documentElement.setAttribute('data-theme', this.theme);
	}
}
</script>
