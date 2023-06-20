<script lang="ts">
	import { onMount } from 'svelte';

	let isOpen = false;
	function toggleMenu() {
		isOpen = !isOpen;
	}

	function handleLogout() {
		isOpen = false;
	}

	onMount(() => {
		const closeMenu = (event: any) => {
			if (!event.target.closest('.menu') && !event.target.closest('.hamburger')) {
				isOpen = false;
			}
		};
		if (typeof window !== 'undefined') {
			window.addEventListener('click', closeMenu);
			return () => {
				window.removeEventListener('click', closeMenu);
			};
		}
	});
</script>

<nav>
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<div class="hamburger" on:click={toggleMenu}>
		<div class="line" class:top-line={isOpen} />
		<div class="line" class:mid-line={isOpen} />
		<div class="line" class:bot-line={isOpen} />
	</div>

	{#if isOpen}
		<div class="menu">
			<a on:click={toggleMenu} href="/">Search</a>
			<a on:click={toggleMenu} href="/browse">Browse Packages</a>
			<!-- <a on:click={toggleMenu} href="/new">New Package</a> -->
			<a on:click={toggleMenu} href="/manage">Manage Packages</a>
			<!-- {#if $isAdmin}
          <a on:click={toggleMenu} href="/admin">Admin</a>
      {/if} -->
			<a on:click={toggleMenu} href="/account">Account Tokens</a>
			<a on:click={handleLogout} href="/auth/signout">Logout</a>
		</div>
	{/if}
</nav>

<style>
	.hamburger {
		display: flex;
		flex-direction: column;
		justify-content: space-around;
		width: 2rem;
		height: 2rem;
		cursor: pointer;
	}

	.line {
		width: 2rem;
		height: 0.25rem;
		background: var(--color-theme-4);
		transition: all 0.3s ease;
	}

	.top-line {
		transform: translateY(0.675rem) rotate(-45deg);
	}

	.mid-line {
		opacity: 0;
	}

	.bot-line {
		transform: translateY(-0.675em) rotate(45deg);
	}

	.menu {
		display: flex;
		flex-direction: column;
		position: absolute;
		right: 0;
		width: 100%;
		background: var(--color-bg-1);
		text-align: center;
		z-index: 1;
	}

	.menu a {
		padding: 1rem;
		color: var(--color-theme-4);
		text-decoration: none;
	}
</style>
