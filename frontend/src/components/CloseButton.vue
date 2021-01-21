<template>
	<a :href="href" class="root" :style="css">
		<div class="right">
			<div class="left"></div>
		</div>
	</a>
</template>

<script>
export default {
	name: 'CloseButton',
	
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
		}
	},

	computed: {
		css() {
			return {
				"--size": this.size,
				"--color": this.color
			}
		}
	}
}
</script>

<style scoped lang="scss">
	.root {
	  width: var(--size);
	  height: var(--size);

		.right {
			height: var(--size);
			width: 2px;
			margin-left: calc(var(--size) / 2);
			background-color: var(--color);
			transform: rotate(45deg);
			z-index: 1;

			.left {
				height: var(--size);
				width: 2px;
				background-color: var(--color);
				transform: rotate(90deg);
				z-index: 2;
			}
		}
	}
</style>
