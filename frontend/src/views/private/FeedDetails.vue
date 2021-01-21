<template>
	<Feed :data="details"/>
	<div class="hr"></div>
	<section name="article-list">
		<Article v-for="article in articles" :key="article" :data="article" @changed="refresh"/>
	</section>
</template>

<script>
import Article from "@/components/Article.vue";
import Feed from "@/components/Feed.vue";

export default {
	name: 'FeedDetails',
	
	components: {
		Article,
		Feed,
	},

	props: {
		id: {
			type: [String],
			default: "invalid",
		}
	},
	
	computed: {
		details() {
			if (this.rawdetails == "") {
				return {}
			}
			return JSON.parse(this.rawdetails)
		},
		articles() {
			if (this.rawarticles == "") {
				return []
			}
			return JSON.parse(this.rawarticles)
		},
		ispaused() {
			return this.details.Paused ? "(updates paused)" : ""
		},
	},
	
	data() {
		return {
			rawdetails: "",
			rawarticles: ""
		}
	},
	
	methods: {
		refresh() {
			let self = this;
			fetch("/api/feed/details?id="+this.id)
				.then(function(res) {
					if (res.ok) {
						return res.text()
					}
					throw new Error(res.status);
				})
				.then(function(res) {
					self.rawdetails = res
				})
				.catch(error => {
					console.error(error.message)
				});
			fetch("/api/feed/articles?id="+this.id)
				.then(function(res) {
					if (res.ok) {
						return res.text()
					}
					throw new Error(res.status);
				})
				.then(function(res) {
					self.rawarticles = res
				})
				.catch(error => {
					console.error(error.message)
				});
		}
	},

	created() {
		this.refresh()
	}
}
</script>

<style lang="scss">
.hr {
	background-color: var(--secondary-color);
	corner-radius: 2px;
	height: 4px;
	margin: 10px 0px;
}
</style>
