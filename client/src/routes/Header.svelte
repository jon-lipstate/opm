<script>
	import { page } from '$app/stores';
	import logo from '$lib/images/odin-logo.svg';
	import github from '$lib/images/github.svg';
</script>

<header>
	<div class="logo">
		<a href="https://odin-lang.org">
			<img src={logo} alt="OdinLogo" />
		</a>
	</div>
	<nav>
		<input type="text" name="search" id="" placeholder="Search" />
		<button type="submit">Search</button>
	</nav>

	<div class="corner">
		<a href="https://github.com/odin-lang/Odin">
			<img src={github} alt="GitHub" />
		</a>
	</div>
</header>

<style>
	header {
		display: flex;
		justify-content: space-between;
	}

	.corner {
		width: 3em;
		height: 3em;
	}

	.corner a {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
	}

	.corner img {
		width: 2em;
		height: 2em;
		object-fit: contain;
	}
	.logo img {
		margin-left: 0.75rem;
		margin-top: 0.75rem;
		width: 6rem;
		height: 2rem;
		object-fit: contain;
	}

	nav {
		display: flex;
		justify-content: center;
		--background: rgba(255, 255, 255, 0.7);
	}

	nav a {
		display: flex;
		height: 100%;
		align-items: center;
		padding: 0 0.5rem;
		color: var(--color-text);
		font-weight: 700;
		font-size: 0.8rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		text-decoration: none;
		transition: color 0.2s linear;
	}

	a:hover {
		color: var(--color-theme-1);
	}
</style>
