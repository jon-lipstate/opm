<script lang="ts">
    import gh from '$lib/images/github.svg';
	import { isLoggedIn } from '$stores/user';
    import { goto } from '$app/navigation';

	function handleGithubLogin() {
		// TODO: Handle GitHub OAuth login
        // Call your Flask server's /login endpoint
		console.warn( `goto ${import.meta.env.VITE_API_HOST}/auth/github`);
        goto(`${import.meta.env.VITE_API_HOST}/auth/github`)
		isLoggedIn.set(true);
	}
</script>

<button on:click={handleGithubLogin}>Login<img src={gh} alt="github logo" class="github-logo"/></button>

<style>
	button {
		margin:auto, 0;
		font-size: 1rem;
		display: flex; 
    	align-items: center; 
    	justify-content: center; 
		/* background-color:var(--color-theme-2); */
		/* color:#f7f7f77f; */
	}
	.github-logo {
		height: 1.4rem;
		padding-left: 0.2rem;
		vertical-align: middle; 
	}
</style>
