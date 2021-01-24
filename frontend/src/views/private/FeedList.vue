<template>
	<section name="body">
		<section name="feedlist">
			<Feed v-for="feed in list" :key="feed" :data="feed" @changed="refresh"/>
		</section>
		<AddFeed @added="refresh"/>
	</section>
</template>

<script>
import AddFeed from "@/components/AddFeed.vue";
import Feed from "@/components/Feed.vue";

export default {
	name: 'FeedList',

	components: {
		AddFeed,
		Feed
	},

	computed: {
		list() {
			if (this.rawlist == "") {
				return []
			}
			return JSON.parse(this.rawlist)
		}
	},
	
	data() {
		return {
			rawlist: ""
		}
	},

	methods: {
		refresh() {
			let self = this;
			fetch("/api/feed/list")
				.then(function(res) {
					if (res.ok) {
						return res.text()
					}
					throw new Error(res.status);
				})
				.then(function(res) {
					self.rawlist = res
				})
				.catch(error => {
					console.error(error.message)
					self.addstate = false
					setTimeout(() => self.addstate = null, 5000)
				});
		}
	},

	created() {
		this.refresh()
	}
}
</script>

<style scoped lang="scss">
	section[name=body] {
		display: flex;
		flex-direction: column;
		flex: 1;
	}
</style>
