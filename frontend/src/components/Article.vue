<template>
	<a :href="data.URL" class="article" @click.prevent="open">
		<span class="row">{{ data.Title}}</span>
		<span class="row">{{ data.URL }}</span>
		<ReadButton class="read" :href="'/mark-read'" @click.stop.prevent="markread" :read="data.Read"/>
	</a>
</template>

<script>
import ReadButton from "@/components/ReadButton.vue";

export default {
	name: 'Article',

	components: {
		ReadButton
	},
	
	emits: {
		changed: null,
	},
	
	props: {
		data: {
			type: [Object],
			default: function() {
				return {
					ID: "",
					Title: "",
					URL: "",
					Published: new Date(),
					Read: false,
				}
			},
			validator: function(value) {
				if (typeof value.ID != "string") {
					return false
				}
				if (typeof value.Title != "string") {
					return false
				}
				if (typeof value.URL != "string") {
					return false
				}
				if (typeof value.Read != "boolean") {
					return false
				}
				return true
			}
		}
	},
	
	data() {
		return {
			url: "",
			name: "",
			addstate: null,
			socket: null,
			list: []
		}
	},
	
	methods: {
		open() {
			window.open(this.data.URL, '_blank');
		},
		markread() {
			if (this.data.Read) {
				fetch("/api/article/unread?id="+this.data.ID).then(() => this.$emit("changed"))
				return
			}
			fetch("/api/article/read?id="+this.data.ID).then(() => this.$emit("changed"))
		}
	}
}
</script>

<style scoped lang="scss">
.article {
	display: flex;
	flex-direction: column;
	position: relative;

	margin: 2px;

	border-radius: 5px;
	border-style: outset;
	border-color: var(--secondary-color);

	color: var(--font-color);
	background-color: var(--bg-color);

	text-decoration: none;
	
	.row {
		color: var(--font-color);
		width: 100%;
	}

	.read {
		position: absolute;
		top: 6px;
		left: calc(100% - 27px);
		background-color: var(--bg-color); /* So it overlays text, not sure if I like this. */
	}
}
</style>
