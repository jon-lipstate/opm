<script>
	import logo from '$lib/images/odin-logo.svg';
	// import github from '$lib/images/git_light.svg';
</script>

<header>
	<nav>
		<!-- Logo-->
		<div class="logo">
			<a href="/">
				<img src={logo} alt="OdinLogo" />
			</a>
		</div>
		<!-- Search Bar-->
		<form class="search-container" action="/search">
			<input
				type="search"
				name="query"
				class="search-input"
				placeholder="Search"
				spellcheck="false"
				autocomplete="off"
			/>
			<div class="shortcut">
				<kbd class="">Ctrl</kbd> <kbd class="">K</kbd>
			</div>
		</form>
		<!-- Right Items -->
		<div class="right-items">
			<a href="https://odin-lang.org/"> ODIN-LANG </a>
			<a href="/login"> Login </a>
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
		margin-left: 0.25rem;
		width: 5rem;
		object-fit: contain;
	}
	.search-container {
		position: relative;
		display: flex;
		align-items: center;
		width: 30rem;
		height: 100%;
		padding: 0 3px;
	}
	.search-input {
		padding: 0.2em 0.5em;
		border: 1px solid var(--back-translucent);
		font-family: inherit;
		font-size: 1.1rem;
		appearance: none;
		-webkit-appearance: none;
		width: 100%;
		height: 1.3em;
		border-radius: var(--border-radius);
	}
	.shortcut {
		position: absolute;
		top: calc(50%-0.8rem);
		right: 1rem;
		width: 100%;
		text-align: right;
		pointer-events: none;
		font-size: 1.2rem;
		text-transform: uppercase;
	}
	.right-items {
		margin-right: 0.25rem;
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
	}
	.right-items > a {
		margin-right: 0.5rem;
	}
	a {
		color: var(--color-text);
	}
</style>
