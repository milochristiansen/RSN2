<template>
	<section name="body">
		<section name="unreadlist">
			<UnreadArticle v-for="article in list" :key="article" :data="article"/>
		</section>
		<AddFeed/>
	</section>
</template>

<script>
import ReconnectingWebSocket from 'reconnecting-websocket';

import UnreadArticle from "@/components/UnreadArticle.vue";
import AddFeed from "@/components/AddFeed.vue";

export default {
	name: 'Unread',

	components: {
		UnreadArticle,
		AddFeed
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
			socket: null,
			rawlist: ""
		}
	},

	methods: {
		refresh(message) {
			this.rawlist = message.data
		}
	},

	created() {
		this.socket = new ReconnectingWebSocket("ws://localhost:3366/api/article/feed", [], {
			connectionTimeout: 20000
		})
		this.socket.addEventListener("message", this.refresh)
	}
}
</script>

<style scoped lang="scss">
	section[name=body] {
		display: flex;
		flex-direction: column;
		height: 100%;
	}
</style>
