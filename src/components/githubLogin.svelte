<script lang="ts">
    import gh from '$lib/images/github.svg';
    import { goto } from '$app/navigation';

	function handleGithubLogin() {
        goto(`/auth/signin`)
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
