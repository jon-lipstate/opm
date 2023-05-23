<script>
	import { isLoggedIn } from '$stores/user';
	import logo from '$lib/images/odin-logo.svg';
	import GithubLogin from '$components/githubLogin.svelte';
	import Menu from '$components/hamburgerMenu.svelte';

</script>

<header>
	<nav>
		<!-- Logo-->
		<div class="logo">
			<a href="/">
				<img src={logo} alt="OdinLogo" />
			</a>
		</div>
		<!-- Right Items -->
		<div class="right-items">
			{#if !$isLoggedIn}
			<GithubLogin />

			{:else}
				<Menu/>
			{/if}
			<!-- <a href="/login"> Login </a> -->
		</div>
	</nav>
</header>

<style>
	nav {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		background-color: var(--color-bg-0);
		margin: 0;
		min-height: 40px;
	}
	.logo {
		margin-left: 0.75rem;
		margin-top: 0.25rem;
		width: 5rem;
		object-fit: contain;
	}
	.logo img {
		max-height: 60px;
	}

	.right-items {
		margin-right: 0.25rem;
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
	}
	/* .right-items > a {
		margin-right: 0.5rem;
	} */
	a {
		color: var(--color-text);
	}
</style>
