<template>
	<a :href="data.URL" class="article" @click.prevent="open">
		<span class="row">{{ data.Title}}</span>
		<span class="row">{{ data.FeedName }}</span>
		<CloseButton class="read" :href="'/mark-read'" @click.stop.prevent="markread" :size="'15px'" :color="'var(--secondary-color)'"/>
	</a>
</template>

<script>
import CloseButton from "@/components/CloseButton.vue";

export default {
	name: 'UnreadArticle',

	components: {
		CloseButton
	},
	
	props: {
		data: {
			type: [Object],
			default: function() {
				return {
					ID: "",
					Title: "",
					URL: "",
					FeedName: "",
					Published: new Date(),
				}
			},
			validator: function(value) {
				if (!value.ID instanceof String) {
					return false
				}
				if (!value.Title instanceof String) {
					return false
				}
				if (!value.URL instanceof String) {
					return false
				}
				if (!value.FeedName instanceof String) {
					return false
				}
				if (value.Published instanceof Date) {
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
			fetch("/api/article/read?id="+this.data.ID)
		},
		markread() {
			fetch("/api/article/read?id="+this.data.ID)
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
		top: 3px;
		left: calc(100% - 18px);
		background-color: var(--bg-color); /* So it overlays text, not sure if I like this. */
	}
}
</style>
