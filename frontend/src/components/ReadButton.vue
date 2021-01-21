<template>
	<a :href="href" :style="css" :class="isread"><div></div></a>
</template>

<script>
export default {
	name: 'ReadButton',
	
	props: {
		href: {
			type: [String],
			default: "#",
		},
		read: {
			type: [Boolean],
			default: false,
		}
	},

	computed: {
		css() {
			return {
				"--color": this.read ? "var(--primary-color)" : "var(--font-color)"
			}
		},
		isread() {
			return this.read ? "" : "read"
		}
	}
}
</script>

<style scoped lang="scss">
	a {
		height: 24px;
		width: 24px;
		position: relative;
		div {
			position: absolute;
			left: 5px;
			
			display: inline-block;
			transform: rotate(45deg);
			height: 16px;
			width: 8px;
			border-bottom: 4px solid var(--color);
			border-right: 4px solid var(--color);	
		}
	}
</style>
