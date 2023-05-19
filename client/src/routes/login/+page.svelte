<script lang="ts">
	let email = '';
	let password = '';

	function handleEmailPasswordLogin() {
		// TODO: Handle email/password login
	}

	function handleGithubLogin() {
		// TODO: Handle GitHub OAuth login
	}
</script>

<section>
	<h1>Login</h1>

	<div class="login-method">
		<h2>Email & Password</h2>

		<input type="email" bind:value={email} placeholder="Email" />
		<input type="password" bind:value={password} placeholder="Password" />

		<button on:click={handleEmailPasswordLogin}>Login with Email</button>
	</div>

	<div class="login-method">
		<h2>GitHub</h2>

		<button on:click={handleGithubLogin}>Login with GitHub</button>
	</div>
</section>

<style>
	.login-method {
		margin-bottom: 2rem;
	}
</style>
