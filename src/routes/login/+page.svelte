<script lang="ts">
	import gh from '$lib/images/github.svg';

	function handleGithubLogin() {
		// TODO: Handle GitHub OAuth login
	}
</script>

<section>
	<h1>Login</h1>

	<div class="login-method">
		<button on:click={handleGithubLogin}
			>Sign in with GitHub <img src={gh} alt="github logo" class="github-logo" /></button
		>
	</div>
</section>

<style>
	button {
		margin: auto, 0;
		font-size: 1rem;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.github-logo {
		height: 2rem;
		padding-left: 0.2rem;
		vertical-align: middle;
	}
</style>
