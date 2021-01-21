<template>
	<div>
		<p>Before logging in you must confirm your email address. You should shortly recive an email containing instructions.</p>

		<p>If you never get your email then <a href="&#109;&#97;&#105;&#108;&#116;&#111;&#58;&#109;&#105;&#108;&#111;&#64;&#104;&#116;&#116;&#112;&#115;&#99;&#111;&#108;&#111;&#110;&#115;&#108;&#97;&#115;&#104;&#115;&#108;&#97;&#115;&#104;&#119;&#119;&#119;&#46;&#99;&#111;&#109;">send me an email</a> and I will figure out what broke this time.</p>

		<router-link to="/">Return to Login Page</router-link>
	</div>
</template>

<script>
export default {
	name: 'ConfirmEmail',
}
</script>

<style lang="scss">
	p {
		color: var(--font-color);
	}
	a {
		color: var(--secondary-color);
	}
</style>
