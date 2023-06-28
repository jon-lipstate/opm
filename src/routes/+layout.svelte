<script>
	import { dev } from '$app/environment';
	import { inject } from '@vercel/analytics';
	inject({ mode: dev ? 'development' : 'production' });
	//
	import Header from './Header.svelte';
	import './styles.css';
	import fjord from '$lib/images/fjord.png';
</script>

<div class="app">
	<Header />
	<span class="alpha-warning">
		<a href="https://github.com/jon-lipstate/opm/issues">Alpha Release</a>
	</span>

	<main>
		<slot />
	</main>

	<footer>
		<p>
			visit <a href="https://pkg.odin-lang.org/">pkg.odin-lang.org</a> for core/vendor packages
		</p>
	</footer>
</div>

<style>
	body {
		background-image: url({fjord});
		background-size: cover;
		background-position: center;
	}

	.alpha-warning {
		color: rgba(217, 0, 0, 0.709);
		text-align: center;
		margin-top: 1rem;
		font-style: italic;
	}

	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
		background-color: rgba(52, 58, 64, 0.8);;
	}

	main {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 1rem;
		width: 100%;
		max-width: 90vw;
		margin: 0 auto;
		box-sizing: border-box;
	}

	footer {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 12px;
	}

	footer a {
		font-weight: bold;
	}

	@media (min-width: 480px) {
		footer {
			padding: 12px 0;
		}
	}
</style>
