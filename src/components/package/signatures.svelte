<script lang="ts">
	export let structs: string = `Context :: struct {
  entities: Entities,
  component_map: map[typeid]Component_List,
}`;
	export let procs: string = `init_ecs :: proc() -> (ctx: Context)`;
</script>

<section class="signatures">
	<div style="color:orangered">TODO: Use the core language doc format and data</div>
	<div>
		<h3>Structures:</h3>
		<code>
			{@html structs}
		</code>
	</div>
	<div>
		<h3>Procedures:</h3>
		<code>
			{@html procs}
		</code>
	</div>
</section>

<style>
	.signatures {
		margin-top: 1rem;
	}
</style>
