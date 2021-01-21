module.exports = {
	devServer: {
		proxy: {
			'^/api/article/feed': {
				target: "ws://localhost:3366/api/article/feed",
				ws: true
			},
			'^/api': {
				target: "http://localhost:3366",
				changeOrigin: true
			}
		}
	}
};
