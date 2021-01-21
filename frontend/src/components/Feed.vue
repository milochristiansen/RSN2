<template>
	<a :href="'/user/feed-details?id=' + data.ID" class="feed">
		<span class="row">{{ data.Name}}<span v-if="data.Paused"> (paused)</span></span>
		<span class="row">{{ data.URL }}</span>
		<CloseButton v-if="deleting == false" class="delete" :href="'/delete'" @click.stop.prevent="predelete" :size="'15px'" :color="'var(--secondary-color)'"/>
		<CloseButton v-else class="delete" :href="'/delete'" @click.stop.prevent="delete" :size="'15px'" :color="'red'"/>
		<PauseButton class="pause" :href="'/pause'" @click.stop.prevent="pause" :size="'12px'" :color="'var(--secondary-color)'" :paused="data.Paused"/>
	</a>
</template>

<script>
import CloseButton from "@/components/CloseButton.vue";
import PauseButton from "@/components/PauseButton.vue";

export default {
	name: 'UnreadArticle',

	components: {
		CloseButton,
		PauseButton
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
					FeedName: "",
					Published: new Date(),
				}
			},
			validator: function(value) {
				if (typeof value.ID != "string") {
					return false
				}
				if (typeof value.Name != "string") {
					return false
				}
				if (typeof value.URL != "string") {
					return false
				}
				return true
			}
		}
	},
	
	data() {
		return {
			deleting: false,
			list: []
		}
	},
	
	methods: {
		predelete() {
			this.deleting = true
			setTimeout(() => this.deleting = false, 1000)
		},
		delete() {
			fetch("/api/feed/unsubscribe?id="+this.data.ID).then(() => this.$emit("changed"))
		},
		pause() {
			
			if (this.data.Paused) {
				fetch("/api/feed/unpause?id="+this.data.ID).then(() => this.$emit("changed"))
				return
			}
			fetch("/api/feed/pause?id="+this.data.ID).then(() => this.$emit("changed"))
		}
	}
}
</script>

<style scoped lang="scss">
.feed {
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

	.delete {
		position: absolute;
		top: 3px;
		left: calc(100% - 18px);
		background-color: var(--bg-color); /* So it overlays text, not sure if I like this. */
	}
	.pause {
		position: absolute;
		top: 22px;
		left: calc(100% - 16px);
		background-color: var(--bg-color); /* So it overlays text, not sure if I like this. */
	}
}
</style>
