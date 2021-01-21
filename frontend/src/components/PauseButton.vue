<template>
	<a :href="href" :style="css" :class="ispaused"></a>
</template>

<script>
export default {
	name: 'PauseButton',
	
	props: {
		size: {
			type: [String],
			default: "25px",
		},
		href: {
			type: [String],
			default: "#",
		},
		color: {
			type: [String],
			default: "black",
		},
		paused: {
			type: [Boolean],
			default: false,
		}
	},

	computed: {
		css() {
			return {
				"--size": this.size,
				"--color": this.color
			}
		},
		ispaused() {
			return this.paused ? "" : "paused"
		}
	}
}
</script>

<style scoped lang="scss">
	a {
		border: 0;
		background: transparent;
		box-sizing: border-box;
		width: 0;
		height: var(--size);
	
		border-color: transparent transparent transparent var(--color);
		transition: 100ms all ease;
		cursor: pointer;
	
		// play state
		border-style: solid;
		border-width: calc(var(--size) / 2) 0 calc(var(--size) / 2) var(--size);
	
		&.paused {
			border-style: double;
			border-width: 0px 0 0px var(--size);
		}
	
		&:hover {
			border-color: transparent transparent transparent var(--color);
		}
	}
</style>
